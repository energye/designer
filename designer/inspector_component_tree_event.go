package designer

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"unsafe"
)

// 设计 - 组件树 - 事件

func (m *FormTab) TreeOnContextPopup(sender lcl.IObject, mousePos types.TPoint, handled *bool) {
	logs.Debug("TreeOnContextPopup pos:", mousePos)
}

func (m *FormTab) TreeOnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	logs.Debug("TreeOnMouseDown x,y:", X, Y)
	if button == types.MbRight {
		selectNode := m.tree.GetNodeAt(X, Y)
		if selectNode.IsValid() {
			m.tree.SetSelected(selectNode)
		}
	}
}

// 数据指针转设计组件
func (m *FormTab) DataToDesigningComponent(data uintptr) *TDesigningComponent {
	dc := (*TDesigningComponent)(unsafe.Pointer(data))
	return dc
}

// 组件树选择事件
func (m *FormTab) TreeOnGetSelectedIndex(sender lcl.IObject, node lcl.ITreeNode) {
	data := node.Data()
	component := m.DataToDesigningComponent(data)
	if component != nil {
		m.switchComponentEditing(component)
	}
	logs.Info("Inspector-component-tree OnGetSelectedIndex name:", node.Text(), "id:", component.id)
}

// 取消选中所有节点
//func (m *InspectorComponentTree) UnSelectedAll() {
//	for _, node := range m.nodeData {
//		node.node.SetSelected(false)
//	}
//}
