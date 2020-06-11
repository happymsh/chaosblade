package metrics

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chaosblade-io/chaosblade/agent/pdata"
)

var request *http.Request
var experimentPlanModel *pdata.ExperimentPlanModel = &pdata.ExperimentPlanModel{}
var nodeStatusModel *NodeStatusModel = &NodeStatusModel{}
var Metrics *MetricWriter = &MetricWriter{ExperimentPlanModel: experimentPlanModel, NodeStatusModel: nodeStatusModel}

func init() {
	request, _ = http.NewRequest("GET", "http://localhost:9100/this_is_a_fake_url", strings.NewReader(""))
	Metrics.Channel = make(chan int, 2)
	Metrics.Channel <- 1
}

//mock HttpServe by implemet ResponseWriter interface
type MetricWriter struct {
	*pdata.ExperimentPlanModel
	*NodeStatusModel
	CpuTotal         float64
	CpuIdle          float64
	CpuTotalLS       float64
	CpuIdleLS        float64
	FsTotal          float64
	FsFree           float64
	ReceiveTotal     float64
	TransmitTotal    float64
	ReceiveTotalLS   float64
	TransmitTotalLS  float64
	DiskReadTotal    float64
	DiskReadTotalLS  float64
	DiskWriteTotal   float64
	DiskWriteTotalLS float64
	Channel          chan int
	sync.RWMutex
}

func (w *MetricWriter) Header() http.Header {
	h := make(http.Header)
	return h
}

func (w *MetricWriter) WriteHeader(int) {
}

func (w *MetricWriter) Write(metricBytes []byte) (int, error) {
	buffer := bytes.NewBuffer(metricBytes)
	for {
		metricLine, _ := buffer.ReadBytes('\n')
		if len(metricLine) == 0 {
			break
		}
		metric := strings.Split(string(metricLine)[:len(metricLine)-1], " ")
		if len(metric) == 1 {
			break
		}
		value := metric[len(metric)-1]
		//cpuSecondsTotal
		if strings.HasPrefix(metric[0], "#") {
			//comment line
		} else if metric[0] == "node_load1" {
			w.Load1 = value
		} else if metric[0] == "node_load5" {
			w.Load5 = value
		} else if metric[0] == "node_load15" {
			w.Load15 = value
		} else if strings.HasPrefix(metric[0], "node_cpu_seconds_total{cpu=\"") {
			if cpuFloat, err := strconv.ParseFloat(value, 64); err == nil {
				w.CpuTotal += cpuFloat
				if strings.Contains(metric[0], "mode=\"idle\"") {
					w.CpuIdle += cpuFloat
				}
			}
		} else if metric[0] == "node_memory_Buffers_bytes" {
			w.MemoryMemBuffersBytes = value
		} else if metric[0] == "node_memory_Cached_bytes" {
			w.MemoryCachedBytes = value
		} else if metric[0] == "node_memory_MemFree_bytes" {
			w.MemoryMemFreeBytes = value
		} else if metric[0] == "node_memory_MemTotal_bytes" {
			w.MemoryMemtotalBytes = value
		} else if metric[0] == "node_memory_free_bytes" { // for Darwin ,see https://github.com/prometheus/node_exporter/issues/869
			w.MemoryMemFreeBytes = value
		} else if metric[0] == "node_memory_inactive_bytes" { // for Darwin
			// inactive can be recovered by Darwin
			w.MemoryCachedBytes = value
			w.MemoryMemBuffersBytes = "0"
		} else if metric[0] == "node_memory_total_bytes" { // for Darwin
			w.MemoryMemtotalBytes = value
		} else if strings.HasPrefix(metric[0], "node_filesystem_free_bytes") {
			bytesFloat, err := strconv.ParseFloat(value, 64)
			if err == nil {
				w.FsFree += bytesFloat
			}
		} else if strings.HasPrefix(metric[0], "node_filesystem_size_bytes") {
			bytesFloat, err := strconv.ParseFloat(value, 64)
			if err == nil {
				w.FsTotal += bytesFloat
			}
		} else if metric[0] == "node_filefd_allocated" {
			v, err := strconv.ParseFloat(value, 64)
			if err == nil {
				w.FilefdAllocated = strconv.FormatFloat(v, 'f', 0, 64)
			}
		} else if metric[0] == "node_filefd_maximum" {
			v, err := strconv.ParseFloat(value, 64)
			if err == nil {
				w.FilefdMaximun = strconv.FormatFloat(v, 'f', 0, 64)
			}
		} else if metric[0] == "node_sockstat_TCP_alloc" {
			v, err := strconv.ParseFloat(value, 64)
			if err == nil {
				w.SockstatTCPAlloc = strconv.FormatFloat(v, 'f', 0, 64)
			}
		} else if metric[0] == "node_sockstat_TCP_tw" {
			v, err := strconv.ParseFloat(value, 64)
			if err == nil {
				w.SockstatTCPTw = strconv.FormatFloat(v, 'f', 0, 64)
			}
		} else if strings.HasPrefix(metric[0], "node_network_receive_bytes_total") {
			bytesFloat, err := strconv.ParseFloat(value, 64)
			if err == nil {
				w.ReceiveTotal += bytesFloat
			}
		} else if strings.HasPrefix(metric[0], "node_network_transmit_bytes_total") {
			bytesFloat, err := strconv.ParseFloat(value, 64)
			if err == nil {
				w.TransmitTotal += bytesFloat
			}
		} else if strings.HasPrefix(metric[0], "node_disk_read_bytes_total") {
			bytesFloat, err := strconv.ParseFloat(value, 64)
			if err == nil {
				w.DiskReadTotal += bytesFloat
			}
		} else if strings.HasPrefix(metric[0], "node_disk_written_bytes_total") {
			bytesFloat, err := strconv.ParseFloat(value, 64)
			if err == nil {
				w.DiskWriteTotal += bytesFloat
			}
		}
	}
	return 0, nil
}

func (w *MetricWriter) Parse() {
	if w.CpuTotal != 0 {
		w.CpuSecondsTotal = strconv.FormatFloat(w.CpuTotal, 'f', 6, 64)
	}
	if w.CpuIdle != 0 {
		w.CpuSecondsIdle = strconv.FormatFloat(w.CpuIdle, 'f', 6, 64)
	}
	//CpuUsage: 1-(CpuSecondsIdle-CpuSecondsIdleS)/(CpuSecondsTotal-CpuSecondsTotalS)
	if total := w.CpuTotal - w.CpuTotalLS; total != 0 {
		w.CpuUsage = strconv.FormatFloat(1-(w.CpuIdle-w.CpuIdleLS)/total, 'f', 6, 64)
	}
	//MemoryUsage
	// 1-(MemoryMemBuffersBytes+MemoryCachedBytes+MemoryMemFreeBytes)/MemoryMemtotalBytes
	buffer, err1 := strconv.ParseFloat(w.MemoryMemBuffersBytes, 64)
	cache, err2 := strconv.ParseFloat(w.MemoryCachedBytes, 64)
	free, err3 := strconv.ParseFloat(w.MemoryMemFreeBytes, 64)
	mtotal, err4 := strconv.ParseFloat(w.MemoryMemtotalBytes, 64)
	if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
		w.MemoryUsage = strconv.FormatFloat(1-free/mtotal, 'f', 6, 64)
	}
	w.MemoryMemBuffersBytes = strconv.FormatFloat(buffer, 'f', 0, 64)
	w.MemoryCachedBytes = strconv.FormatFloat(cache, 'f', 0, 64)
	w.MemoryMemFreeBytes = strconv.FormatFloat(free, 'f', 0, 64)
	w.MemoryMemtotalBytes = strconv.FormatFloat(mtotal, 'f', 0, 64)

	if w.FsTotal != 0 {
		w.FilesystemSizeBytes = strconv.FormatFloat(w.FsTotal, 'f', 0, 64)
		w.FilesystemFreeBytes = strconv.FormatFloat(w.FsFree, 'f', 0, 64)
		w.FilesystemUsage = strconv.FormatFloat(1-w.FsFree/w.FsTotal, 'f', 6, 64)
	}

	w.NetworkReceiveBytesTotal = strconv.FormatFloat(w.ReceiveTotal, 'f', 0, 64)
	w.NetworkReceiveBytes = strconv.FormatFloat(w.ReceiveTotal-w.ReceiveTotalLS, 'f', 0, 64)

	w.NetworkTransmitBytesTotal = strconv.FormatFloat(w.TransmitTotal, 'f', 0, 64)
	w.NetworkTransmitBytes = strconv.FormatFloat(w.TransmitTotal-w.TransmitTotalLS, 'f', 0, 64)

	w.DiskReadBytesTotal = strconv.FormatFloat(w.DiskReadTotal, 'f', 0, 64)
	w.DiskReadBytes = strconv.FormatFloat(w.DiskReadTotal-w.DiskReadTotalLS, 'f', 0, 64)

	w.DiskWriteBytesTotal = strconv.FormatFloat(w.DiskWriteTotal, 'f', 0, 64)
	w.DiskWriteBytes = strconv.FormatFloat(w.DiskWriteTotal-w.DiskWriteTotalLS, 'f', 0, 64)

	w.SampleTime = time.Now().UnixNano() / int64(time.Millisecond)
}

func (w *MetricWriter) GetMetrics(interval int64) {
	w.Lock()
	_ = <-w.Channel
	defer w.Unlock()
	w.SampleInterval = interval
	w.PrePars()
	MetricsHandler.ServeHTTP(w, request)
	w.Parse()
	w.Channel <- 1
}

func (w *MetricWriter) PrePars() {
	w.CpuSecondsIdleLS = w.CpuSecondsIdle
	w.CpuSecondsTotalLS = w.CpuSecondsTotal
	w.CpuTotalLS = w.CpuTotal
	w.CpuIdleLS = w.CpuIdle
	w.NetworkReceiveBytesTotalLS = w.NetworkReceiveBytesTotal
	w.NetworkTransmitBytesTotalLS = w.NetworkTransmitBytesTotal
	w.TransmitTotalLS = w.TransmitTotal
	w.ReceiveTotalLS = w.ReceiveTotal
	w.DiskReadBytesTotalLS = w.DiskReadBytesTotal
	w.DiskWriteBytesTotalLS = w.DiskWriteBytesTotal
	w.DiskReadTotalLS = w.DiskReadTotal
	w.DiskWriteTotalLS = w.DiskWriteTotal
	w.CpuTotal = 0
	w.CpuIdle = 0
	w.FsTotal = 0
	w.FsFree = 0
	w.ReceiveTotal = 0
	w.TransmitTotal = 0
	w.DiskReadTotal = 0
	w.DiskWriteTotal = 0
}
