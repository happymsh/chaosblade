package metrics

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"

	"github.com/chaosblade-io/chaosblade/data"
)

//mock HttpServe by implemet ResponseWriter interface
type MetricWriter struct {
	*data.NodeStatusModel
	CpuTotal      float64
	FsTotal       float64
	FsFree        float64
	ReceiveTotal  float64
	TransmitTotal float64
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
					w.CpuSecondsIdle = value
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
			w.FilefdAllocated = value
		} else if metric[0] == "node_filefd_maximum" {
			w.FilefdMaximun = value
		} else if metric[0] == "node_sockstat_TCP_alloc" {
			w.SockstatTCPAlloc = value
		} else if metric[0] == "node_sockstat_TCP_tw" {
			w.SockstatTCPTw = value
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
		}
	}
	return 0, nil
}

func (w *MetricWriter) Parse() {
	if w.CpuTotal != 0 {
		w.CpuSecondsTotal = strconv.FormatFloat(w.CpuTotal, 'f', 2, 64)
	}
	//CpuUsage
	//MemoryUsage
	if w.FsTotal != 0 {
		w.FilesystemSizeBytes = strconv.FormatFloat(w.FsTotal, 'f', 2, 64)
		w.FilesystemFreeBytes = strconv.FormatFloat(w.FsFree, 'f', 2, 64)
		w.FilesystemUsage = strconv.FormatFloat(1-w.FsFree/w.FsTotal, 'f', 2, 64)
	}
	if w.ReceiveTotal != 0 {
		w.NetworkReceiveBytesTotal = strconv.FormatFloat(w.ReceiveTotal, 'f', 2, 64)
	}
	if w.TransmitTotal != 0 {
		w.NetworkTransmitBytesTotal = strconv.FormatFloat(w.TransmitTotal, 'f', 2, 64)
	}
}

func (w *MetricWriter) GetMetrics() {
	request, _ := http.NewRequest("GET", "http://localhost:9100/fakeurl", strings.NewReader(""))
	MetricsHandler.ServeHTTP(w, request)
	w.Parse()
}
