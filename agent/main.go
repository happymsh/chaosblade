package main

import (
	"fmt"
	"github.com/chaosblade-io/chaosblade/agent/cli"
	"github.com/chaosblade-io/chaosblade/cli/cmd"
	"os"
)

func main() {
	baseCommand := cmd.CmdInit()
	//add agent command, must implement the cmd.Command interface of chaosblade
	baseCommand.AddCommand(&cli.AgentCommand{})
	baseCommand.AddCommand(&cli.AgentVersionCommand{})

	if err := baseCommand.CobraCmd().Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
