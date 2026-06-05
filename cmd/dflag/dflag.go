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
	"strings"
	"unicode"
)

// New 创建并返回一个新的命令实例
func New() *dFlag {
	return &dFlag{commands: tool.NewArrayMap[string, *Command]()}
}

type Args struct {
	list        *tool.HashMap[string, string]
	positionals []string
}

// Get 获取指定名称的参数值
//
//	name: 参数名称, "-path --path" > path
func (m *Args) Get(name string) string {
	return m.list.Get(name)
}

// Contains 检查指定名称的参数是否存在
//
//	name: 参数名称, "-path --path" > path
func (m *Args) Contains(name string) bool {
	return m.list.ContainsKey(name)
}

func (m *Args) Positionals() []string {
	result := make([]string, len(m.positionals))
	copy(result, m.positionals)
	return result
}

type Command struct {
	Name string
	Long string
	Run  func(args Args)
}

type dFlag struct {
	commands *tool.ArrayMap[string, *Command]
}

func (m *dFlag) Add(cmd *Command) {
	if cmd == nil || cmd.Name == "" || cmd.Run == nil {
		return
	}
	m.commands.Add(cmd.Name, cmd)
}

func (m *dFlag) Help() {
	keys := m.commands.Keys()
	println("energy command")
	for _, name := range keys {
		println("  ", name)
		cmd := m.commands.Get(name)
		println("    ", cmd.Long)
	}
}

func (m *dFlag) Parse() {
	newArgs := os.Args[1:]
	if len(newArgs) == 0 {
		m.Help()
		return
	}
	if cmd := m.commands.Get(newArgs[0]); cmd != nil {
		inArgs := &Args{list: tool.NewHashMap[string, string]()}
		cmdArgs := newArgs[1:]
		for i := 0; i < len(cmdArgs); i++ {
			el := strings.TrimSpace(cmdArgs[i])
			if el == "" {
				continue
			}
			if el[0] == '-' {
				v := ""
				if eqEl := strings.SplitN(el, "=", 2); len(eqEl) > 1 {
					el = eqEl[0]
					v = eqEl[1]
				} else if i+1 <= len(cmdArgs)-1 {
					tmpV := strings.TrimSpace(cmdArgs[i+1])
					if tmpV != "" && tmpV[0] != '-' {
						v = tmpV
						i++
					}
				}
				name := strings.TrimFunc(el, func(r rune) bool {
					if r == '-' {
						return true
					}
					return false
				})
				inArgs.list.Add(name, v)
			} else {
				inArgs.positionals = append(inArgs.positionals, el)
				inArgs.list.Add(el, "")
			}
		}
		if cmd.Run != nil {
			cmd.Run(*inArgs)
		}
	}
}

// ParseCommandLine 解析命令行字符串，正确处理带引号的参数
func ParseCommandLine(line string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	inSingle := false
	escape := false
	for _, r := range line {
		if escape {
			current.WriteRune(r)
			escape = false
			continue
		}
		switch {
		case r == '\\':
			escape = true
		case r == '"' && !inSingle:
			inQuote = !inQuote
		case r == '\'' && !inQuote:
			inSingle = !inSingle
		case unicode.IsSpace(r) && !inQuote && !inSingle:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
