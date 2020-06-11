package cli

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/chaosblade-io/chaosblade/agent/config"
	"github.com/chaosblade-io/chaosblade/agent/httpserver"
	"strconv"
)

const httpAgent = "http_agent"

type AgentCommand struct {
	agentBaseCommand
}

var testFlag bool          //agent测试标识，为ture时，不会校验传入的environment等信息，并且不会向etcd上报心跳
var environment *string    //1-虚拟机-VIRTUAL 2-容器-CONTAINER
var hostIp *string         //宿主机IP
var hostPort *string       //宿主机端口
var containerId *string    //environment=CONTAINER时，需要传递容器ID
var containerPort *string  //environment=CONTAINER时，需要传递容器内部端口
var sampleInterval *string //节点系统资源采样间隔，单位秒
var port string            //agent 启动端口
var appShortname *string   //应用简称

var EnvErr error = errors.New("environment must be 1(VIRTUAL) or 2(CONTAINER)")
var HostErr error = errors.New("hostIp and hostPort can not be empty")
var ContainerErr error = errors.New("containerId and containerPort can not be empty when environment is CONTAINER")

func (hc *AgentCommand) Init() {
	hc.command = &cobra.Command{
		Use:   "httpagent",
		Short: "start a http daemon server to receive chaos experiments",
		Long: `
start a http daemon server to receive chaos experiments. the default port is 36661.
1) TEST MODE: if environment is empty, then start agent with port by 36661 and sample-interval by 1s, other parameters will be ignore
2) VIRTUAL: parameters required(environment=1 && host-ip=xxx.xxx.xxx.xxx && host-port=xx)
3) CONTAINER: parameters required(environment=2 && host-ip=xxx.xxx.xxx.xxx && host-port=xx && container-id=xxx && container-port=xx)

Note:in TEST MODE, agent would not regist information(include environment,host-ip,host-port,container-id,container-port) to ETCD.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config.ConfInit()
			if *environment == "" {
				testFlag = true
			}
			if !testFlag {
				if *environment != config.VIRTUAL && *environment != config.CONTAINER {
					return EnvErr
				}
				if *hostIp == "" || *hostPort == "" {
					return HostErr
				}
				if (*containerId == "" || *containerPort == "") && *environment == config.CONTAINER {
					return ContainerErr
				}
				if *environment == config.CONTAINER {
					port = *containerPort
				} else {
					port = *hostPort
				}
			} else {
				port = config.PORT //for test mode default
			}

			interval, err := strconv.ParseInt(*sampleInterval, 10, 64)
			if err != nil {
				return errors.New("sampleInterval is not a valid seconds:" + err.Error())
			}
			if *appShortname == "" {
				//todo 临时去除校验
				//return errors.New("appShortname can not be empty")
			}
			//write config file.
			viper.Set("AGENT", map[string]interface{}{
				"interval":      interval,
				"testFlag":      testFlag,
				"environment":   *environment,
				"hostIp":        *hostIp,
				"hostPort":      *hostPort,
				"containerId":   *containerId,
				"containerPort": *containerPort,
				"appShortname":  *appShortname,
			})
			err = viper.WriteConfig()
			if err != nil {
				logrus.Error(err)
			}
			return httpserver.AgentInit(port, 60, interval, testFlag,
				*environment, *hostIp, *hostPort, *containerId, *containerPort, *appShortname)
		},
	}
	environment = hc.command.Flags().StringP(
		"environment", "e", "", "specify the environment(1-VIRTUAL,2-CONTAINER) where agent run.")
	hostIp = hc.command.Flags().StringP(
		"host-ip", "i", "", "specify the host IP of machine where agent run.")
	hostPort = hc.command.Flags().StringP(
		"host-port", "p", "", "specify the host port of machine where agent run.")
	containerId = hc.command.Flags().StringP(
		"container-id", "c", "", "specify the ID of container where agent run.")
	containerPort = hc.command.Flags().StringP(
		"container-port", "o", config.PORT, "specify the port of container where agent run.")
	sampleInterval = hc.command.Flags().StringP(
		"sample-interval", "s", "1", "specify the metrics sample interval.")
	appShortname = hc.command.Flags().StringP(
		"app-shortname", "a", "", "specify the application short name. eg F-APIP")
}
