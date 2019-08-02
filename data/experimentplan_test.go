package data

import (
	"database/sql"
	"testing"
	"time"
)

const testdataFile = "chaosblade_test.dat"

func TestSource_ExperimentPlanTableExists(t *testing.T) {
	database, err := sql.Open("sqlite3", testdataFile)
	defer database.Close()
	if err != nil {
		t.Error(err)
		return
	}
	source := &Source{
		DB: database,
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
	source := &Source{
		DB: database,
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
	source := &Source{
		DB: database,
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
}

func TestSource_UpdateExperimentPlanModelByExpId(t *testing.T) {
	database, err := sql.Open("sqlite3", testdataFile)
	defer database.Close()
	if err != nil {
		t.Error(err)
		return
	}
	source := &Source{
		DB: database,
	}
	err = source.UpdateExperimentPlanModelByExpId("E3", "2", "cpuburn execute error.")
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
		t.Error("UpdateExperimentPlanModelByExpId not worked for status.")
	}
	if result.Error == "" {
		t.Error("UpdateExperimentPlanModelByExpId not worked for error.")
	}
}
