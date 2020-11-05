package hzcp

import (
	"context"
	"fmt"

	"path"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
)

type oomActionCommand struct {
	spec.BaseExpActionCommandSpec
}

func (*oomActionCommand) Name() string {
	return "oom"
}

func (*oomActionCommand) Aliases() []string {
	return []string{"oom"}
}

func (*oomActionCommand) ShortDesc() string {
	return "hzcp-jvm oom"
}

func (*oomActionCommand) LongDesc() string {
	return "hzcp-jvm oom"
}

func (*oomActionCommand) Matchers() []spec.ExpFlagSpec {
	return []spec.ExpFlagSpec{}
}

func (*oomActionCommand) Flags() []spec.ExpFlagSpec {
	return []spec.ExpFlagSpec{}
}

type jvmOomExecutor struct {
	channel spec.Channel
}

func (ce *jvmOomExecutor) Name() string {
	return "oom"
}

func (ce *jvmOomExecutor) SetChannel(channel spec.Channel) {
	ce.channel = channel
}

func (ce *jvmOomExecutor) Exec(uid string, ctx context.Context, model *spec.ExpModel) *spec.Response {
	err := checkJvmExpEnv()
	if err != nil {
		return spec.ReturnFail(spec.Code[spec.CommandNotFound], err.Error())
	}
	if ce.channel == nil {
		return spec.ReturnFail(spec.Code[spec.ServerError], "channel is nil")
	}
	if _, ok := spec.IsDestroy(ctx); ok {
		return ce.stop(ctx)
	}
	keyWord := model.ActionFlags["keyword"]
	port := model.ActionFlags["port"]

	return ce.start(ctx, keyWord, port)
}

const hzcpjvmoom = "hzcpjvmoom"

// start burn jvm
func (ce *jvmOomExecutor) start(ctx context.Context, keyWord string, port string) *spec.Response {
	args := fmt.Sprintf("--start --keyword=%s --port=%s --debug=%t", keyWord, port, util.Debug)
	return ce.channel.Run(ctx, path.Join(ce.channel.GetScriptPath(), hzcpjvmoom), args)
}

// stop burn jvm
func (ce *jvmOomExecutor) stop(ctx context.Context) *spec.Response {
	return ce.channel.Run(ctx, path.Join(ce.channel.GetScriptPath(), hzcpjvmoom),
		fmt.Sprintf("--stop --debug=%t", util.Debug))
}
