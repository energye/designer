package editor

import (
	"encoding/json"
	"github.com/energye/energy/v3/ipc"
)

type JSDefinition struct {
	File  string  `json:"file"`
	Range JSRange `json:"range"`
}

var onGoToDefinition func(filePath string, line, character int)

func SetOnGoToDefinition(handler func(filePath string, line, character int)) {
	onGoToDefinition = handler
}

func (m *TWebviewEditor) initDefinitionIPC() {
	type DefinitionParams struct {
		RequestID int    `json:"requestID"`
		File      string `json:"file"`
		Line      int    `json:"line"`
		Column    int    `json:"column"`
	}

	ipc.On("gopls-definition", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("null")
			return
		}

		var params DefinitionParams
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &params); err != nil {
			context.Result("null")
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
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("ok")
			return
		}

		var def JSDefinition
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &def); err != nil {
			context.Result("ok")
			return
		}

		if onGoToDefinition != nil {
			onGoToDefinition(def.File, def.Range.Start.Line, def.Range.Start.Character)
		}
		context.Result("ok")
	})
}
