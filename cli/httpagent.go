package main

import (
	"github.com/chaosblade-io/chaosblade/exec/httpagent"
	"github.com/spf13/cobra"
)

const httpAgent = "http_agent"

type HttpCommand struct {
	baseCommand
}

func (hc *HttpCommand) Init() {
	hc.command = &cobra.Command{
		Use:   "httpagent",
		Short: "start a http daemon server to receive chaos experiments.",
		Long:  "start a http daemon server to receive chaos experiments. the default port is 13500.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return httpagent.AgentInit("36661", 5, 5, 1)
			//var channel = exec.NewLocalChannel()
			//port := "8080"
			//nohupArgs := fmt.Sprintf(`%s %s &`, path.Join(channel.GetScriptPath(), httpAgent), port)
			//response := channel.Run(context.Background(), "nohup", nohupArgs)
			//if !response.Success {
			//	errors.New(response.Error())
			//}
		},
	}
}
