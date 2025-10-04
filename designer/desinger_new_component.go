package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 组件设计创建管理

// 设计组件
type DesigningComponent struct {
	owner        *FormTab                // 所属设计面板
	originObject any                     // 原始组件对象
	object       lcl.IWinControl         // 组件
	drag         *drag                   // 拖拽控制
	dx, dy       int32                   // 拖拽控制
	dcl, dct     int32                   // 拖拽控制
	isDown       bool                    // 拖拽控制
	propertyList []*vtedit.TEditNodeData // 组件属性
	eventList    []*vtedit.TEditNodeData // 组件事件
	isDesigner   bool                    // 组件是否正在设计
}

// 设计组件鼠标移动
func (m *DesigningComponent) OnMouseMove(sender lcl.IObject, shift types.TShiftState, X int32, Y int32) {
	br := m.object.BoundsRect()
	hint := fmt.Sprintf(`%v: %v
	Left: %v Top: %v
	Width: %v Height: %v`, m.object.Caption(), m.object.ToString(), br.Left, br.Top, br.Width(), br.Height())
	m.object.SetHint(hint)
	if m.isDown {
		m.drag.Hide()
		point := m.object.ClientToParent(types.TPoint{X: X, Y: Y}, m.owner.designerBox.object)
		x := point.X - m.dx
		y := point.Y - m.dy
		m.object.SetBounds(m.dcl+x, m.dct+y, br.Width(), br.Height())
	}
}

func (m *DesigningComponent) OnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	logs.Debug("OnMouseDown 设计组件", m.object.ToString())
	if !m.owner.placeComponent(m.object, X, Y) {
		m.isDown = true
		point := m.object.ClientToParent(types.TPoint{X: X, Y: Y}, m.owner.designerBox.object)
		m.dx, m.dy = point.X, point.Y
		m.dcl = m.object.Left()
		m.dct = m.object.Top()
		// 更新设计查看器的属性信息
		go lcl.RunOnMainThreadAsync(func(id uint32) {
			m.LoadPropertyToInspector()
			m.owner.hideAllDrag()
			m.drag.Show()
		})
	}
}

func (m *DesigningComponent) OnMouseUp(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if m.isDown {
		m.drag.Show()
	}
	m.isDown = false
}

func (m *DesigningComponent) SetObject(object any) {
	m.object = lcl.AsWinControl(object)
	m.originObject = object
}

// 加载组件属性到设计器
func (m *DesigningComponent) LoadPropertyToInspector() {
	// 显示设计组件拖拽
	m.drag.Show()
	// 加到到设计器
	inspector.LoadComponent(m)
}

func (m *DesigningComponent) SetParent(value lcl.IWinControl) {
	m.object.SetParent(value)
	m.drag.SetParent(value)
}

// 创建设计窗体-隐藏
func NewFormDesigner(designerForm *FormTab) *DesigningComponent {
	m := new(DesigningComponent)
	comp := lcl.NewForm(nil)
	comp.SetWidth(defaultWidth)
	comp.SetHeight(defaultHeight)
	comp.SetCaption(designerForm.name)
	comp.SetName(designerForm.name)
	comp.SetVisible(false) // 隐藏
	m.object = lcl.AsWinControl(comp)
	return m
}

// 创建设计按钮
func NewButtonDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := new(DesigningComponent)
	m.owner = designerForm
	designerForm.addDesignerComponent(m)
	comp := lcl.NewButton(designerForm.designerBox.object)
	//comp.SetParent(designerForm.designerBox.object)
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetCaption(designerForm.GetComponentCaptionName("Button"))
	comp.SetName(comp.Caption())
	comp.SetShowHint(true)
	comp.SetOnMouseMove(m.OnMouseMove)
	comp.SetOnMouseDown(m.OnMouseDown)
	comp.SetOnMouseUp(m.OnMouseUp)
	m.drag = newDrag(designerForm.designerBox.object, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 创建设计编辑框
func NewEditDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := new(DesigningComponent)
	m.owner = designerForm
	designerForm.addDesignerComponent(m)
	comp := lcl.NewEdit(designerForm.designerBox.object)
	//comp.SetParent(designerForm.designerBox.object)
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetCaption(designerForm.GetComponentCaptionName("Edit"))
	comp.SetName(comp.Caption())
	comp.SetText(comp.Caption())
	comp.SetShowHint(true)
	comp.SetOnMouseMove(m.OnMouseMove)
	comp.SetOnMouseDown(m.OnMouseDown)
	comp.SetOnMouseUp(m.OnMouseUp)
	m.drag = newDrag(designerForm.designerBox.object, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

func NewCheckBoxDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := new(DesigningComponent)
	m.owner = designerForm
	designerForm.addDesignerComponent(m)
	comp := lcl.NewCheckBox(designerForm.designerBox.object)
	//comp.SetParent(designerForm.designerBox.object)
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetCaption(designerForm.GetComponentCaptionName("CheckBox"))
	comp.SetName(comp.Caption())
	comp.SetChecked(false)
	comp.SetShowHint(true)
	comp.SetOnMouseMove(m.OnMouseMove)
	comp.SetOnMouseDown(m.OnMouseDown)
	comp.SetOnMouseUp(m.OnMouseUp)
	comp.SetOnChange(func(sender lcl.IObject) {
		comp.SetChecked(false)
	})
	m.drag = newDrag(designerForm.designerBox.object, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

func NewPanelDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := new(DesigningComponent)
	m.owner = designerForm
	designerForm.addDesignerComponent(m)
	comp := lcl.NewPanel(designerForm.designerBox.object)
	//comp.SetParent(designerForm.designerBox.object)
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetCaption(designerForm.GetComponentCaptionName("Panel"))
	comp.SetName(comp.Caption())
	comp.SetShowHint(true)
	comp.SetOnMouseMove(m.OnMouseMove)
	comp.SetOnMouseDown(m.OnMouseDown)
	comp.SetOnMouseUp(m.OnMouseUp)
	m.drag = newDrag(designerForm.designerBox.object, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}
