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

package webview

import (
	"encoding/json"
	"github.com/energye/designer/designer/editor"
	"github.com/energye/designer/designer/editor/gopls"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/energy/v3/ipc"
	"github.com/energye/lcl/lcl"
	"os"
	"strings"
	"time"
)

// JSCompletionItem 发送给前端的补全项
type JSCompletionItem struct {
	Label               string       `json:"label"`
	Kind                int          `json:"kind"`
	Detail              string       `json:"detail,omitempty"`
	Documentation       string       `json:"documentation,omitempty"`
	SortText            string       `json:"sortText,omitempty"`
	FilterText          string       `json:"filterText,omitempty"`
	InsertText          string       `json:"insertText,omitempty"`
	InsertTextFormat    int          `json:"insertTextFormat,omitempty"`
	AdditionalTextEdits []JSTextEdit `json:"additionalTextEdits,omitempty"`
	Preselect           bool         `json:"preselect,omitempty"`
	Deprecated          bool         `json:"deprecated,omitempty"`
}

type JSTextEdit struct {
	Range   JSRange `json:"range"`
	NewText string  `json:"newText"`
}

type JSRange struct {
	Start JSPosition `json:"start"`
	End   JSPosition `json:"end"`
}

type JSPosition struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type JSParameter struct {
	Label         string `json:"label"`
	Documentation string `json:"documentation,omitempty"`
}

type JSSignature struct {
	Label         string        `json:"label"`
	Documentation string        `json:"documentation,omitempty"`
	Parameters    []JSParameter `json:"parameters,omitempty"`
}

type JSSignatureHelpResult struct {
	Signatures      []JSSignature `json:"signatures"`
	ActiveSignature int           `json:"activeSignature"`
	ActiveParameter int           `json:"activeParameter"`
}

type JSWorkspaceEdit struct {
	Changes map[string][]JSTextEdit `json:"changes"`
}

type JSCodeAction struct {
	Title       string           `json:"title"`
	Kind        string           `json:"kind,omitempty"`
	IsPreferred bool             `json:"isPreferred,omitempty"`
	Edit        *JSWorkspaceEdit `json:"edit,omitempty"`
}

// FileData 文件数据结构，用于IPC文件读写
type FileData struct {
	File     string `json:"file"`
	Content  string `json:"content"`
	Language string `json:"language"`
	ModTime  int64  `json:"modTime"`
	ReadOnly bool   `json:"readOnly"`
}

// CompletionParams gopls补全请求参数
type CompletionParams struct {
	RequestID   int    `json:"requestID"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	TriggerKind int    `json:"triggerKind"`
	TriggerChar string `json:"triggerChar"`
}

// SignatureHelpParams gopls签名帮助请求参数
type SignatureHelpParams struct {
	RequestID int    `json:"requestID"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
}

// CodeActionParams gopls代码操作请求参数
type CodeActionParams struct {
	RequestID   int                `json:"requestID"`
	File        string             `json:"file"`
	StartLine   int                `json:"startLine"`
	StartChar   int                `json:"startChar"`
	EndLine     int                `json:"endLine"`
	EndChar     int                `json:"endChar"`
	Kinds       string             `json:"kinds,omitempty"`
	Diagnostics []gopls.Diagnostic `json:"diagnostics,omitempty"`
}

// DidOpenParams gopls文件打开通知参数
type DidOpenParams struct {
	File       string `json:"file"`
	LanguageID string `json:"languageId"`
	Content    string `json:"content"`
	Version    int    `json:"version"`
}

// DidChangeParams gopls文件变更通知参数
type DidChangeParams struct {
	File    string `json:"file"`
	Content string `json:"content"`
	Version int    `json:"version"`
}

// FormattingParams gopls格式化请求参数
type FormattingParams struct {
	RequestID int    `json:"requestID"`
	File      string `json:"file"`
}

// RegData 文件注册数据
type RegData struct {
	File    string `json:"file"`
	ModTime int64  `json:"modTime"`
}

// DirtyData 文件脏状态数据
type DirtyData struct {
	File    string `json:"file"`
	IsDirty bool   `json:"isDirty"`
}

// readFileData 读取文件并返回序列化的FileData JSON，checkText为true时检查是否为文本文件
func readFileData(filePath string, checkText bool) (string, bool) {
	if checkText && !editor.IsTextFile(filePath) {
		return "", false
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", false
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", false
	}
	result := FileData{
		File:     filePath,
		Content:  string(content),
		Language: editor.DetectLanguage(filePath),
		ModTime:  fileInfo.ModTime().UnixMilli(),
		ReadOnly: editor.IsFileReadOnly(filePath),
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", false
	}
	return string(jsonData), true
}

func (m *TWebviewEditor) initIPCEvent() {
	ipc.On("monaco-inited", func(context ipc.IContext) {
		logs.Info("ipc monaco-inited BrowserId:", context.BrowserId(), context.Data())
		go func() {
			for i := 0; i < 30; i++ {
				if editor.IsPLSReady() || editor.IsPLSFailed() {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
			ready, failed := editor.PLSStatus()
			status := "ready"
			if failed {
				status = "unavailable"
			} else if !ready {
				status = "loading"
			}
			lcl.RunOnMainThreadAsync(func(id uint32) {
				ipc.Emit("gopls-status", status)
			})
		}()
		context.Result("")
	})

	// Completion
	ipc.On("gopls-completion", func(context ipc.IContext) {
		var params CompletionParams
		if !parseIPCParams(context, &params, "[]") {
			return
		}

		plsClient := editor.PLSClient()
		if plsClient == nil {
			context.Result("[]")
			return
		}

		go func() {
			fileURI := editor.FilePathToURI(params.File)
			result, err := plsClient.Completion(fileURI, params.Line, params.Column, params.TriggerKind, params.TriggerChar)
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
	ipc.On("gopls-signatureHelp", func(context ipc.IContext) {
		var params SignatureHelpParams
		if !parseIPCParams(context, &params, "") {
			return
		}

		plsClient := editor.PLSClient()
		if plsClient == nil {
			context.Result("")
			return
		}

		go func() {
			fileURI := editor.FilePathToURI(params.File)
			result, err := plsClient.SignatureHelp(fileURI, params.Line, params.Column)
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
	ipc.On("gopls-codeAction", func(context ipc.IContext) {
		var params CodeActionParams
		if !parseIPCParams(context, &params, "[]") {
			return
		}

		plsClient := editor.PLSClient()
		if plsClient == nil {
			context.Result("[]")
			return
		}

		var kinds []string
		if params.Kinds != "" {
			kinds = strings.Split(params.Kinds, ",")
		}

		go func() {
			fileURI := editor.FilePathToURI(params.File)
			result, err := plsClient.CodeAction(fileURI, params.StartLine, params.StartChar, params.EndLine, params.EndChar, kinds, params.Diagnostics)
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
						filePath := editor.URIToFilePath(uri)
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

	// DidOpen
	ipc.On("gopls-didOpen", func(context ipc.IContext) {
		var params DidOpenParams
		if !parseIPCParams(context, &params, "ok") {
			return
		}

		if plsClient := editor.PLSClient(); plsClient != nil {
			go func() {
				fileURI := editor.FilePathToURI(params.File)
				plsClient.DidOpen(fileURI, params.LanguageID, params.Content, params.Version)
			}()
		}
		context.Result("ok")
	})

	// DidChange - 异步：避免阻塞UI线程
	ipc.On("gopls-didChange", func(context ipc.IContext) {
		var params DidChangeParams
		if !parseIPCParams(context, &params, "ok") {
			return
		}

		go func() {
			if plsClient := editor.PLSClient(); plsClient != nil {
				fileURI := editor.FilePathToURI(params.File)
				if err := plsClient.DidChange(fileURI, params.Version, params.Content); err != nil {
					logs.Error("gopls-didChange 发送失败:", err)
				}
			}
		}()
		context.Result("ok")
	})

	// DidClose
	ipc.On("gopls-didClose", func(context ipc.IContext) {
		filePath, ok := parseIPCString(context, "ok")
		if !ok {
			return
		}

		if plsClient := editor.PLSClient(); plsClient != nil {
			go func() {
				fileURI := editor.FilePathToURI(filePath)
				plsClient.DidClose(fileURI)
			}()
		}
		context.Result("ok")
	})

	// Open file request
	ipc.On("open-file-request", func(context ipc.IContext) {
		filePath, ok := parseIPCString(context, "")
		if !ok {
			return
		}

		result, ok := readFileData(filePath, true)
		if !ok {
			context.Result("")
			return
		}
		context.Result(result)
	})

	// Save file
	ipc.On("save-file", func(context ipc.IContext) {
		var fileData FileData
		if !parseIPCParams(context, &fileData, "error: invalid data") {
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
			editor.UpdateSavedModTime(fileData.File)
		}

		context.Result("ok")

		if plsClient := editor.PLSClient(); plsClient != nil {
			go func() {
				fileURI := editor.FilePathToURI(fileData.File)
				plsClient.DidSave(fileURI, fileData.Content)
			}()
		}
	})

	// Reload file request
	ipc.On("reload-file-request", func(context ipc.IContext) {
		filePath, ok := parseIPCString(context, "")
		if !ok {
			return
		}

		result, ok := readFileData(filePath, false)
		if !ok {
			context.Result("")
			return
		}
		context.Result(result)
	})

	// Register opened file
	ipc.On("register-opened-file", func(context ipc.IContext) {
		var regData RegData
		if !parseIPCParams(context, &regData, "error") {
			return
		}

		m.fileManager.RegisterFile(regData.File)
		context.Result("ok")
	})

	// Unregister opened file
	ipc.On("unregister-opened-file", func(context ipc.IContext) {
		filePath, ok := parseIPCString(context, "error")
		if !ok {
			return
		}

		m.fileManager.UnregisterFile(filePath)
		context.Result("ok")
	})

	// Set file dirty
	ipc.On("set-file-dirty", func(context ipc.IContext) {
		var dirtyData DirtyData
		if !parseIPCParams(context, &dirtyData, "error") {
			return
		}

		m.fileManager.SetDirty(dirtyData.File, dirtyData.IsDirty)
		context.Result("ok")
	})

	// Formatting
	ipc.On("gopls-formatting", func(context ipc.IContext) {
		var params FormattingParams
		if !parseIPCParams(context, &params, "[]") {
			return
		}

		plcClient := editor.PLSClient()
		if plcClient == nil {
			context.Result("[]")
			return
		}

		go func() {
			fileURI := editor.FilePathToURI(params.File)
			result, err := plcClient.Formatting(fileURI)
			if err != nil || result == nil {
				lcl.RunOnMainThreadAsync(func(id uint32) {
					ipc.Emit("gopls-formatting-response", params.RequestID, "[]")
				})
				return
			}

			jsEdits := make([]JSTextEdit, len(result))
			for i, te := range result {
				jsEdits[i] = JSTextEdit{
					NewText: te.NewText,
					Range: JSRange{
						Start: JSPosition{Line: te.Range.Start.Line, Character: te.Range.Start.Character},
						End:   JSPosition{Line: te.Range.End.Line, Character: te.Range.End.Character},
					},
				}
			}

			respData, _ := json.Marshal(jsEdits)
			lcl.RunOnMainThreadAsync(func(id uint32) {
				ipc.Emit("gopls-formatting-response", params.RequestID, string(respData))
			})
		}()

		context.Result("")
	})
}
