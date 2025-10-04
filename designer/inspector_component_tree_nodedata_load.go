package designer

import "github.com/energye/designer/pkg/logs"

// 设计 - 组件的设计树数据加载

func (m *InspectorComponentTree) Load(component *DesigningComponent) {
	logs.Debug("")
	m.AddTreeNodeItem(nil, "Form1: TForm", -1)
}
