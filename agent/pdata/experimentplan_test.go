package pdata

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/chaosblade-io/chaosblade/data"
)

const testdataFile = "chaosblade_test.dat"

func TestSource_ExperimentPlanTableExists(t *testing.T) {
	database, err := sql.Open("sqlite3", testdataFile)
	defer database.Close()
	if err != nil {
		t.Error(err)
		return
	}
	source := &PSource{
		Source: data.Source{DB: database},
	}
	exists, err := source.ExperimentPlanTableExists()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("TestSource_ExperimentPlanTableExists:", exists)
}

func TestSource_InitExperimentPlanTable(t *testing.T) {
	database, err := sql.Open("sqlite3", testdataFile)
	defer database.Close()
	if err != nil {
		t.Error(err)
		return
	}
	source := &PSource{
		Source: data.Source{DB: database},
	}
	exists, err := source.ExperimentPlanTableExists()
	if err != nil {
		t.Error(err)
		return
	}
	if exists {
		t.Log("drop table experiment_plan")
		source.DB.Exec("drop table experiment_plan")
	}
	err = source.InitExperimentPlanTable()
	if err != nil {
		t.Error(err)
		return
	}
	exists, err = source.ExperimentPlanTableExists()
	if err != nil {
		t.Error(err)
		return
	}
	if !exists {
		t.Error("InitExperimentPlanTable Init failed.")
	}
}

func TestSource_InsertExperimentPlanModel(t *testing.T) {
	expPlan := &ExperimentPlanModel{
		ExpId:       "E3",
		JobId:       "J1",
		NodetimeGap: 6000,
		Command:     "cpu",
		SubCommand:  "fullload",
		Flag:        "--cpu-count=2 --cpu-percent=20",
		Status:      "0",
		Error:       "",
		StartTime:   time.Now(),
		EndTime:     time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	database, err := sql.Open("sqlite3", testdataFile)
	defer database.Close()
	if err != nil {
		t.Error(err)
		return
	}
	source := &PSource{
		Source: data.Source{DB: database},
	}
	source.CheckAndInitExperimentPlanTable()
	source.DB.Exec("delete from experiment_plan where exp_id='E3'")
	tx, err := source.DB.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	err = InsertExperimentPlanModel(tx, expPlan)
	if err != nil {
		t.Error(err)
		return
	}
	tx.Commit()
	result, err := source.QueryExperimentPlanModelByExpId("E3")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)

	experimentPlanModels, err := source.QueryExperimentPlanModelByStatus("0")
	if err != nil {
		t.Error(err)
		return
	}
	if len(experimentPlanModels) == 0 {
		t.Error("experimentPlanModels len is 0.")
	}
	for _, model := range experimentPlanModels {
		if model.Status != "0" {
			t.Error("experimentPlanModels has model with status 0.")
		}
	}
}

func TestSource_UpdateExperimentPlanModelByExpId(t *testing.T) {
	database, err := sql.Open("sqlite3", testdataFile)
	defer database.Close()
	if err != nil {
		t.Error(err)
		return
	}
	source := &PSource{
		Source: data.Source{DB: database},
	}
	err = source.UpdateExperimentPlanModelStatusByExpId("E3", "sljdfnsxojdf", "2", "cpuburn execute error.")
	if err != nil {
		t.Error(err)
		return
	}
	result, err := source.QueryExperimentPlanModelByExpId("E3")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
	if result.Status != "2" {
		t.Error("UpdateExperimentPlanModelStatusByExpId not worked for status.")
	}
	if result.Error == "" {
		t.Error("UpdateExperimentPlanModelStatusByExpId not worked for error.")
	}
}

func TestPSource_QueryExperimentPlanModelOrderByStartTime(t *testing.T) {
	database, err := sql.Open("sqlite3", testdataFile)
	defer database.Close()
	if err != nil {
		t.Error(err)
		return
	}
	source := &PSource{
		Source: data.Source{DB: database},
	}
	E0 := &ExperimentPlanModel{
		ExpId:       "E0",
		JobId:       "J1",
		NodetimeGap: 6000,
		Command:     "cpu",
		SubCommand:  "fullload",
		Flag:        "--cpu-count=1 --cpu-percent=20",
		Status:      "0",
		Error:       "",
		StartTime:   time.Now().Add(time.Second * 30),
		EndTime:     time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	E1 := &ExperimentPlanModel{
		ExpId:       "E1",
		JobId:       "J1",
		NodetimeGap: 6000,
		Command:     "cpu",
		SubCommand:  "fullload",
		Flag:        "--cpu-count=1 --cpu-percent=20",
		Status:      "1",
		Error:       "",
		StartTime:   time.Now().Add(-time.Second * 3),
		EndTime:     time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	E2 := &ExperimentPlanModel{
		ExpId:       "E2",
		JobId:       "J1",
		NodetimeGap: 6000,
		Command:     "cpu",
		SubCommand:  "fullload",
		Flag:        "--cpu-count=2 --cpu-percent=20",
		Status:      "0",
		Error:       "",
		StartTime:   time.Now().Add(-time.Second * 4),
		EndTime:     time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	E3 := &ExperimentPlanModel{
		ExpId:       "E3",
		JobId:       "J1",
		NodetimeGap: 6000,
		Command:     "cpu",
		SubCommand:  "fullload",
		Flag:        "--cpu-count=1 --cpu-percent=20",
		Status:      "3",
		Error:       "",
		StartTime:   time.Now().Add(-time.Second * 2),
		EndTime:     time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	E4 := &ExperimentPlanModel{
		ExpId:       "E4",
		JobId:       "J1",
		NodetimeGap: 6000,
		Command:     "cpu",
		SubCommand:  "fullload",
		Flag:        "--cpu-count=2 --cpu-percent=20",
		Status:      "2",
		Error:       "",
		StartTime:   time.Now().Add(-time.Second * 1),
		EndTime:     time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	experimentPlanModels := make([]*ExperimentPlanModel, 0)
	experimentPlanModels = append(experimentPlanModels, E0)
	experimentPlanModels = append(experimentPlanModels, E1)
	experimentPlanModels = append(experimentPlanModels, E2)
	experimentPlanModels = append(experimentPlanModels, E3)
	experimentPlanModels = append(experimentPlanModels, E4)
	_, _ = source.DB.Exec("delete from experiment_plan")
	tx, err := source.DB.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	err = InsertExperimentPlanModel(tx, E0)
	if err != nil {
		t.Error(err)
		return
	}
	err = InsertExperimentPlanModel(tx, E1)
	if err != nil {
		t.Error(err)
		return
	}
	err = InsertExperimentPlanModel(tx, E2)
	if err != nil {
		t.Error(err)
		return
	}
	err = InsertExperimentPlanModel(tx, E3)
	if err != nil {
		t.Error(err)
		return
	}
	err = InsertExperimentPlanModel(tx, E4)
	if err != nil {
		t.Error(err)
		return
	}
	_ = tx.Commit()
	type args struct {
		status []string
	}
	var tests = []struct {
		name    string
		args    args
		want    []*ExperimentPlanModel
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{[]string{"0", "1"}},
			want:    experimentPlanModels,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := source.QueryToProcessExperimentPlanModelOrderByStartTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryToProcessExperimentPlanModelOrderByStartTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(fmt.Sprintf("got[0] is %+v", got[0]))
			t.Log(fmt.Sprintf("got[1] is %+v", got[1]))
			if got[0].ExpId != "E2" {
				t.Errorf("QueryToProcessExperimentPlanModelOrderByStartTime() got = %v, want E2", got[0])
			}
			if got[1].ExpId != "E1" {
				t.Errorf("QueryToProcessExperimentPlanModelOrderByStartTime() got = %v, want E1", got[1])
			}
		})
	}
}
