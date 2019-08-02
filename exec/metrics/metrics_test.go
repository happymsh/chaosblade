package metrics

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/chaosblade-io/chaosblade/data"
)

func TestMetricWriter_GetMetrics(t *testing.T) {
	runFunc = func(arg []string) { run([]string{""}) }
	nodeStatusModel := &data.NodeStatusModel{}
	w := MetricWriter{NodeStatusModel: nodeStatusModel}
	w.GetMetrics()
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
