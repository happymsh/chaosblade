package cli

import "github.com/spf13/cobra"

// agentBaseCommand
type agentBaseCommand struct {
	command *cobra.Command
}

func (bc *agentBaseCommand) Init() {
}

func (bc *agentBaseCommand) CobraCmd() *cobra.Command {
	return bc.command
}

func (bc *agentBaseCommand) Name() string {
	return bc.command.Name()
}
