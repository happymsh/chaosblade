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

	"github.com/chaosblade-io/chaosblade-spec-go/channel"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
)

type Executor struct {
	executors map[string]spec.Executor
}

func NewExecutor() spec.Executor {
	return &Executor{
		executors: getAllOsExecutors(),
	}
}

func getAllOsExecutors() (executors map[string]spec.Executor) {
	executors = make(map[string]spec.Executor, 0)
	expModels := GetAllExpModels()
	for _, expModel := range expModels {
		executorMap := ExtractExecutorFromExpModel(expModel)
		for key, value := range executorMap {
			executors[key] = value
		}
	}
	return executors
}

func GetAllExpModels() []spec.ExpModelCommandSpec {
	return []spec.ExpModelCommandSpec{
		NewCpuCommandModelSpec(),
		NewJvmCommandModelSpec(),
	}
}

func ExtractExecutorFromExpModel(expModel spec.ExpModelCommandSpec) map[string]spec.Executor {
	executors := make(map[string]spec.Executor)
	for _, actionModel := range expModel.Actions() {
		executors[expModel.Name()+actionModel.Name()] = actionModel.Executor()
	}
	return executors
}

func (*Executor) Name() string {
	return "hzcp"
}

func (e *Executor) Exec(uid string, ctx context.Context, model *spec.ExpModel) *spec.Response {
	key := model.Target + model.ActionName
	executor := e.executors[key]
	if executor == nil {
		return spec.ReturnFail(spec.Code[spec.HandlerNotFound], fmt.Sprintf("the hzcp executor not found, %s", key))
	}
	executor.SetChannel(channel.NewLocalChannel())
	return executor.Exec(uid, ctx, model)
}

func (*Executor) SetChannel(channel spec.Channel) {
}
