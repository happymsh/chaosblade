package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type NodeStatusModel struct {
	// system load in last 1 min.
	Load1 string
	// system load in last 5 min.
	Load5 string
	// system load in last 15 min.
	Load15 string
	// Total CPU seconds from system startup
	CpuSecondsTotal string
	// Idle CPU seconds from system startup
	CpuSecondsIdle string
	// last sample for CpuSecondsTotal
	CpuSecondsTotalLS string
	// last sample for CpuSecondsIdle
	CpuSecondsIdleLS string
	// 1-(CpuSecondsIdle-CpuSecondsIdleS)/(CpuSecondsTotal-CpuSecondsTotalS)
	CpuUsage string
	// Total Memory
	MemoryMemtotalBytes string
	// Buffer Memory
	MemoryMemBuffersBytes string
	// Cached Memory
	MemoryCachedBytes string
	// Free Memory
	MemoryMemFreeBytes string
	// 1-(MemoryMemBuffersBytes+MemoryCachedBytes+MemoryMemFreeBytes)/MemoryMemtotalBytes
	MemoryUsage string
	// Total Filesystem
	FilesystemSizeBytes string
	// Free Filesystem
	FilesystemFreeBytes string
	// 1-FilesystemFreeBytes/FilesystemSizeBytes
	FilesystemUsage string
	// Open file descriptor
	FilefdAllocated string
	// Maximun of open file descriptor
	FilefdMaximun string
	// TCPs
	SockstatTCPAlloc string
	// TCPs of time_wait
	SockstatTCPTw string
	// Network receive bytes from system startup
	NetworkReceiveBytesTotal string
	// Network transmit bytes from system startup
	NetworkTransmitBytesTotal string
	// last sample for NetworkReceiveBytesTotal
	NetworkReceiveBytesTotalLS string
	// last sample for NetworkTransmitBytesTotal
	NetworkTransmitBytesTotalLS string
	// if this record is queried, the QueryFlag will be set
	QueryFlag string
	// sample time , always systime
	SampleTime time.Time
	// set when QueryFlag changed
	UpdateTime time.Time
}

type NodeStatusSource interface {
	// CheckAndInitNodeStatusTable, if NodeStatus table not exists, then init it
	CheckAndInitNodeStatusTable()

	// NodeStatusTableExists return true if NodeStatus exists
	NodeStatusTableExists() (bool, error)

	// InitNodeStatusTable for first executed
	InitNodeStatusTable() error

	// InsertNodeStatusModel for creating chaos experiment_plan
	InsertNodeStatusModel(model *NodeStatusModel) error

	// UpdateNodeStatusModelByUid
	UpdateNodeStatusModelByExpId(expId, status, errMsg string) error
}

const nodeStatusTableDDL = `CREATE TABLE IF NOT EXISTS node_status (
	exp_id VARCHAR(20) PRIMARY KEY,
	flow_id VARCHAR(20) NOT NULL,
	task_id VARCHAR(20) NOT NULL,

	Load1                     string
	Load5                     string
	Load15                    string
	CpuTotal           string
	CpuSecondsIdle            string
	CpuUsage                  string
	MemoryMemtotalBytes       string
	MemoryMemBuffersBytes     string
	MemoryCachedBytes         string
	MemoryMemFreeBytes        string
	MemoryUsage               string
	FilesystemSizeBytes       string
	FilesystemFreeBytes       string
	FilesystemUsage           string
	FilefdAllocated           string
	FilefdMaximun             string
	SockstatTCPAlloc          string
	SockstatTCPTw             string
	NetworkReceiveBytesTotal  string
	NetworkTransmitBytesTotal string
	QueryFlag                 string
	SampleTime                time.Time
	UpdateTime                time.Time
)`

var nodeStatusIndexDDL = []string{
	`CREATE INDEX expplan_command_idx ON experiment_plan (command)`,
	`CREATE INDEX expplan_status_idx ON experiment_plan (status)`,
}

var insertNodeStatusDML = `INSERT INTO 
	node_status (exp_id, flow_id, task_id, command, sub_command, flag, status, error, start_time, end_time, create_time, update_time)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func (s *Source) CheckAndInitNodeStatusTable() {
	exists, err := s.NodeStatusTableExists()
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	if !exists {
		err = s.InitNodeStatusTable()
		if err != nil {
			logrus.Fatalf(err.Error())
		}
	}
}

func (s *Source) NodeStatusTableExists() (bool, error) {
	stmt, err := s.DB.Prepare(tableExistsDQL)
	if err != nil {
		return false, fmt.Errorf("select node_status table exists err when invoke db prepare, %s", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query("node_status")
	if err != nil {
		return false, fmt.Errorf("select node_status table exists or not err, %s", err)
	}
	defer rows.Close()
	var c int
	for rows.Next() {
		rows.Scan(&c)
		break
	}
	return c != 0, nil
}

func (s *Source) InitNodeStatusTable() error {
	_, err := s.DB.Exec(nodeStatusTableDDL)
	if err != nil {
		return fmt.Errorf("create node_status table err, %s", err)
	}
	for _, sql := range nodeStatusIndexDDL {
		_, err = s.DB.Exec(sql)
		if err != nil {
			return fmt.Errorf("create node_status Index err, %s", err)
		}
	}
	return nil
}

func (s *Source) InsertNodeStatusModel(model *NodeStatusModel) error {
	stmt, err := s.DB.Prepare(insertNodeStatusDML)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		//model.ExpId,
		//model.JobId,
		//model.TaskId,
		//model.Command,
		//model.SubCommand,
		//model.Flag,
		//model.Status,
		//model.Error,
		//model.StartTime,
		//model.EndTime,
		model.SampleTime,
		model.UpdateTime,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Source) UpdateNodeStatusModelByExpId(expId, status, errMsg string) error {
	stmt, err := s.DB.Prepare(`UPDATE node_status
	SET status = ?, error = ?, update_time = ?
	WHERE exp_id = ?
`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(status, errMsg, time.Now(), expId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Source) QueryNodeStatusModelByExpId(expId string) (*NodeStatusModel, error) {
	stmt, err := s.DB.Prepare(`SELECT * FROM node_status WHERE exp_id = ?`)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(expId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	models, err := getNodeStatusModelsFrom(rows)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, nil
	}
	return models[0], nil
}

func getNodeStatusModelsFrom(rows *sql.Rows) ([]*NodeStatusModel, error) {
	models := make([]*NodeStatusModel, 0)
	for rows.Next() {
		var expId, flowId, taskId, command, subCommand, flag, status, error string
		var startTime, endTime, createTime, updateTime time.Time
		err := rows.Scan(&expId, &flowId, &taskId, &command, &subCommand, &flag, &status, &error, &startTime, &endTime, &createTime, &updateTime)
		if err != nil {
			return nil, err
		}
		model := &NodeStatusModel{
			//ExpId:      expId,
			//JobId:     flowId,
			//TaskId:     taskId,
			//Command:    command,
			//SubCommand: subCommand,
			//Flag:       flag,
			//Status:     status,
			//Error:      error,
			//StartTime:  startTime,
			//EndTime:    endTime,
			SampleTime: createTime,
			UpdateTime: updateTime,
		}
		models = append(models, model)
	}
	return models, nil
}
