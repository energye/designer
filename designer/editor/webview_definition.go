package editor

import (
	"encoding/json"
	"sync"

	"github.com/energye/energy/v3/ipc"
)

type JSDefinition struct {
	File  string  `json:"file"`
	Range JSRange `json:"range"`
}

var (
	onGoToDefinition func(filePath string, line, character int)
	defMu            sync.RWMutex
)

func SetOnGoToDefinition(handler func(filePath string, line, character int)) {
	defMu.Lock()
	onGoToDefinition = handler
	defMu.Unlock()
}

// DefinitionParams gopls定义跳转请求参数
type DefinitionParams struct {
	RequestID int    `json:"requestID"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
}

func (m *TWebviewEditor) initDefinitionIPC() {
	ipc.On("gopls-definition", func(context ipc.IContext) {
		var params DefinitionParams
		if !parseIPCParams(context, &params, "null") {
			return
		}

		if gPLSClient == nil {
			context.Result("null")
			return
		}

		fileURI := filePathToURI(params.File)
		locations, err := gPLSClient.Definition(fileURI, params.Line, params.Column)
		if err != nil || len(locations) == 0 {
			context.Result("null")
			return
		}

		loc := locations[0]
		filePath := uriToFilePath(loc.URI)
		if filePath == "" {
			context.Result("null")
			return
		}

		result := JSDefinition{
			File: filePath,
			Range: JSRange{
				Start: JSPosition{Line: loc.Range.Start.Line, Character: loc.Range.Start.Character},
				End:   JSPosition{Line: loc.Range.End.Line, Character: loc.Range.End.Character},
			},
		}

		respData, _ := json.Marshal(result)
		context.Result(string(respData))
	})

	ipc.On("go-to-definition", func(context ipc.IContext) {
		var def JSDefinition
		if !parseIPCParams(context, &def, "ok") {
			return
		}

		defMu.RLock()
		handler := onGoToDefinition
		defMu.RUnlock()
		if handler != nil {
			handler(def.File, def.Range.Start.Line, def.Range.Start.Character)
		}
		context.Result("ok")
	})
}
