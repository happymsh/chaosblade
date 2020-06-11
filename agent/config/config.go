package config

import (
	"github.com/chaosblade-io/chaosblade-spec-go/util"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path"
	"strings"
	"time"
)

const conf = "conf/agent_conf.yml"

//for config
const (
	ENV_CHAOS_ETCD_ENDPOINT           = "ENV_CHAOS_ETCD_ENDPOINT"
	ENV_CHAOS_ETCD_BEAT_TIME          = "ENV_CHAOS_ETCD_BEAT_TIME"
	ENV_CHAOS_ETCD_DIAL_TIMEOUT       = "ENV_CHAOS_ETCD_DIAL_TIMEOUT"
	ENV_CHAOS_ETCD_DIAL_KEEPALIVETIME = "ENV_CHAOS_ETCD_DIAL_KEEPALIVETIME"
)

//for command flag
const (
	VIRTUAL   = "1" //虚拟机
	CONTAINER = "2" //容器

	PORT = "36661"
)

var EtcdEndpoints []string
var EtcdDialTimeout time.Duration
var EtcdDialKeepAliveTime time.Duration
var EtcdBeatTime int64

func varInit() {
	EtcdBeatTime = viper.GetInt64(ENV_CHAOS_ETCD_BEAT_TIME)
	EtcdDialTimeout = time.Duration(viper.GetInt64(ENV_CHAOS_ETCD_DIAL_TIMEOUT)) * time.Second
	EtcdDialKeepAliveTime = time.Duration(viper.GetInt64(ENV_CHAOS_ETCD_DIAL_KEEPALIVETIME)) * time.Second
	EtcdEndpoints = strings.Split(viper.GetString(ENV_CHAOS_ETCD_ENDPOINT), ",")
	if EtcdBeatTime == 0 {
		EtcdBeatTime = 5
	}
	if EtcdDialTimeout == time.Duration(0) {
		EtcdDialTimeout = time.Duration(30) * time.Second
	}
	if EtcdDialKeepAliveTime == time.Duration(0) {
		EtcdDialKeepAliveTime = time.Duration(10) * time.Second
	}
}

func confInit() {
	viper.SetConfigFile(path.Join(util.GetProgramPath(), conf))
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error("read config file error:", err.Error())
		panic(err)
	}
	varInit()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Info("config file changed:", e.Name)
		//varInit()
	})
}

var ConfInit = confInit
