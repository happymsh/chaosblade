package metrics

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMetricWriter_GetMetrics(t *testing.T) {
	runFunc = func(arg []string) { run([]string{""}) }
	nodeStatusModel := &NodeStatusModel{}
	w := MetricWriter{NodeStatusModel: nodeStatusModel}
	w.Channel = make(chan int, 2)
	w.Channel <- 1
	w.GetMetrics(1)
	fmt.Printf("%+v\n", w)
	fmt.Printf("%+v\n", w.NodeStatusModel)
}

func TestSscanf(t *testing.T) {
	s := "98374598"
	var f float64
	_, err := fmt.Sscanf(s, "%e", &f)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(f)
	}
}

func TestParseFloat(t *testing.T) {
	s := "9.8374598e+07"
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(f)
	}
}
