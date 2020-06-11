package cli

import (
	"github.com/chaosblade-io/chaosblade/version"
	"github.com/spf13/cobra"
)

var AgtVer string = "unknow"

type AgentVersionCommand struct {
	agentBaseCommand
}

func (hc *AgentVersionCommand) Init() {
	hc.command = &cobra.Command{
		Use:   "icbc-version",
		Short: "Print ICBC agent version info",
		Long:  "Print ICBC agent version info",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("ICBC agent version: %s\n", AgtVer)
			cmd.Printf("chaosblade version: %s\n", version.Ver)
			cmd.Printf("env: %s\n", version.Env)
			cmd.Printf("build-time: %s\n", version.BuildTime)
			return
		},
	}
}
