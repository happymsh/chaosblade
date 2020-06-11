package httpserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/chaosblade-io/chaosblade/agent/config"
	"github.com/chaosblade-io/chaosblade/agent/etcd"
	"github.com/chaosblade-io/chaosblade/agent/metrics"
	"net/http"
	"time"

	"github.com/chaosblade-io/chaosblade/agent/task"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
)

//scrapeInterval for metrics
var scrapeInterval int64
var ctx_ = context.Background()

const TaskKey = "TaskKey"

func notFoundHandler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte(spec.Return(spec.Code[spec.HandlerNotFound]).ToString()))
	if err != nil {
		logrus.Error("call notFoundHandler error:", err.Error())
	}
}

func AgentInit(port string, rtimeout, interval int64, testFlag bool, environment string, hostIp string, hostPort string, containerId string, containerPort string, appShortname string) error {
	logrus.Info("AgentInit parameters:")
	logrus.Info("[port]", port)
	logrus.Info("[rtimeout]", rtimeout)
	logrus.Info("[interval]", interval)
	logrus.Info("[testFlag]", testFlag)
	logrus.Info("[environment]", environment)
	logrus.Info("[hostIp]", hostIp)
	logrus.Info("[hostPort]", hostPort)
	logrus.Info("[containerId]", containerId)
	logrus.Info("[containerPort]", containerPort)
	logrus.Info("[appShortname]", appShortname)
	logrus.Info("[EtcdEndpoints]", config.EtcdEndpoints)
	logrus.Info("[EtcdBeatTime]", config.EtcdBeatTime)
	logrus.Info("[EtcdDialTimeout]", config.EtcdDialTimeout)
	logrus.Info("[EtcdDialKeepAliveTime]", config.EtcdDialKeepAliveTime)

	errCh := make(chan error)

	//do experiment task
	logrus.Info("do experiment task")
	ctx := context.WithValue(ctx_, TaskKey, "ExpTaskArrange")
	go func() {
		for {
			_ = task.ExpTaskArrange(ctx, false)
			time.Sleep(time.Second)
		}
	}()
	//do metrics
	//ctx = context.WithValue(ctx_, TaskKey, "MetricsCollect")
	logrus.Info("do metrics")
	go func() {
		metrics.OpenReportFile()
		defer metrics.CloseReportFile()
		metrics.Metrics.RecordTitle()
		for {
			metrics.Metrics.GetMetrics(interval)
			metrics.Metrics.Record()
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}()

	//do regist to etcd
	if !testFlag {
		logrus.Info("do regist to etcd")
		var key string
		if environment == config.VIRTUAL {
			if appShortname != "" {
				//chaos/virtual/appShortName/hostIp/hostPort
				key = fmt.Sprintf("/chaos/virtual/%s/%s/%s", appShortname, hostIp, hostPort)
			} else {
				key = fmt.Sprintf("/chaos/virtual/%s/%s", hostIp, hostPort)
			}
		} else if environment == config.CONTAINER {
			if appShortname != "" {
				//chaos/container/appShortName/hostIp/hostPort/containerId/containerPort
				key = fmt.Sprintf("/chaos/container/%s/%s/%s/%s/%s", appShortname, hostIp, hostPort, containerId, containerPort)
			} else {
				key = fmt.Sprintf("/chaos/container/%s/%s/%s/%s", hostIp, hostPort, containerId, containerPort)
			}

		} else {
			return errors.New("environment invalid")
		}
		ctx = context.WithValue(ctx_, TaskKey, "Available")
		logrus.Info("initial etcd client")
		err := etcd.ClientInit(ctx)
		if err != nil {
			return err
		}
		defer etcd.Client.Close()
		//listen Lease KeepAlive and save it to log
		logrus.Info("listen Lease KeepAlive and save it to log")
		go func() {
			for {
				select {
				case leaseKeepResp := <-etcd.LeaseRespChan:
					if leaseKeepResp == nil {
						logrus.Error("listen etcd Lease KeepAlive: release failed.")
					} else {
						logrus.Debug("listen etcd Lease KeepAlive: release success.")
					}
				}
				time.Sleep(time.Duration(config.EtcdBeatTime) * time.Second)
			}
		}()
		err = etcd.Available(ctx, key)
		if err != nil {
			return err
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", notFoundHandler)
	mux.HandleFunc("/chaosagent/api/exp/dispatch", experimentPlanHandler)
	mux.HandleFunc("/chaosagent/api/exp/terminate", terminalHandler)
	mux.HandleFunc("/chaosagent/api/metrics", statusHandler)
	mux.HandleFunc("/chaosagent/api/report", reportHandler)
	mux.Handle("/chaosagent/api/metrics/node_exporter", &metrics.MetricsHandler)
	srv := &http.Server{
		Addr:        fmt.Sprintf(":%s", port),
		Handler:     mux,
		ReadTimeout: time.Second * time.Duration(rtimeout),
		//WriteTimeout: time.Second * time.Duration(wtimeout),
	}

	//start httpserver
	logrus.Info("start httpserver")
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			errCh <- err
		}
	}()

	//handle http error
	logrus.Info("handle http error")
	select {
	case e := <-errCh:
		logrus.Error("agent error:", e.Error())
		return e
	}
}
