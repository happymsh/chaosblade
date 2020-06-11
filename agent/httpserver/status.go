package httpserver

import (
	"encoding/json"
	"fmt"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/chaosblade-io/chaosblade/agent/metrics"
	"github.com/chaosblade-io/chaosblade/agent/pdata"
	"net/http"
	"time"
)

type expStatus struct {
	ExpId     string `json:"event_id"`
	ExpStatus string `json:"event_status"`
	Uid       string `json:"uid"`
	ErrMsg    string `json:"err_msg"`
}

type statusResponse struct {
	JobId       string                  `json:"job_id"`
	NodetimeGap int64                   `json:"nodetime_gap"`
	Finished    string                  `json:"finished"`
	Events      []expStatus             `json:"events"`
	Metrics     metrics.NodeStatusModel `json:"metrics"`
}

const (
	FLOW_TODO      = "0"
	FLOW_DOING     = "1"
	FLOW_DONE      = "2"
	FLOW_TERMINATE = "3"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(spec.ReturnFail(spec.Code[spec.ServerError],
			fmt.Sprintf("statusHandler ParseForm error, %s", err.Error())).ToString()))
		return
	}

	jobId := r.Form.Get("job_id")
	pSource := pdata.GetSource()
	planModels, err := pSource.QueryExperimentPlanModelByJobId(jobId)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(spec.ReturnFail(spec.Code[spec.ServerError],
			fmt.Sprintf("statusHandler QueryExperimentPlanModelByJobId error, %s", err.Error())).ToString()))
		return
	}
	sr := &statusResponse{}
	select {
	case _ = <-metrics.Metrics.Channel:
		sr.Metrics = *metrics.Metrics.NodeStatusModel
		metrics.Metrics.Channel <- 1
	}
	if len(planModels) != 0 {
		sr.JobId = planModels[0].JobId
		sr.NodetimeGap = planModels[0].NodetimeGap
		sr.Events = make([]expStatus, 0)
		for _, planModel := range planModels {
			exps := expStatus{
				ExpId:     planModel.ExpId,
				Uid:       planModel.Uid,
				ExpStatus: planModel.Status,
				ErrMsg:    planModel.Error,
			}
			sr.Events = append(sr.Events, exps)
		}
	}
	//add by kfzx-wumg's request.
	if sr.JobId == "" {
		sr.JobId = jobId
	}
	if sr.Events == nil {
		sr.Events = make([]expStatus, 0)
	}

	sr.Finished = FLOW_TODO
	// calculate Finished by events' status
	isAllExpDone := true
	for _, v := range sr.Events {
		if v.ExpStatus == pdata.ToBeProcessed || v.ExpStatus == pdata.Processing {
			isAllExpDone = false
		}
	}
	// calculate Finished by duration
	var startTime int64 = -1
	startTime = viper.GetInt64("FlOW_START_TIME")
	var endTime int64 = -1
	endTime = viper.GetInt64("FlOW_END_TIME")
	gp := viper.GetInt64("GAP")
	startTime += gp
	endTime += gp
	nowTime := time.Now().UnixNano() / int64(time.Millisecond)
	if nowTime < startTime {
		sr.Finished = FLOW_TODO
	} else if nowTime > endTime {
		if isAllExpDone {
			sr.Finished = FLOW_DONE
		} else {
			sr.Finished = FLOW_DOING
		}
	} else {
		sr.Finished = FLOW_DOING
	}
	if viper.GetBool("FLOW_TERMINATE") {
		sr.Finished = FLOW_TERMINATE
	}

	bytes, err := json.MarshalIndent(sr, "", "\t")
	if err != nil {
		logrus.Error("MetricsCollect by Marshal :", fmt.Sprintf("%+v", sr))
		w.WriteHeader(500)
		_, _ = w.Write([]byte(spec.ReturnFail(spec.Code[spec.ServerError],
			fmt.Sprintf("statusResponse Marshal error, %s", err.Error())).ToString()))
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		logrus.Error("statusHandler ResponseWriter error :", err.Error())
	}
}
