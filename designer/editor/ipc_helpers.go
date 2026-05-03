// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package editor

import (
	"encoding/json"
	"github.com/energye/energy/v3/ipc"
)

// parseIPCParams extracts the first argument from IPC data and unmarshals it into target.
// Returns false if data could not be parsed, in which case context.Result(defaultResult) is called.
func parseIPCParams(context ipc.IContext, target any, defaultResult string) bool {
	data := context.Data()
	arr, ok := data.([]any)
	if !ok || len(arr) == 0 {
		context.Result(defaultResult)
		return false
	}
	jsonData, _ := json.Marshal(arr[0])
	if err := json.Unmarshal(jsonData, target); err != nil {
		context.Result(defaultResult)
		return false
	}
	return true
}

// parseIPCString extracts the first argument as a plain string from IPC data.
// Returns ("", false) if data could not be parsed, in which case context.Result(defaultResult) is called.
func parseIPCString(context ipc.IContext, defaultResult string) (string, bool) {
	data := context.Data()
	arr, ok := data.([]any)
	if !ok || len(arr) == 0 {
		context.Result(defaultResult)
		return "", false
	}
	s, ok := arr[0].(string)
	if !ok {
		context.Result(defaultResult)
		return "", false
	}
	return s, true
}
