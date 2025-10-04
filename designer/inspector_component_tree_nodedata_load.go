package designer

import "github.com/energye/designer/pkg/logs"

// 设计 - 组件设计树数据加载

func (m *InspectorComponentTree) Load(component *DesigningComponent) {
	if component == nil {
		logs.Error("组件设计树加载节点失败, 设计组件为空")
		return
	}
	switch component.componentType {
	case CtForm:
		m.AddTreeNodeItem(nil, "Form1: TForm", -1)
	case CtOther:
	default:
		logs.Error("组件设计树加载节点失败, 未知类型", component.componentType, "组件:", component.object.ToString())
	}
}
