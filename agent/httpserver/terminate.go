package httpserver

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/chaosblade-io/chaosblade/agent/task"
	"net/http"
)

func terminalHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("terminate url is called, now to kill all chaos experiment")
	response := task.ExpTaskArrange(ctx_, true)
	_, err := w.Write([]byte(response.ToString()))
	if err != nil {
		logrus.Error("terminalHandler write error:", err.Error())
	}
	viper.Set("FLOW_TERMINATE", true)
	err = viper.WriteConfig()
	if err != nil {
		logrus.Error(err)
	}
}
