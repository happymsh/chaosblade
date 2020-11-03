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

package exec

import (
	"context"
	"fmt"
	"path"
	"strconv"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
)

type BlockActionSpec struct {
	spec.BaseExpActionCommandSpec
}

func NewBlockActionSpec() spec.ExpActionCommandSpec {
	return &BlockActionSpec{
		spec.BaseExpActionCommandSpec{
			ActionMatchers: []spec.ExpFlagSpec{},
			ActionFlags: []spec.ExpFlagSpec{
				&spec.ExpFlag{
					Name: "rbyte",
					Desc: "bytes limit for disk read",
				},
				&spec.ExpFlag{
					Name: "wbyte",
					Desc: "bytes limit for disk write",
				},
			},
			ActionExecutor: &BlockIOExecutor{},
		},
	}
}

func (*BlockActionSpec) Name() string {
	return "block"
}

func (*BlockActionSpec) Aliases() []string {
	return []string{}
}
func (*BlockActionSpec) ShortDesc() string {
	return "Limit disk read or write bytes in a second"
}

func (*BlockActionSpec) LongDesc() string {
	return "Limit disk read or write bytes in a second"
}

type BlockIOExecutor struct {
	channel spec.Channel
}

func (*BlockIOExecutor) Name() string {
	return "Block"
}

var BlockIOBin = "chaos_blockio"

func (be *BlockIOExecutor) Exec(uid string, ctx context.Context, model *spec.ExpModel) *spec.Response {
	err := checkDiskExpEnv()
	if err != nil {
		return spec.ReturnFail(spec.Code[spec.CommandNotFound], err.Error())
	}
	if be.channel == nil {
		return spec.ReturnFail(spec.Code[spec.ServerError], "channel is nil")
	}

	if _, ok := spec.IsDestroy(ctx); ok {
		return be.stop(ctx)
	}

	rbyte := model.ActionFlags["rbyte"]
	wbyte := model.ActionFlags["wbyte"]
	if rbyte=="" && wbyte==""{
		return spec.ReturnFail(spec.Code[spec.IllegalParameters],
			"--rbyte and --wbyte : at least one of the two is not empty")
	}
	if rbyte!=""{
		rbyteI, err := strconv.Atoi(rbyte)
		if err != nil {
			return spec.ReturnFail(spec.Code[spec.IllegalParameters],
				"--rbyte value must be a positive integer")
		}
		if rbyteI <= 0 {
			return spec.ReturnFail(spec.Code[spec.IllegalParameters],
				"--rbyte value must be a prositive integer and bigger than 0")
		}
	}
	if wbyte!=""{
		wbyteI, err := strconv.Atoi(wbyte)
		if err != nil {
			return spec.ReturnFail(spec.Code[spec.IllegalParameters],
				"--wbyte value must be a positive integer")
		}
		if wbyteI <= 0 {
			return spec.ReturnFail(spec.Code[spec.IllegalParameters],
				"--wbyte value must be a prositive integer and bigger than 0")
		}
	}
	return be.start(ctx, rbyte, wbyte)
}

func (be *BlockIOExecutor) start(ctx context.Context, readBytes, writeBytes string) *spec.Response {
	return be.channel.Run(ctx, path.Join(be.channel.GetScriptPath(), BlockIOBin),
		fmt.Sprintf("--rbyte=%s --wbyte=%s --start --debug=%t", readBytes, writeBytes, util.Debug))
}

func (be *BlockIOExecutor) stop(ctx context.Context) *spec.Response {
	return be.channel.Run(ctx, path.Join(be.channel.GetScriptPath(), BlockIOBin),
		fmt.Sprintf("--stop --debug=%t", util.Debug))
}

func (be *BlockIOExecutor) SetChannel(channel spec.Channel) {
	be.channel = channel
}
