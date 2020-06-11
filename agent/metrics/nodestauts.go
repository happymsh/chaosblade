package metrics

import (
	"fmt"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type NodeStatusModel struct {
	// sample time , always systime
	SampleTime int64 `json:"sample_time"`
	// sample Interval
	SampleInterval int64 `json:"sample_interval"`
	// system load in last 1 min.
	Load1 string `json:"load1"`
	// system load in last 5 min.
	Load5 string `json:"load5"`
	// system load in last 15 min.
	Load15 string `json:"load15"`
	// Total CPU seconds from system startup
	CpuSecondsTotal string `json:"cpu_seconds_total"`
	// Idle CPU seconds from system startup
	CpuSecondsIdle string `json:"cpu_seconds_idle"`
	// last sample for CpuSecondsTotal
	CpuSecondsTotalLS string `json:"cpu_seconds_total_ls"`
	// last sample for CpuSecondsIdle
	CpuSecondsIdleLS string `json:"cpu_seconds_idle_ls"`
	// 1-(CpuSecondsIdle-CpuSecondsIdleS)/(CpuSecondsTotal-CpuSecondsTotalS)
	CpuUsage string `json:"cpu_usage"`
	// Total Memory
	MemoryMemtotalBytes string `json:"memory_memtotal_bytes"`
	// Buffer Memory
	MemoryMemBuffersBytes string `json:"memory_mem_buffers_bytes"`
	// Cached Memory
	MemoryCachedBytes string `json:"memory_cached_bytes"`
	// Free Memory
	MemoryMemFreeBytes string `json:"memory_mem_free_bytes"`
	// 1-(MemoryMemBuffersBytes+MemoryCachedBytes+MemoryMemFreeBytes)/MemoryMemtotalBytes
	MemoryUsage string `json:"memory_usage"`
	// Total Filesystem
	FilesystemSizeBytes string `json:"filesystem_size_bytes"`
	// Free Filesystem
	FilesystemFreeBytes string `json:"filesystem_free_bytes"`
	// 1-FilesystemFreeBytes/FilesystemSizeBytes
	FilesystemUsage string `json:"filesystem_usage"`
	// Open file descriptor
	FilefdAllocated string `json:"filefd_allocated"`
	// Maximun of open file descriptor
	FilefdMaximun string `json:"filefd_maximun"`
	// TCPs
	SockstatTCPAlloc string `json:"sockstat_tcp_alloc"`
	// TCPs of time_wait
	SockstatTCPTw string `json:"sockstat_tcp_tw"`
	// Network receive bytes from system startup
	NetworkReceiveBytesTotal string `json:"network_receive_bytes_total"`
	// Network transmit bytes from system startup
	NetworkTransmitBytesTotal string `json:"network_transmit_bytes_total"`
	// last sample for NetworkReceiveBytesTotal
	NetworkReceiveBytesTotalLS string `json:"network_receive_bytes_total_ls"`
	// last sample for NetworkTransmitBytesTotal
	NetworkTransmitBytesTotalLS string `json:"network_transmit_bytes_total_ls"`
	// Network receive bytes from last sample
	NetworkReceiveBytes string `json:"network_receive_bytes"`
	// Network transmit bytes from last sample
	NetworkTransmitBytes string `json:"network_transmit_bytes"`
	// Disk read bytes from system startup
	DiskReadBytesTotal string `json:"disk_read_bytes_total"`
	// Disk write bytes from system startup
	DiskWriteBytesTotal string `json:"disk_write_bytes_total"`
	// last smple for DiskReadBytesTotal
	DiskReadBytesTotalLS string `json:"disk_read_bytes_total_ls"`
	// last smple for DiskWriteBytesTotal
	DiskWriteBytesTotalLS string `json:"disk_write_bytes_total_ls"`
	// Disk Read bytes/s
	DiskReadBytes string `json:"disk_read_bytes"`
	// Disk Write bytes/s
	DiskWriteBytes string `json:"disk_write_bytes"`
}

var once sync.Once
var file *os.File
var ReportFileName string

const ReportDir = "report"
const ReportFile = "node_report.csv"
const ReportFileBak = "node_report.csv.bak"

var GBLength int64 = 1024 * 1024 * 1024

func OpenReportFile() {
	dirName := path.Join(util.GetProgramPath(), ReportDir)
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		logrus.Error("NodeStatusModel.Record, create report dir ", dirName, err.Error())
		return
	}
	ReportFileName = path.Join(dirName, ReportFile)
	file, err = os.OpenFile(ReportFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logrus.Error("NodeStatusModel.OpenReportFile, open file node.report err:", err.Error())
		return
	}
}

func CloseReportFile() {
	if file != nil {
		_ = file.Close()
	}
}

func (ns *NodeStatusModel) RecordTitle() {
	var title string
	title = fmt.Sprintf("%s,%s", title, "sample_time_for_human")
	t := reflect.TypeOf(*ns)
	for i := 0; i < t.NumField(); i++ {
		title = fmt.Sprintf("%s,%s", title, t.Field(i).Tag.Get("json"))
	}
	_, err := fmt.Fprintf(file, "%s\n", title[1:])
	if err != nil {
		logrus.Error("NodeStatusModel.RecordTitle, Fprintf file ", err.Error())
	}
}

func (ns *NodeStatusModel) Record() {
	//if file >=1GB, rotate file
	dirName := path.Join(util.GetProgramPath(), ReportDir)
	fileInfos, err := ioutil.ReadDir(dirName)
	if err != nil {
		logrus.Error("NodeStatusModel.Record, ReadDir, ", err.Error())
		return
	}
	for _, v := range fileInfos {
		if v.Name() == ReportFile {
			if v.Size() > GBLength {
				CloseReportFile()
				bakFileName := path.Join(dirName, ReportFileBak)
				err := os.Remove(bakFileName)
				if err != nil {
					logrus.Error("NodeStatusModel.Record, Remove file, ", err.Error())
					//return
				}
				err = os.Rename(ReportFileName, bakFileName)
				if err != nil {
					logrus.Error("NodeStatusModel.Record, Rename file, ", err.Error())
					return
				}
				file, err = os.OpenFile(ReportFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
				if err != nil {
					logrus.Error("NodeStatusModel.Record, open file node.report err:", err.Error())
					return
				}
			}
		}
	}

	//format line
	var record string
	t := time.Unix(ns.SampleTime/1000, (ns.SampleTime%1000)*int64(time.Millisecond))
	record = fmt.Sprintf("%s,%s", record, t.String())
	v := reflect.ValueOf(*ns)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Int64:
			record = fmt.Sprintf("%s,%s", record, strconv.FormatInt(field.Int(), 10))
		case reflect.String:
			record = fmt.Sprintf("%s,%s", record, field.String())
		}
	}

	//write file
	_, err = fmt.Fprintf(file, "%s\n", record[1:])
	if err != nil {
		logrus.Error("NodeStatusModel.Record, Fprintf file ", err.Error())
	}
}
