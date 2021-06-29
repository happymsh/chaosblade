/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hzcp

import (
	"context"
	"fmt"
	"path"

	"github.com/chaosblade-io/chaosblade-spec-go/channel"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
)

type CpuCommandModelSpec struct {
	spec.BaseExpModelCommandSpec
}

func NewCpuCommandModelSpec() spec.ExpModelCommandSpec {
	return &CpuCommandModelSpec{
		spec.BaseExpModelCommandSpec{
			ExpActions: []spec.ExpActionCommandSpec{
				&fullLoadActionCommand{
					spec.BaseExpActionCommandSpec{
						ActionMatchers: []spec.ExpFlagSpec{},
						ActionFlags:    []spec.ExpFlagSpec{},
						ActionExecutor: &cpuExecutor{},
					},
				},
			},
			ExpFlags: []spec.ExpFlagSpec{
				&spec.ExpFlag{
					Name:     "keyword",
					Desc:     "using for matching target Java process",
					Required: true,
				},
				&spec.ExpFlag{
					Name:     "port",
					Desc:     "using for attaching target Java process",
					Required: false,
				},
			},
		},
	}
}

func (*CpuCommandModelSpec) Name() string {
	return "hzcp-cpu"
}

func (*CpuCommandModelSpec) ShortDesc() string {
	return "hzcp-cpu experiment"
}

func (*CpuCommandModelSpec) LongDesc() string {
	return "hzcp-cpu experiment, for example full load"
}

func (*CpuCommandModelSpec) Example() string {
	return "blade create hzcp-cpu fullload "
}

type fullLoadActionCommand struct {
	spec.BaseExpActionCommandSpec
}

func (*fullLoadActionCommand) Name() string {
	return "fullload"
}

func (*fullLoadActionCommand) Aliases() []string {
	return []string{"fl"}
}

func (*fullLoadActionCommand) ShortDesc() string {
	return "hzcp-cpu load"
}

func (*fullLoadActionCommand) LongDesc() string {
	return "hzcp-cpu load"
}

func (*fullLoadActionCommand) Matchers() []spec.ExpFlagSpec {
	return []spec.ExpFlagSpec{}
}

func (*fullLoadActionCommand) Flags() []spec.ExpFlagSpec {
	return []spec.ExpFlagSpec{}
}

type cpuExecutor struct {
	channel spec.Channel
}

func (ce *cpuExecutor) Name() string {
	return "hzcp-cpu"
}

func (ce *cpuExecutor) SetChannel(channel spec.Channel) {
	ce.channel = channel
}

func (ce *cpuExecutor) Exec(uid string, ctx context.Context, model *spec.ExpModel) *spec.Response {
	err := checkCpuExpEnv()
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

const hzcpcpufullload = "hzcpcpufullload"

// start burn cpu
func (ce *cpuExecutor) start(ctx context.Context, keyWord string, port string) *spec.Response {
	args := fmt.Sprintf("--start --keyword=%s --port=%s --debug=%t", keyWord, port, util.Debug)
	return ce.channel.Run(ctx, path.Join(ce.channel.GetScriptPath(), hzcpcpufullload), args)
}

// stop burn cpu
func (ce *cpuExecutor) stop(ctx context.Context) *spec.Response {
	return ce.channel.Run(ctx, path.Join(ce.channel.GetScriptPath(), hzcpcpufullload),
		fmt.Sprintf("--stop --debug=%t", util.Debug))
}

func checkCpuExpEnv() error {
	commands := []string{"ps", "awk", "grep", "kill", "nohup", "tr"}
	for _, command := range commands {
		if !channel.IsCommandAvailable(command) {
			return fmt.Errorf("%s command not found", command)
		}
	}
	return nil
}