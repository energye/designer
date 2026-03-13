// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package dflag

import (
	"github.com/energye/designer/pkg/tool"
	"os"
)

func New() *command {
	return &command{commands: tool.NewArrayMap[string, *Command]()}
}

type Command struct {
	Name string
	Long string
	Run  func(args []string)
}

type command struct {
	commands *tool.ArrayMap[string, *Command]
}

func (m *command) Add(cmd *Command) {
	if cmd == nil || cmd.Name == "" || cmd.Run == nil {
		return
	}
	m.commands.Add(cmd.Name, cmd)
}

func (m *command) Parse() {
	args := os.Args[1:]
	if len(args) == 0 {
		keys := m.commands.Keys()
		for _, name := range keys {
			println(name)
			cmd := m.commands.Get(name)
			println("  ", cmd.Long)
		}
		return
	}
	if cmd := m.commands.Get(args[0]); cmd != nil {
		cmd.Run(args[1:])
	}
}
