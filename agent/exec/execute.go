package exec

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/chaosblade-io/chaosblade-spec-go/channel"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
)

var channel_ spec.Channel = channel.NewLocalChannel()

const BLADE = "blade"

var CRun = channel_.Run

type Result struct {
	Result string `json:"result"`
}

func ExecuteExp(ctx context.Context, command, subCommand, flags, timeout string) *spec.Response {
	args := fmt.Sprintf("create %s %s", command, subCommand)
	for _, flag := range strings.Fields(flags) {
		//args = fmt.Sprintf("%s --%s", args, flag)
		args = fmt.Sprintf("%s %s", args, flag)
	}
	if timeout != "" {
		args = fmt.Sprintf("%s --timeout=%s", args, timeout)
	}
	return CRun(ctx, path.Join(util.GetProgramPath(), BLADE), args)
}
