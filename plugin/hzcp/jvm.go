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
	"fmt"

	"github.com/chaosblade-io/chaosblade-spec-go/channel"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
)

type JvmCommandModelSpec struct {
	spec.BaseExpModelCommandSpec
}

func NewJvmCommandModelSpec() spec.ExpModelCommandSpec {
	return &JvmCommandModelSpec{
		spec.BaseExpModelCommandSpec{
			ExpActions: []spec.ExpActionCommandSpec{
				&threadfullActionCommand{
					spec.BaseExpActionCommandSpec{
						ActionMatchers: []spec.ExpFlagSpec{},
						ActionFlags:    []spec.ExpFlagSpec{},
						ActionExecutor: &jvmThreadfullExecutor{},
					},
				},
				&oomActionCommand{
					spec.BaseExpActionCommandSpec{
						ActionMatchers: []spec.ExpFlagSpec{},
						ActionFlags:    []spec.ExpFlagSpec{},
						ActionExecutor: &jvmOomExecutor{},
					},
				},
			},
		},
	}
}

func (*JvmCommandModelSpec) Name() string {
	return "hzcp-jvm"
}

func (*JvmCommandModelSpec) ShortDesc() string {
	return "hzcp-jvm experiment"
}

func (*JvmCommandModelSpec) LongDesc() string {
	return "hzcp-jvm experiment, for example oom"
}

func (*JvmCommandModelSpec) Example() string {
	return "blade create hzcp-jvm oom "
}

func checkJvmExpEnv() error {
	commands := []string{"ps", "awk", "grep", "kill", "nohup", "tr"}
	for _, command := range commands {
		if !channel.IsCommandAvailable(command) {
			return fmt.Errorf("%s command not found", command)
		}
	}
	return nil
}
