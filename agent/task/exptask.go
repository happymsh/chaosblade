package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
	"path"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/chaosblade-io/chaosblade/agent/exec"
	"github.com/chaosblade-io/chaosblade/agent/pdata"
)

var mutex sync.RWMutex = sync.RWMutex{}

var pSource = pdata.GetSource()
var queryAllToProcessExperimentPlanModel = pSource.QueryAllToProcessExperimentPlanModel
var queryExperimentPlanModelOrderByStartTime = pSource.QueryToProcessExperimentPlanModelOrderByStartTime
var updateExperimentPlanModelByExpId = pSource.UpdateExperimentPlanModelStatusByExpId
var queryExperimentModelByUid = pSource.QueryExperimentModelByUid
var executeExp = exec.ExecuteExp
var cRun = exec.CRun

//non-support multi-goroutine
func ExpTaskArrange(ctx context.Context, terminalFlag bool) *spec.Response {
	mutex.Lock()
	defer mutex.Unlock()

	if terminalFlag {
		var destroyUidList []string
		var unDestroyUidList []string
		planModels, err := queryAllToProcessExperimentPlanModel()
		if err == nil {
			for _, planModel := range planModels {
				index := fmt.Sprintf("%s|%s", planModel.ExpId, planModel.Uid)
				if planModel.Uid != "" {
					args := fmt.Sprintf("destroy %s", planModel.Uid)
					ctxTimeout, _ := context.WithTimeout(ctx, 60*time.Second)
					response := cRun(ctxTimeout, path.Join(util.GetProgramPath(), exec.BLADE), args)
					if response.Success {
						err = updateExperimentPlanModelByExpId(planModel.ExpId, planModel.Uid, pdata.Terminate, "")
						if err != nil {
							unDestroyUidList = append(unDestroyUidList, index)
							logrus.Error("ExpTaskArrange error by updateExperimentPlanModelByExpId when terminate a exp ", planModel.Uid, err)
						} else {
							destroyUidList = append(destroyUidList, index)
						}
					} else {
						unDestroyUidList = append(unDestroyUidList, index)
						logrus.Error("ExpTaskArrange error by timinate Exp ", planModel.Uid, response)
					}
				} else {
					err = updateExperimentPlanModelByExpId(planModel.ExpId, planModel.Uid, pdata.Terminate, "")
					if err != nil {
						unDestroyUidList = append(unDestroyUidList, index)
						logrus.Error("ExpTaskArrange error by updateExperimentPlanModelByExpId when terminate a exp ", planModel.Uid, err)
					} else {
						destroyUidList = append(destroyUidList, index)
					}
				}

			}
		} else {
			unDestroyUidList = append(unDestroyUidList, "ALL")
			logrus.Error("ExpTaskArrange error by queryAllToProcessExperimentPlanModel when terminate.", err)
		}
		result := fmt.Sprintf("Terminated exp(event_id|uid) list is %+v , UnTerminated exp(event_id|uid) list is %+v. note: uid can be empty", destroyUidList, unDestroyUidList)
		if unDestroyUidList == nil {
			return spec.ReturnSuccess(result)
		} else {
			return spec.ReturnFail(spec.Code[spec.DeployError], result)
		}
	}

	//result <-- select exp_plan where systime>starttime and status in (0,1) order by starttime
	//for range result:
	//	if status is 1 then destroy the exp and update status
	//	if status is 0 then create the exp and update status
	experimentPlanModels, err := queryExperimentPlanModelOrderByStartTime()
	if err != nil {
		logrus.Error("ExpTaskArrange error by queryExperimentPlanModelOrderByStartTime:", pdata.ToBeProcessed, err)
		return nil
	}
	for _, planModel := range experimentPlanModels {
		if planModel.Status == pdata.Processing {
			if planModel.EndTime.Before(time.Now()) {
				logrus.Info("destroy a exp:", fmt.Sprintf("%+v", planModel))
				if planModel.Uid == "" {
					logrus.Error("destroy failed, beacause planModel.Uid is empty")
					break
				}
				experimentModel, err := queryExperimentModelByUid(planModel.Uid)
				if err != nil {
					logrus.Error("ExpTaskArrange error by queryExperimentModelByUid,", planModel.Uid, err)
					break
				}
				if experimentModel == nil {
					logrus.Error("ExpTaskArrange error by queryExperimentModelByUid, experimentModel is empty ")
					break
				}
				args := fmt.Sprintf("destroy %s", experimentModel.Uid)
				ctxTimeout, _ := context.WithTimeout(ctx, 60*time.Second)
				response := cRun(ctxTimeout, path.Join(util.GetProgramPath(), exec.BLADE), args)
				if !response.Success {
					logrus.Error("ExpTaskArrange error by destroy Exp:", response)
					break
				}
				err = updateExperimentPlanModelByExpId(planModel.ExpId, planModel.Uid, pdata.SuccessfulProcessing, "")
				if err != nil {
					logrus.Error("ExpTaskArrange error by updateExperimentPlanModelByExpId when Destroyed a exp.", err)
					break
				}
			}
		} else if planModel.Status == pdata.ToBeProcessed {
			if planModel.EndTime.Before(time.Now()) {
				logrus.Info("this exp is expire:", fmt.Sprintf("%+v", planModel))
				err = updateExperimentPlanModelByExpId(planModel.ExpId, "", pdata.ProcessingFailure, fmt.Sprintf("exp is expire in %+v", time.Now()))
				if err != nil {
					logrus.Error("ExpTaskArrange error by updateExperimentPlanModelByExpId when exp is expire.", err)
					break
				}
				break
			}
			logrus.Info("start a exp:", fmt.Sprintf("%+v", planModel))
			ctxTimeout, _ := context.WithTimeout(ctx, 60*time.Second)
			response := executeExp(ctxTimeout, planModel.Command, planModel.SubCommand, planModel.Flag, "")
			if !response.Success {
				logrus.Error("ExpTaskArrange error by ExecuteExp:", response)
				break
			}
			result, ok := response.Result.(string)
			if !ok {
				logrus.Error("ExpTaskArrange error by response.Result assert")
				break
			}
			r := &exec.Result{}
			err = json.Unmarshal([]byte(result), r)
			if err != nil {
				logrus.Error("ExpTaskArrange error by Unmarshal response.result:", err)
				break
			}
			err = updateExperimentPlanModelByExpId(planModel.ExpId, r.Result, pdata.Processing, "")
			if err != nil {
				logrus.Error("ExpTaskArrange error by updateExperimentPlanModelByExpId.", err)
				break
			}
		}
	}
	return nil
}
