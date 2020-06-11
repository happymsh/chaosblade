//plan data
package pdata

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type ExperimentPlanModel struct {
	ExpId       string `json:"event_id"`
	JobId       string
	NodetimeGap int64
	Uid         string
	Command     string `json:"command"`
	SubCommand  string `json:"sub_command"`
	Flag        string `json:"flag"`
	Status      string //0-ToBeProcessed 1-Processing 2-SuccessfulProcessing 3-ProcessingFailure
	Error       string
	StartOffset int64     `json:"start_offset"` //Millisecond
	EndOffset   int64     `json:"end_offset"`   //Millisecond
	StartTime   time.Time //Exp StartTime(node)
	EndTime     time.Time //Exp EndTime(node)
	CreateTime  time.Time
	UpdateTime  time.Time
}

const (
	ToBeProcessed        = "0"
	Processing           = "1"
	SuccessfulProcessing = "2"
	ProcessingFailure    = "3"
	Terminate            = "4"
)

//type ExperimentPlanSource interface {
//	// CheckAndInitExperimentPlanTable, if experimentPlan table not exists, then init it
//	CheckAndInitExperimentPlanTable()
//
//	// ExperimentPlanTableExists return true if experimentPlan exists
//	ExperimentPlanTableExists() (bool, error)
//
//	// InitExperimentPlanTable for first executed
//	InitExperimentPlanTable() error
//
//	// InsertExperimentPlanModel for creating chaos experiment_plan
//	//InsertExperimentPlanModel(model *ExperimentPlanModel) error
//
//	// UpdateExperimentPlanModelStatusByExpId
//	UpdateExperimentPlanModelStatusByExpId(expId, status, errMsg string) error
//
//	QueryExperimentPlanModelByExpId(expId string) (*ExperimentPlanModel, error)
//
//	QueryExistExperimentPlanModelByStatus(status string) (bool, error)
//}

const expPlanTableDDL = `CREATE TABLE IF NOT EXISTS experiment_plan (
	exp_id VARCHAR(32) PRIMARY KEY,
	job_id VARCHAR(20),
	nodetime_gap integer,
	uid VARCHAR(32),
	command VARCHAR NOT NULL,
	sub_command VARCHAR,
	flag VARCHAR,
	status VARCHAR NOT NULL,
	error VARCHAR,
	start_offset integer,
	end_offset integer ,
	start_time DATE NOT NULL,
	end_time DATE NOT NULL,
	create_time DATE NOT NULL,
	update_time DATE NOT NULL
)`

var expPlanIndexDDL = []string{
	`CREATE INDEX experiment_plan_command_idx ON experiment_plan (command)`,
	`CREATE INDEX experiment_plan_status_idx ON experiment_plan (status)`,
	`CREATE INDEX experiment_plan_uid_idx ON experiment_plan (uid)`,
	`CREATE INDEX experiment_plan_jobid_idx ON experiment_plan (job_id)`,
	`CREATE INDEX experiment_plan_expid_idx ON experiment_plan (exp_id)`,
}

var insertExpPlanDML = `INSERT INTO 
	experiment_plan (exp_id, job_id, nodetime_gap, uid, command, sub_command, flag, status, error, start_offset, end_offset, start_time, end_time, create_time, update_time)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func (s *PSource) CheckAndInitExperimentPlanTable() {
	exists, err := s.ExperimentPlanTableExists()
	if err != nil {
		logrus.Error(err.Error())
	}
	if !exists {
		err = s.InitExperimentPlanTable()
		if err != nil {
			logrus.Error(err.Error())
		}
	}
}

func (s *PSource) ExperimentPlanTableExists() (bool, error) {
	stmt, err := s.DB.Prepare(tableExistsDQL)
	if err != nil {
		return false, fmt.Errorf("select experiment_plan table exists err when invoke db prepare, %s", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query("experiment_plan")
	if err != nil {
		return false, fmt.Errorf("select experiment_plan table exists or not err, %s", err)
	}
	defer rows.Close()
	var c int
	for rows.Next() {
		rows.Scan(&c)
		break
	}
	return c != 0, nil
}

func (s *PSource) InitExperimentPlanTable() error {
	_, err := s.DB.Exec(expPlanTableDDL)
	if err != nil {
		return fmt.Errorf("create experiment_plan table err, %s", err)
	}
	for _, sqlStatement := range expPlanIndexDDL {
		_, err = s.DB.Exec(sqlStatement)
		if err != nil {
			return fmt.Errorf("create experiment_plan Index err, %s", err)
		}
	}
	return nil
}

func InsertExperimentPlanModel(tx *sql.Tx, model *ExperimentPlanModel) error {
	stmt, err := tx.Prepare(insertExpPlanDML)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		model.ExpId,
		model.JobId,
		model.NodetimeGap,
		model.Uid,
		model.Command,
		model.SubCommand,
		model.Flag,
		model.Status,
		model.Error,
		model.StartOffset,
		model.EndOffset,
		model.StartTime,
		model.EndTime,
		model.CreateTime,
		model.UpdateTime,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PSource) UpdateExperimentPlanModelStatusByExpId(expId, uid, status, errMsg string) error {
	stmt, err := s.DB.Prepare(`UPDATE experiment_plan
	SET uid=? ,status = ?, error = ?, update_time = ?
	WHERE exp_id = ?
`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(uid, status, errMsg, time.Now(), expId)
	if err != nil {
		return err
	}
	return nil
}

func (s *PSource) QueryExperimentPlanModelByExpId(expId string) (*ExperimentPlanModel, error) {
	stmt, err := s.DB.Prepare(`SELECT * FROM experiment_plan WHERE exp_id = ?`)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(expId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	models, err := getExperimentPlanModelsFrom(rows)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, nil
	}
	return models[0], nil
}

func getExperimentPlanModelsFrom(rows *sql.Rows) ([]*ExperimentPlanModel, error) {
	models := make([]*ExperimentPlanModel, 0)
	for rows.Next() {
		var nodetimeGap, startOffset, endOffset int64
		var expId, jobId, uid, command, subCommand, flag, status, error string
		var startTime, endTime, createTime, updateTime time.Time
		err := rows.Scan(&expId, &jobId, &nodetimeGap, &uid, &command, &subCommand, &flag, &status, &error, &startOffset, &endOffset, &startTime, &endTime, &createTime, &updateTime)
		if err != nil {
			return nil, err
		}
		model := &ExperimentPlanModel{
			ExpId:       expId,
			Uid:         uid,
			JobId:       jobId,
			NodetimeGap: nodetimeGap,
			Command:     command,
			SubCommand:  subCommand,
			Flag:        flag,
			Status:      status,
			Error:       error,
			StartOffset: startOffset,
			EndOffset:   endOffset,
			StartTime:   startTime,
			EndTime:     endTime,
			CreateTime:  createTime,
			UpdateTime:  updateTime,
		}
		models = append(models, model)
	}
	return models, nil
}

func (s *PSource) QueryExistExperimentPlanModelByStatus(status string) (bool, error) {
	stmt, err := s.DB.Prepare(`SELECT * FROM experiment_plan WHERE status = ?`)
	if err != nil {
		return false, err
	}
	rows, err := stmt.Query(status)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	models, err := getExperimentPlanModelsFrom(rows)
	if err != nil {
		return false, err
	}
	if len(models) == 0 {
		return false, nil
	}
	return true, nil
}

func (s *PSource) QueryExperimentPlanModelByStatus(status string) ([]*ExperimentPlanModel, error) {
	stmt, err := s.DB.Prepare(`SELECT * FROM experiment_plan WHERE Status = ?`)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	models, err := getExperimentPlanModelsFrom(rows)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, nil
	}
	return models, nil
}

func (s *PSource) QueryToProcessExperimentPlanModelOrderByStartTime() ([]*ExperimentPlanModel, error) {
	//stmt, err := s.DB.Prepare(`SELECT * FROM experiment_plan WHERE Status in (?)  and start_time < ? order by start_time`)
	stmt, err := s.DB.Prepare(`SELECT * FROM experiment_plan WHERE Status in ('0','1') and start_time <= datetime('now','localtime') order by start_time `)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	models, err := getExperimentPlanModelsFrom(rows)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, nil
	}
	return models, nil
}

func (s *PSource) QueryAllToProcessExperimentPlanModel() ([]*ExperimentPlanModel, error) {
	//stmt, err := s.DB.Prepare(`SELECT * FROM experiment_plan WHERE Status in (?)  and start_time < ? order by start_time`)
	stmt, err := s.DB.Prepare(`SELECT * FROM experiment_plan WHERE Status in ('0','1') `)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	models, err := getExperimentPlanModelsFrom(rows)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, nil
	}
	return models, nil
}

func (s *PSource) QueryExperimentPlanModelByJobId(jobId string) ([]*ExperimentPlanModel, error) {
	stmt, err := s.DB.Prepare(`SELECT * FROM experiment_plan WHERE job_id = ?`)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(jobId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	models, err := getExperimentPlanModelsFrom(rows)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, nil
	}
	return models, nil
}
