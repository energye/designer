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
	"os"
	"strings"

	"github.com/energye/designer/designer/editor"
	"github.com/energye/designer/designer/editor/gopls"
	"github.com/energye/energy/v3/ipc"
	"github.com/energye/lcl/lcl"
)

// emitPLSResponse 统一把异步结果推回代码编辑页；抽成变量便于单元测试注入。
var emitPLSResponse = func(reqID int, event, payload string) {
	lcl.RunOnMainThreadAsync(func(id uint32) {
		ipc.Emit(event, reqID, payload)
	})
}

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

// JSHoverContent 发给前端的一段 hover 内容；Language 非空时前端包成代码块。
type JSHoverContent struct {
	Value    string `json:"value"`
	Language string `json:"language,omitempty"`
}

// JSHover 发给前端的悬停提示，range 为 0-based LSP 坐标，前端转 1-based。
type JSHover struct {
	Contents []JSHoverContent `json:"contents"`
	Range    *JSRange         `json:"range,omitempty"`
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

// HoverParams gopls悬停提示请求参数
type HoverParams struct {
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
func readFileData(filePath string, checkText bool) (*FileData, bool) {
	if checkText && !editor.IsTextFile(filePath) {
		return nil, false
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, false
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, false
	}
	result := &FileData{
		File:     filePath,
		Content:  string(content),
		Language: editor.DetectLanguage(filePath),
		ModTime:  fileInfo.ModTime().UnixMilli(),
		ReadOnly: editor.IsFileReadOnly(filePath),
	}
	return result, true
}

// 下面这些是纯转换函数：把 IPC 参数变成派发队列里的通知 / 发给前端的数据，
// 方便在不启动 LCL/WebView 的单元测试里覆盖路径转 URI、编辑项映射等逻辑。
func didOpenNotification(p DidOpenParams) editor.PLSNotification {
	return editor.PLSNotification{
		Kind:       editor.PLSNotifyOpen,
		URI:        editor.FilePathToURI(p.File),
		LanguageID: p.LanguageID,
		Content:    p.Content,
		Version:    p.Version,
	}
}

func didChangeNotification(p DidChangeParams) editor.PLSNotification {
	return editor.PLSNotification{
		Kind:    editor.PLSNotifyChange,
		URI:     editor.FilePathToURI(p.File),
		Content: p.Content,
		Version: p.Version,
	}
}

func didCloseNotification(filePath string) editor.PLSNotification {
	return editor.PLSNotification{
		Kind: editor.PLSNotifyClose,
		URI:  editor.FilePathToURI(filePath),
	}
}

func didSaveNotification(d FileData) editor.PLSNotification {
	return editor.PLSNotification{
		Kind:    editor.PLSNotifySave,
		URI:     editor.FilePathToURI(d.File),
		Content: d.Content,
	}
}

// hoverToJS 把 gopls 的 hover 结果转成前端结构；纯函数便于单测。
func hoverToJS(h *gopls.Hover) *JSHover {
	if h == nil || len(h.Contents) == 0 {
		return nil
	}
	result := &JSHover{Contents: make([]JSHoverContent, len(h.Contents))}
	for i, c := range h.Contents {
		result.Contents[i] = JSHoverContent{Value: c.Value, Language: c.Language}
	}
	if h.Range != nil {
		result.Range = &JSRange{
			Start: JSPosition{Line: h.Range.Start.Line, Character: h.Range.Start.Character},
			End:   JSPosition{Line: h.Range.End.Line, Character: h.Range.End.Character},
		}
	}
	return result
}

func jsTextEdits(edits []gopls.TextEdit) []JSTextEdit {
	result := make([]JSTextEdit, len(edits))
	for i, te := range edits {
		result[i] = JSTextEdit{
			NewText: te.NewText,
			Range: JSRange{
				Start: JSPosition{Line: te.Range.Start.Line, Character: te.Range.Start.Character},
				End:   JSPosition{Line: te.Range.End.Line, Character: te.Range.End.Character},
			},
		}
	}
	return result
}

// workspaceEditToJS 把 WorkspaceEdit 的 URI key 还原成本地路径，前端按
// 归一化后的 filePath 查 Monaco model；URI 转换失败时保留原始 key。
func workspaceEditToJS(edit *gopls.WorkspaceEdit) *JSWorkspaceEdit {
	if edit == nil || edit.Changes == nil {
		return nil
	}
	result := &JSWorkspaceEdit{Changes: make(map[string][]JSTextEdit, len(edit.Changes))}
	for uri, edits := range edit.Changes {
		key := editor.URIToFilePath(uri)
		if key == "" {
			key = uri
		}
		result.Changes[key] = jsTextEdits(edits)
	}
	return result
}

func (m *TWebviewEditor) initIPCEvent() {
	// Completion
	ipc.On("gopls-completion", func(context ipc.IContext) {
		var params CompletionParams
		if !parseIPCParams(context, &params, "[]") {
			return
		}

		plsClient := editor.PLSClient()
		if plsClient == nil {
			// 前端只等 *-response 事件，不用 emit 回调；必须显式回空结果，
			// 否则 pending 会一直挂到超时。
			emitPLSResponse(params.RequestID, "gopls-completion-response", "[]")
			context.Result("")
			return
		}

		go func() {
			fileURI := editor.FilePathToURI(params.File)
			result, err := plsClient.Completion(fileURI, params.Line, params.Column, params.TriggerKind, params.TriggerChar)
			if err != nil {
				emitPLSResponse(params.RequestID, "gopls-completion-response", "[]")
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
					jsItems[i].AdditionalTextEdits = jsTextEdits(item.AdditionalTextEdits)
				}
			}

			respData, err := json.Marshal(jsItems)
			if err != nil {
				emitPLSResponse(params.RequestID, "gopls-completion-response", "[]")
				return
			}
			emitPLSResponse(params.RequestID, "gopls-completion-response", string(respData))
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
			emitPLSResponse(params.RequestID, "gopls-signatureHelp-response", "")
			context.Result("")
			return
		}

		go func() {
			fileURI := editor.FilePathToURI(params.File)
			result, err := plsClient.SignatureHelp(fileURI, params.Line, params.Column)
			if err != nil || result == nil {
				emitPLSResponse(params.RequestID, "gopls-signatureHelp-response", "")
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
			emitPLSResponse(params.RequestID, "gopls-signatureHelp-response", string(respData))
		}()

		context.Result("")
	})

	// Hover（鼠标悬停提示）
	ipc.On("gopls-hover", func(context ipc.IContext) {
		var params HoverParams
		if !parseIPCParams(context, &params, "null") {
			return
		}

		plsClient := editor.PLSClient()
		if plsClient == nil {
			emitPLSResponse(params.RequestID, "gopls-hover-response", "null")
			context.Result("")
			return
		}

		go func() {
			fileURI := editor.FilePathToURI(params.File)
			result, err := plsClient.Hover(fileURI, params.Line, params.Column)
			if err != nil || result == nil {
				emitPLSResponse(params.RequestID, "gopls-hover-response", "null")
				return
			}
			respData, _ := json.Marshal(hoverToJS(result))
			emitPLSResponse(params.RequestID, "gopls-hover-response", string(respData))
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
			// 保存流程依赖这个响应继续往下走；gopls 不可用时也要回空数组，
			// 否则 manualSave/autoSave 会卡在 organize imports 阶段。
			emitPLSResponse(params.RequestID, "gopls-codeAction-response", "[]")
			context.Result("")
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
				emitPLSResponse(params.RequestID, "gopls-codeAction-response", "[]")
				return
			}

			jsActions := make([]JSCodeAction, len(result))
			for i, action := range result {
				jsActions[i] = JSCodeAction{
					Title:       action.Title,
					Kind:        action.Kind,
					IsPreferred: action.IsPreferred,
					Edit:        workspaceEditToJS(action.Edit),
				}
			}

			respData, _ := json.Marshal(jsActions)
			emitPLSResponse(params.RequestID, "gopls-codeAction-response", string(respData))
		}()

		context.Result("")
	})

	// DidOpen / DidChange / DidClose / DidSave 一律走同一个 FIFO 队列。
	// 旧实现每个通知各起 goroutine，重载时 didClose→didOpen 可能被在途的
	// didChange 插队，gopls 会拿到顺序颠倒且 version 回退的文档状态。
	// 队列还会在 gopls 完成 initialize 之前保留通知，不再静默丢弃。
	ipc.On("gopls-didOpen", func(context ipc.IContext) {
		var params DidOpenParams
		if !parseIPCParams(context, &params, "ok") {
			return
		}
		editor.EnqueuePLSNotification(didOpenNotification(params))
		context.Result("ok")
	})

	ipc.On("gopls-didChange", func(context ipc.IContext) {
		var params DidChangeParams
		if !parseIPCParams(context, &params, "ok") {
			return
		}
		editor.EnqueuePLSNotification(didChangeNotification(params))
		context.Result("ok")
	})

	ipc.On("gopls-didClose", func(context ipc.IContext) {
		filePath, ok := parseIPCString(context, "ok")
		if !ok {
			return
		}
		editor.EnqueuePLSNotification(didCloseNotification(filePath))
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

		// 只读保护：前端标记的只读（设计器生成代码）与磁盘层面的只读都拒绝写入。
		if fileData.ReadOnly || editor.IsFileReadOnly(fileData.File) {
			context.Result("error: file is read-only")
			return
		}

		// 覆盖已有文件时保留原权限，避免把用户设置的权限重置为 0644。
		perm := os.FileMode(0644)
		if fi, err := os.Stat(fileData.File); err == nil {
			perm = fi.Mode().Perm()
		}
		if err := os.WriteFile(fileData.File, []byte(fileData.Content), perm); err != nil {
			context.Result("error: " + err.Error())
			return
		}

		m.fileManager.SetDirty(fileData.File, false)
		fileInfo, _ := os.Stat(fileData.File)
		if fileInfo != nil {
			m.fileManager.UpdateModTime(fileData.File, fileInfo.ModTime())
			editor.UpdateSavedModTime(fileData.File)
		}
		// 若此前因未保存内容暂缓了代码生成，脏状态清空后重放。
		editor.NotifyFileSaved()

		context.Result("ok")

		// DidSave 同样入队，保证排在之前的 didChange 之后到达 gopls。
		editor.EnqueuePLSNotification(didSaveNotification(fileData))
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
			emitPLSResponse(params.RequestID, "gopls-formatting-response", "[]")
			context.Result("")
			return
		}

		go func() {
			fileURI := editor.FilePathToURI(params.File)
			result, err := plcClient.Formatting(fileURI)
			if err != nil || result == nil {
				emitPLSResponse(params.RequestID, "gopls-formatting-response", "[]")
				return
			}

			respData, _ := json.Marshal(jsTextEdits(result))
			emitPLSResponse(params.RequestID, "gopls-formatting-response", string(respData))
		}()

		context.Result("")
	})
}
