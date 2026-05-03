//go:build synedit

package synedit

import (
	"github.com/energye/designer/designer/editor"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// TSynEditEditor 原生 LCL SynEdit 编辑器实现
// 通过 build tag "synedit" 控制，默认不编译
type TSynEditEditor struct {
	synEdit     lcl.ISynEdit
	fileManager *editor.FileManager
}

func NewSynEditEditor(owner lcl.IWinControl) editor.IEditor {
	editor.InitPLS()
	m := &TSynEditEditor{
		fileManager: editor.NewFileManager(),
	}
	m.synEdit = lcl.NewSynEdit(owner)
	m.synEdit.SetAlign(types.AlClient)
	// TODO: 通过 LCL 事件直接调用 editor.PLSClient() 实现代码补全、诊断等功能
	editor.SetCurrentEditor(m)
	return m
}

func init() {
	editor.RegisterEditorFactory(editor.EtSynEdit, NewSynEditEditor)
}

func (m *TSynEditEditor) Type() editor.EditType { return editor.EtSynEdit }

func (m *TSynEditEditor) OpenFile(filePath string, readOnly ...bool) {
	// TODO: 读取文件内容设置到 SynEdit
}

func (m *TSynEditEditor) CloseFile(filePath string) {
	// TODO: 关闭文件，清理状态
}

func (m *TSynEditEditor) SaveCurrentFile() {
	// TODO: 保存当前文件
}

func (m *TSynEditEditor) FileManager() *editor.FileManager {
	return m.fileManager
}

func (m *TSynEditEditor) Stop() {
	// TODO: 清理资源
}
