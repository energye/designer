//go:build synedit

package editor

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// TSynEditEditor 原生 LCL SynEdit 编辑器实现
// 通过 build tag "synedit" 控制，默认不编译
type TSynEditEditor struct {
	synEdit     lcl.ISynEdit
	fileManager *FileManager
}

func NewSynEditEditor(owner lcl.IWinControl) IEditor {
	InitPLS()
	m := &TSynEditEditor{
		fileManager: NewFileManager(),
	}
	m.synEdit = lcl.NewSynEdit(owner)
	m.synEdit.SetAlign(types.AlClient)
	// TODO: 通过 LCL 事件直接调用 PLSClient() 实现代码补全、诊断等功能
	SetCurrentEditor(m)
	return m
}

func (m *TSynEditEditor) Type() EditType { return EtSynEdit }

func (m *TSynEditEditor) OpenFile(filePath string, readOnly ...bool) {
	// TODO: 读取文件内容设置到 SynEdit
}

func (m *TSynEditEditor) CloseFile(filePath string) {
	// TODO: 关闭文件，清理状态
}

func (m *TSynEditEditor) SaveCurrentFile() {
	// TODO: 保存当前文件
}

func (m *TSynEditEditor) FileManager() *FileManager {
	return m.fileManager
}

func (m *TSynEditEditor) Stop() {
	// TODO: 清理资源
}
