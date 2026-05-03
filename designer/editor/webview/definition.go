package webview

import (
	"encoding/json"
	"sync"

	"github.com/energye/designer/designer/editor"
	"github.com/energye/energy/v3/ipc"
	"github.com/energye/lcl/lcl"
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

		plcClient := editor.PLSClient()
		if plcClient == nil {
			context.Result("null")
			return
		}

		context.Result("")

		go func() {
			fileURI := editor.FilePathToURI(params.File)
			locations, err := plcClient.Definition(fileURI, params.Line, params.Column)
			if err != nil || len(locations) == 0 {
				lcl.RunOnMainThreadAsync(func(id uint32) {
					ipc.Emit("gopls-definition-response", params.RequestID, "null")
				})
				return
			}

			loc := locations[0]
			filePath := editor.URIToFilePath(loc.URI)
			if filePath == "" {
				lcl.RunOnMainThreadAsync(func(id uint32) {
					ipc.Emit("gopls-definition-response", params.RequestID, "null")
				})
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
			lcl.RunOnMainThreadAsync(func(id uint32) {
				ipc.Emit("gopls-definition-response", params.RequestID, string(respData))
			})
		}()
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
