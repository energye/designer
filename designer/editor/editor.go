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
	"github.com/energye/designer/designer/editor/gopls"
	"github.com/energye/lcl/lcl"
)

type EditType int

const (
	EtWebview = iota
	EtSynEdit
)

// LSPService defines the interface for Language Server Protocol operations.
// Implementations communicate with an LSP server (e.g., gopls) to provide
// code intelligence features like completion, diagnostics, and navigation.
type LSPService interface {
	Completion(fileURI string, line, column int, triggerKind int, triggerChar string) ([]gopls.CompletionItem, error)
	SignatureHelp(fileURI string, line, column int) (*gopls.SignatureHelpResult, error)
	CodeAction(fileURI string, startLine, startChar, endLine, endChar int, kinds []string, diagnostics []gopls.Diagnostic) ([]gopls.CodeAction, error)
	Definition(fileURI string, line, column int) ([]gopls.Location, error)
	DidOpen(fileURI, languageID, content string, version int) error
	DidChange(fileURI string, version int, content string) error
	DidSave(fileURI string, text string) error
	DidClose(fileURI string) error
	SetDiagnosticsHandler(handler func(uri string, diagnostics []gopls.Diagnostic))
}

type IEditor interface {
	Type() EditType
	OpenFile(filePath string, readOnly ...bool)
	CloseFile(filePath string)
	SaveCurrentFile()
	FileManager() *FileManager
	Stop()
}

func NewEditor(owner lcl.IWinControl, editType ...EditType) IEditor {
	var et EditType = EtWebview
	if len(editType) > 0 {
		et = editType[0]
	}
	switch et {
	case EtSynEdit:
		// SynEdit editor requires build tag "synedit"
		return nil
	default:
		return NewWebviewEditor(owner)
	}
}
