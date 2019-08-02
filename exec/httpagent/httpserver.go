package httpagent

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chaosblade-io/chaosblade/transport"

	"github.com/chaosblade-io/chaosblade/exec/metrics"
)

//scrapeInterval for metrics
var scrapeInterval int64

func notFoundHandler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte(transport.Return(transport.Code[transport.HandlerNotFound]).ToString()))
	if err != nil {
		fmt.Sprintf("call notFoundHandler error. %s", err.Error())
	}
}

func todoHandler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte(`{"hello, this API will be provide soon."}`))
	if err != nil {
		fmt.Sprintf("call todoHandler error. %s", err.Error())
	}
}

func AgentInit(port string, rtimeout, wtimeout, interval int64) error {
	scrapeInterval = interval
	mux := http.NewServeMux()
	mux.HandleFunc("/", notFoundHandler)
	mux.HandleFunc("/chaosagent/api/exp/dispatch", experimentPlanHandler)
	mux.HandleFunc("/chaosagent/api/exp/terminate", todoHandler)
	mux.HandleFunc("/chaosagent/api/metrics", todoHandler)
	mux.HandleFunc("/chaosagent/api/metrics/node_exporter", todoHandler)
	mux.Handle("/node_exporter", &metrics.MetricsHandler)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  time.Millisecond * time.Duration(rtimeout),
		WriteTimeout: time.Millisecond * time.Duration(wtimeout),
	}
	return srv.ListenAndServe()
}
