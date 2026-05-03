package webview

import (
	"encoding/json"
	"github.com/energye/energy/v3/ipc"
)

// parseIPCParams extracts the first argument from IPC data and unmarshals it into target.
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
