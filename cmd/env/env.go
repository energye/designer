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

package env

import "fmt"

var env = make(envs)

type envs map[string]string

func Delete(name string) {
	delete(env, name)
}

func HasName(name string) bool {
	_, ok := env[name]
	return ok
}

func Put(name, value string) {
	env[name] = value
}

func Get(name string) string {
	return env[name]
}

func ToEnviron() []string {
	var s []string
	for k, v := range env {
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}
	return s
}

func Clear() {
	env = make(envs)
}
