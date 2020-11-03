package hzcp

import (
	"context"
	"fmt"

	"path"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
)

type threadfullActionCommand struct {
	spec.BaseExpActionCommandSpec
}

func (*threadfullActionCommand) Name() string {
	return "threadfull"
}

func (*threadfullActionCommand) Aliases() []string {
	return []string{"tf"}
}

func (*threadfullActionCommand) ShortDesc() string {
	return "hzcp-jvm threadfull"
}

func (*threadfullActionCommand) LongDesc() string {
	return "hzcp-jvm threadfull"
}

func (*threadfullActionCommand) Matchers() []spec.ExpFlagSpec {
	return []spec.ExpFlagSpec{}
}

func (*threadfullActionCommand) Flags() []spec.ExpFlagSpec {
	return []spec.ExpFlagSpec{}
}

type jvmThreadfullExecutor struct {
	channel spec.Channel
}

func (ce *jvmThreadfullExecutor) Name() string {
	return "threadfull"
}

func (ce *jvmThreadfullExecutor) SetChannel(channel spec.Channel) {
	ce.channel = channel
}

func (ce *jvmThreadfullExecutor) Exec(uid string, ctx context.Context, model *spec.ExpModel) *spec.Response {
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
	//var jvmPercent int

	//jvmPercentStr := model.ActionFlags["jvm-percent"]
	//if jvmPercentStr != "" {
	//	var err error
	//	jvmPercent, err = strconv.Atoi(jvmPercentStr)
	//	if err != nil {
	//		return spec.ReturnFail(spec.Code[spec.IllegalParameters],
	//			"--jvm-percent value must be a positive integer")
	//	}
	//	if jvmPercent > 100 || jvmPercent < 0 {
	//		return spec.ReturnFail(spec.Code[spec.IllegalParameters],
	//			"--jvm-percent value must be a prositive integer and not bigger than 100")
	//	}
	//} else {
	//	jvmPercent = 100
	//}

	return ce.start(ctx)
}

// start burn jvm
func (ce *jvmThreadfullExecutor) start(ctx context.Context) *spec.Response {
	args := fmt.Sprintf("--start --debug=%t", util.Debug)
	return ce.channel.Run(ctx, path.Join(ce.channel.GetScriptPath(), hzcpcpufullload), args)
}

// stop burn jvm
func (ce *jvmThreadfullExecutor) stop(ctx context.Context) *spec.Response {
	return ce.channel.Run(ctx, path.Join(ce.channel.GetScriptPath(), hzcpcpufullload),
		fmt.Sprintf("--stop --debug=%t", util.Debug))
}
