package httpagent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/chaosblade-io/chaosblade/data"
)

func experimentPlanHandler(w http.ResponseWriter, req *http.Request) {
	source := data.GetSource()
	flag, err := source.QueryExistExperimentPlanModelByStatus("0")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Experiment query error, %s", err.Error())))
		return
	}
	if flag {
		w.WriteHeader(500)
		w.Write([]byte("Exist Experiment with status 0."))
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Parse http request error, %s", err.Error())))
		return
	}
	wf := &workflow{}
	if err = json.Unmarshal(body, wf); err != nil {
		w.WriteHeader(415)
		w.Write([]byte(fmt.Sprintf("Parse json parameters error, %s", err.Error())))
		return
	}
	nowTime := time.Now()
	if err = wf.Parse(nowTime); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	if err = wf.insertPlan(); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Set experiment plan error, %s", err.Error())))
		return
	}
	w.Write([]byte(`ok`))
}

type workflow struct {
	JobId         string                     `json:"job_id"`         //chaos-job id
	PortalTime    int64                      `json:"portal_time"`    //chaos platform systime
	FlowStarttime int64                      `json:"flow_starttime"` //chaos flow start time
	Events        []data.ExperimentPlanModel `json:"events"`         //ExperimentPlan list
	//Prepare []prepare `json:"prepare"`
}

const Ms = 1000

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
	nodetimeGap := now/int64(time.Millisecond) - wf.PortalTime
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
		if wf.Events[i].StartOffset == 0 {
			return errors.New("Parse parameters error, StartOffset is invalid.")
		}
		if wf.Events[i].EndOffset == 0 {
			return errors.New("Parse parameters error, EndOffset is invalid.")
		}
		wf.Events[i].Status = "0"
		//starttime(node)
		startTime := wf.FlowStarttime + wf.Events[i].StartOffset + nodetimeGap
		startSec := startTime / 1000
		startNsec := (startTime % 1000) * int64(time.Millisecond)
		wf.Events[i].StartTime = time.Unix(startSec, startNsec)
		//endtime(node)
		endTime := wf.FlowStarttime + wf.Events[i].EndOffset + nodetimeGap
		endSec := endTime / 1000
		endNsec := (endTime % 1000) * int64(time.Millisecond)
		wf.Events[i].EndTime = time.Unix(endSec, endNsec)
		//createtime updatetime
		wf.Events[i].CreateTime = time.Now()
		wf.Events[i].UpdateTime = time.Now()

		wf.Events[i].JobId = wf.JobId
		wf.Events[i].NodetimeGap = nodetimeGap
	}
	return nil
}

func (wf *workflow) insertPlan() error {
	source := data.GetSource()
	src, ok := source.(*data.Source)
	if !ok {
		return errors.New("Source is not implements SourceI")
	}
	tx, err := src.DB.Begin()
	if err != nil {
		return err
	}
	for _, v := range wf.Events {
		err := data.InsertExperimentPlanModel(tx, &v)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
