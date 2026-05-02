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

package editor

import (
	"encoding/json"
	"github.com/energye/designer/designer/editor/gopls"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/energy/v3/ipc"
	"github.com/energye/lcl/lcl"
	"os"
	"strings"
)

func (m *TWebviewEditor) initIPCEvent() {
	ipc.On("monaco-inited", func(context ipc.IContext) {
		logs.Info("ipc monaco-inited BrowserId:", context.BrowserId(), context.Data())
	})

	type CompletionParams struct {
		RequestID   int    `json:"requestID"`
		File        string `json:"file"`
		Line        int    `json:"line"`
		Column      int    `json:"column"`
		TriggerKind int    `json:"triggerKind"`
		TriggerChar string `json:"triggerChar"`
	}

	ipc.On("gopls-completion", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("[]")
			return
		}

		var params CompletionParams
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &params); err != nil {
			logs.Error("gopls-completion: 解析参数失败:", err)
			context.Result("[]")
			return
		}

		if gLSPClient == nil {
			context.Result("[]")
			return
		}

		go func() {
			fileURI := filePathToURI(params.File)
			result, err := gLSPClient.Completion(fileURI, params.Line, params.Column, params.TriggerKind, params.TriggerChar)
			if err != nil {
				lcl.RunOnMainThreadAsync(func(id uint32) {
					ipc.Emit("gopls-completion-response", params.RequestID, "[]")
				})
				return
			}

			jsItems := make([]JSCompletionItem, len(result))
			for i, item := range result {
				jsItems[i] = JSCompletionItem{
					Label:            item.Label,
					Kind:             item.Kind,
					Detail:           item.GetDetail(),
					Documentation:    item.GetDocumentation(),
					SortText:         item.SortText,
					FilterText:       item.FilterText,
					InsertText:       item.InsertText,
					InsertTextFormat: item.InsertTextFormat,
					Preselect:        item.Preselect,
					Deprecated:       item.Deprecated,
				}
				if jsItems[i].InsertText == "" {
					jsItems[i].InsertText = item.Label
				}
				if len(item.AdditionalTextEdits) > 0 {
					jsItems[i].AdditionalTextEdits = make([]JSTextEdit, len(item.AdditionalTextEdits))
					for j, te := range item.AdditionalTextEdits {
						jsItems[i].AdditionalTextEdits[j] = JSTextEdit{
							NewText: te.NewText,
							Range: JSRange{
								Start: JSPosition{Line: te.Range.Start.Line, Character: te.Range.Start.Character},
								End:   JSPosition{Line: te.Range.End.Line, Character: te.Range.End.Character},
							},
						}
					}
				}
			}

			respData, err := json.Marshal(jsItems)
			if err != nil {
				lcl.RunOnMainThreadAsync(func(id uint32) {
					ipc.Emit("gopls-completion-response", params.RequestID, "[]")
				})
				return
			}
			lcl.RunOnMainThreadAsync(func(id uint32) {
				ipc.Emit("gopls-completion-response", params.RequestID, string(respData))
			})
		}()

		context.Result("")
	})

	// Signature Help
	type SignatureHelpParams struct {
		RequestID int    `json:"requestID"`
		File      string `json:"file"`
		Line      int    `json:"line"`
		Column    int    `json:"column"`
	}

	ipc.On("gopls-signatureHelp", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("")
			return
		}

		var params SignatureHelpParams
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &params); err != nil {
			context.Result("")
			return
		}

		if gLSPClient == nil {
			context.Result("")
			return
		}

		go func() {
			fileURI := filePathToURI(params.File)
			result, err := gLSPClient.SignatureHelp(fileURI, params.Line, params.Column)
			if err != nil || result == nil {
				lcl.RunOnMainThreadAsync(func(id uint32) {
					ipc.Emit("gopls-signatureHelp-response", params.RequestID, "")
				})
				return
			}

			jsResult := JSSignatureHelpResult{
				ActiveSignature: result.ActiveSignature,
				ActiveParameter: result.ActiveParameter,
			}
			for _, sig := range result.Signatures {
				jsSig := JSSignature{
					Label:         sig.Label,
					Documentation: sig.GetDocumentation(),
				}
				for _, p := range sig.Parameters {
					jsSig.Parameters = append(jsSig.Parameters, JSParameter{
						Label:         p.GetLabel(),
						Documentation: gopls.ExtractString(p.Documentation),
					})
				}
				jsResult.Signatures = append(jsResult.Signatures, jsSig)
			}

			respData, _ := json.Marshal(jsResult)
			lcl.RunOnMainThreadAsync(func(id uint32) {
				ipc.Emit("gopls-signatureHelp-response", params.RequestID, string(respData))
			})
		}()

		context.Result("")
	})

	// Code Action
	type CodeActionParams struct {
		RequestID   int                     `json:"requestID"`
		File        string                  `json:"file"`
		StartLine   int                     `json:"startLine"`
		StartChar   int                     `json:"startChar"`
		EndLine     int                     `json:"endLine"`
		EndChar     int                     `json:"endChar"`
		Kinds       string                  `json:"kinds,omitempty"`
		Diagnostics []gopls.DiagnosticInput `json:"diagnostics,omitempty"`
	}

	ipc.On("gopls-codeAction", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("[]")
			return
		}

		var params CodeActionParams
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &params); err != nil {
			context.Result("[]")
			return
		}

		if gLSPClient == nil {
			context.Result("[]")
			return
		}

		var kinds []string
		if params.Kinds != "" {
			kinds = strings.Split(params.Kinds, ",")
		}

		go func() {
			fileURI := filePathToURI(params.File)
			result, err := gLSPClient.CodeAction(fileURI, params.StartLine, params.StartChar, params.EndLine, params.EndChar, kinds, params.Diagnostics)
			if err != nil || result == nil {
				lcl.RunOnMainThreadAsync(func(id uint32) {
					ipc.Emit("gopls-codeAction-response", params.RequestID, "[]")
				})
				return
			}

			jsActions := make([]JSCodeAction, len(result))
			for i, action := range result {
				jsActions[i] = JSCodeAction{
					Title:       action.Title,
					Kind:        action.Kind,
					IsPreferred: action.IsPreferred,
				}
				if action.Edit != nil && action.Edit.Changes != nil {
					jsActions[i].Edit = &JSWorkspaceEdit{
						Changes: make(map[string][]JSTextEdit),
					}
					for uri, edits := range action.Edit.Changes {
						filePath := uriToFilePath(uri)
						key := filePath
						if key == "" {
							key = uri
						}
						jsEdits := make([]JSTextEdit, len(edits))
						for j, te := range edits {
							jsEdits[j] = JSTextEdit{
								NewText: te.NewText,
								Range: JSRange{
									Start: JSPosition{Line: te.Range.Start.Line, Character: te.Range.Start.Character},
									End:   JSPosition{Line: te.Range.End.Line, Character: te.Range.End.Character},
								},
							}
						}
						jsActions[i].Edit.Changes[key] = jsEdits
					}
				}
			}

			respData, _ := json.Marshal(jsActions)
			lcl.RunOnMainThreadAsync(func(id uint32) {
				ipc.Emit("gopls-codeAction-response", params.RequestID, string(respData))
			})
		}()

		context.Result("")
	})

	type DidOpenParams struct {
		File       string `json:"file"`
		LanguageID string `json:"languageId"`
		Content    string `json:"content"`
		Version    int    `json:"version"`
	}

	ipc.On("gopls-didOpen", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("ok")
			return
		}

		var params DidOpenParams
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &params); err != nil {
			context.Result("ok")
			return
		}

		if gLSPClient != nil {
			go func() {
				fileURI := filePathToURI(params.File)
				gLSPClient.DidOpen(fileURI, params.LanguageID, params.Content, params.Version)
			}()
		}
		context.Result("ok")
	})

	type DidChangeParams struct {
		File    string `json:"file"`
		Content string `json:"content"`
		Version int    `json:"version"`
	}

	ipc.On("gopls-didChange", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("ok")
			return
		}

		var params DidChangeParams
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &params); err != nil {
			context.Result("ok")
			return
		}

		if gLSPClient != nil {
			// Synchronous: must complete before gopls processes subsequent
			// completion/codeAction requests with the updated file content
			fileURI := filePathToURI(params.File)
			if err := gLSPClient.DidChange(fileURI, params.Version, params.Content); err != nil {
				logs.Error("gopls-didChange 同步发送失败:", err)
			}
		}
		context.Result("ok")
	})

	ipc.On("gopls-didClose", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("ok")
			return
		}

		filePath, ok := arr[0].(string)
		if !ok {
			context.Result("ok")
			return
		}

		if gLSPClient != nil {
			go func() {
				fileURI := filePathToURI(filePath)
				gLSPClient.DidClose(fileURI)
			}()
		}
		context.Result("ok")
	})

	ipc.On("open-file-request", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("")
			return
		}

		filePath, ok := arr[0].(string)
		if !ok {
			context.Result("")
			return
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			context.Result("")
			return
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			context.Result("")
			return
		}

		result := FileData{
			File:     filePath,
			Content:  string(content),
			Language: detectLanguage(filePath),
			ModTime:  fileInfo.ModTime().UnixMilli(),
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			context.Result("")
			return
		}

		context.Result(string(jsonData))
	})

	ipc.On("save-file", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("error: invalid data")
			return
		}

		var fileData FileData
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &fileData); err != nil {
			context.Result("error: " + err.Error())
			return
		}

		if err := os.WriteFile(fileData.File, []byte(fileData.Content), 0644); err != nil {
			context.Result("error: " + err.Error())
			return
		}

		fileInfo, _ := os.Stat(fileData.File)
		if fileInfo != nil {
			m.fileManager.UpdateModTime(fileData.File, fileInfo.ModTime())
			m.fileManager.SetDirty(fileData.File, false)
		}

		context.Result("ok")

		// Notify gopls that file was saved so it can do full analysis
		if gLSPClient != nil {
			go func() {
				fileURI := filePathToURI(fileData.File)
				gLSPClient.DidSave(fileURI, fileData.Content)
			}()
		}
	})

	ipc.On("reload-file-request", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("")
			return
		}

		filePath, ok := arr[0].(string)
		if !ok {
			context.Result("")
			return
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			context.Result("")
			return
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			context.Result("")
			return
		}

		result := FileData{
			File:     filePath,
			Content:  string(content),
			Language: detectLanguage(filePath),
			ModTime:  fileInfo.ModTime().UnixMilli(),
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			context.Result("")
			return
		}

		context.Result(string(jsonData))
	})

	ipc.On("register-opened-file", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("error")
			return
		}

		type RegData struct {
			File    string `json:"file"`
			ModTime int64  `json:"modTime"`
		}
		var regData RegData
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &regData); err != nil {
			context.Result("error")
			return
		}

		m.fileManager.RegisterFile(regData.File)
		context.Result("ok")
	})

	ipc.On("unregister-opened-file", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("error")
			return
		}

		filePath, ok := arr[0].(string)
		if !ok {
			context.Result("error")
			return
		}

		m.fileManager.UnregisterFile(filePath)
		context.Result("ok")
	})

	ipc.On("set-file-dirty", func(context ipc.IContext) {
		data := context.Data()
		arr, ok := data.([]any)
		if !ok || len(arr) == 0 {
			context.Result("error")
			return
		}

		type DirtyData struct {
			File    string `json:"file"`
			IsDirty bool   `json:"isDirty"`
		}
		var dirtyData DirtyData
		jsonData, _ := json.Marshal(arr[0])
		if err := json.Unmarshal(jsonData, &dirtyData); err != nil {
			context.Result("error")
			return
		}

		m.fileManager.SetDirty(dirtyData.File, dirtyData.IsDirty)
		context.Result("ok")
	})
}
