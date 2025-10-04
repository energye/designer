package designer

import "github.com/energye/designer/pkg/logs"

// 设计 - 组件设计树数据加载

// 添加设计组件到组件树
func (m *InspectorComponentTree) AddComponentToTree(component *DesigningComponent) {
	if component == nil {
		logs.Error("组件设计树加载节点失败, 设计组件为空")
		return
	}
	switch component.componentType {
	case CtForm:
		// 窗体表单根节点
		m.AddFormNode(nil, "Form1: TForm", -1)
	case CtOther:
		// 控件子节点

	default:
		logs.Error("组件设计树加载节点失败, 未知类型", component.componentType, "组件:", component.object.ToString())
	}
}
