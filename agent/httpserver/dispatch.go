package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/chaosblade-io/chaosblade/agent/pdata"
)

func experimentPlanHandler(w http.ResponseWriter, req *http.Request) {
	source := pdata.GetSource()
	flag, err := source.QueryExistExperimentPlanModelByStatus("0")
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(spec.ReturnFail(spec.Code[spec.ServerError],
			fmt.Sprintf("Experiment query error, %s", err.Error())).ToString()))
		return
	}
	if flag {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(spec.ReturnFail(spec.Code[spec.IllegalParameters],
			"Exist Experiment with status 0 witch means exp to be processed, can not receive other Experiment").ToString()))
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(spec.ReturnFail(spec.Code[spec.ServerError],
			fmt.Sprintf("Parse http request error, %s", err.Error())).ToString()))
		return
	}
	wf := &workflow{}
	if err = json.Unmarshal(body, wf); err != nil {
		w.WriteHeader(415)
		_, _ = w.Write([]byte(spec.ReturnFail(spec.Code[spec.ServerError],
			fmt.Sprintf("Parse json parameters error, %s", err.Error())).ToString()))
		return
	}
	nowTime := time.Now()
	if err = wf.Parse(nowTime); err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(spec.ReturnFail(spec.Code[spec.ServerError],
			err.Error()).ToString()))
		return
	}
	if err = wf.insertPlan(); err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(spec.ReturnFail(spec.Code[spec.ServerError],
			fmt.Sprintf("Set experiment plan error [ experiment_plan.exp_id is corresponding event_id], %s", err.Error())).ToString()))
		return
	}
	//write config file.
	viper.Set("FLOW", wf)
	viper.Set("GAP", NodetimeGap)
	viper.Set("FlOW_START_TIME", wf.FlowStarttime)
	viper.Set("FlOW_END_TIME", wf.FlowStarttime+wf.Duration)
	viper.Set("FLOW_TERMINATE", false)
	err = viper.WriteConfig()
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("experimentPlanHandler request body", string(body))
	_, _ = w.Write([]byte(spec.ReturnSuccess("ok").ToString()))
}

type workflow struct {
	JobId         string                      `json:"job_id"`         //chaos-job id
	PortalTime    int64                       `json:"portal_time"`    //chaos platform systime
	FlowStarttime int64                       `json:"flow_starttime"` //chaos flow start time
	Duration      int64                       `json:"duration"`       //chaos flow duration time
	Events        []pdata.ExperimentPlanModel `json:"events"`         //ExperimentPlan list
	//Prepare []prepare `json:"prepare"`
}

const Ms = 1000

var NodetimeGap int64

func (wf *workflow) Parse(t time.Time) error {
	if wf.JobId == "" {
		return errors.New("Parse parameters error, JobId is empty.")
	}
	if wf.PortalTime == 0 {
		return errors.New("Parse parameters error, PortalTime is not init.")
	}
	if wf.FlowStarttime == 0 {
		return errors.New("Parse parameters error, FlowStarttime is not init.")
	}
	if len(wf.Events) == 0 {
		return errors.New("Parse parameters error, Events is empty.")
	}
	now := t.UnixNano()
	//calculate NodetimeGap , ms
	NodetimeGap = now/int64(time.Millisecond) - wf.PortalTime
	for i, _ := range wf.Events {
		if wf.Events[i].ExpId == "" {
			return errors.New("Parse parameters error, ExpId is empty.")
		}
		if wf.Events[i].Command == "" {
			return errors.New("Parse parameters error, Command is empty.")
		}
		if wf.Events[i].SubCommand == "" {
			return errors.New("Parse parameters error, SubCommand is empty.")
		}
		if wf.Events[i].Flag == "" {
			return errors.New("Parse parameters error, Flag is empty.")
		}
		//if wf.Events[i].StartOffset == 0 {
		//	return errors.New("Parse parameters error, StartOffset is invalid.")
		//}
		if wf.Events[i].EndOffset == 0 {
			return errors.New("Parse parameters error, EndOffset is invalid.")
		}
		//starttime(node)
		startTime := wf.FlowStarttime + wf.Events[i].StartOffset + NodetimeGap
		startSec := startTime / 1000
		startNsec := (startTime % 1000) * int64(time.Millisecond)
		wf.Events[i].StartTime = time.Unix(startSec, startNsec)
		//endtime(node)
		endTime := wf.FlowStarttime + wf.Events[i].EndOffset + NodetimeGap
		endSec := endTime / 1000
		endNsec := (endTime % 1000) * int64(time.Millisecond)
		wf.Events[i].EndTime = time.Unix(endSec, endNsec)
		//createtime updatetime
		wf.Events[i].CreateTime = time.Now()
		wf.Events[i].UpdateTime = time.Now()

		wf.Events[i].JobId = wf.JobId
		wf.Events[i].NodetimeGap = NodetimeGap
		wf.Events[i].Status = "0"
	}
	return nil
}

func (wf *workflow) insertPlan() error {
	psource := pdata.GetSource()
	tx, err := psource.DB.Begin()
	if err != nil {
		return err
	}
	for _, v := range wf.Events {
		err := pdata.InsertExperimentPlanModel(tx, &v)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
