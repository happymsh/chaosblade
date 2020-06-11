package task

import (
	"context"
	"testing"
	"time"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"

	"github.com/chaosblade-io/chaosblade/data"

	"github.com/chaosblade-io/chaosblade/agent/pdata"
)

func TestExpTaskArrange(t *testing.T) {
	now := time.Now()
	startTime := now.Add(-time.Second * 10)
	endTime := now.Add(time.Second * 10)
	planModel := &pdata.ExperimentPlanModel{
		ExpId:       "E1",
		JobId:       "J1",
		NodetimeGap: 0,
		Command:     "cpu",
		SubCommand:  "fullload",
		Flag:        "cpu-percent=10 cpu-count=1",
		Status:      pdata.ToBeProcessed,
		Error:       "",
		StartOffset: 10000,
		EndOffset:   20000,
		StartTime:   startTime,
		EndTime:     endTime,
		CreateTime:  time.Time{},
		UpdateTime:  time.Time{},
	}
	planModels := make([]*pdata.ExperimentPlanModel, 0, 1)
	planModels = append(planModels, planModel)
	model := &data.ExperimentModel{}
	queryExperimentPlanModelOrderByStartTime = func() (models []*pdata.ExperimentPlanModel, e error) {
		if planModel.Status == pdata.ToBeProcessed || planModel.Status == pdata.Processing {
			return planModels, nil
		}
		return nil, nil
	}
	updateExperimentPlanModelByExpId = func(expId, uid, status, errMsg string) error {
		if expId == planModel.ExpId {
			planModel.Uid = uid
			planModel.Status = status
			planModel.Error = errMsg
		}
		return nil
	}
	queryExperimentModelByUid = func(uid string) (*data.ExperimentModel, error) {
		if model.Uid == uid {
			return model, nil
		}
		return nil, nil
	}
	executeExp = func(ctx context.Context, command, subcommand, flags, timeout string) *spec.Response {
		model.Uid = "c19aca7985dbdd8d"
		model.Status = "Success"
		return &spec.Response{
			Code:    200,
			Success: true,
			Err:     "",
			Result:  `{"code":200,"success":true,"result":"c19aca7985dbdd8d"}`,
		}
	}
	cRun = func(ctx context.Context, script, args string) *spec.Response {
		model.Status = "Destroyed"
		return &spec.Response{
			Code:    200,
			Success: true,
			Err:     "",
			Result:  `{"code":200,"success":true,"result":"c19aca7985dbdd8d"}`,
		}
	}

	ExpTaskArrange(context.TODO(), false)
	if planModel.Status != pdata.Processing {
		t.Error("model.Status is not Processing ", planModel)
	}
	if planModel.Uid != "c19aca7985dbdd8d" {
		t.Error("model.Uid is wrong ", planModel)
	}

	ExpTaskArrange(context.TODO(), false)
	t.Log("sencond call.")
	if planModel.Status != pdata.Processing {
		t.Error("model.Status is not Processing ", planModel)
	}
	if planModel.Uid != "c19aca7985dbdd8d" {
		t.Error("model.Uid is wrong ", planModel)
	}

	time.Sleep(time.Second * 11)

	ExpTaskArrange(context.TODO(), false)

	if planModel.Status != pdata.SuccessfulProcessing {
		t.Error("model.Status is not SuccessfulProcessing ", planModel)
	}
}
