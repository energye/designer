package designer

import (
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 组件设计创建管理

type DesigningComponent struct {
	owner    *FormTab
	object   lcl.IWinControl
	drag     *drag
	dx, dy   int32
	dcl, dct int32
	isDown   bool
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
		point := m.object.ClientToParent(types.TPoint{X: X, Y: Y}, m.owner.designerBox)
		x := point.X - m.dx
		y := point.Y - m.dy
		m.object.SetBounds(m.dcl+x, m.dct+y, br.Width(), br.Height())
	}
}

func (m *DesigningComponent) OnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	m.owner.hideAllDrag()
	m.drag.Show()
	m.isDown = true
	point := m.object.ClientToParent(types.TPoint{X: X, Y: Y}, m.owner.designerBox)
	m.dx, m.dy = point.X, point.Y
	m.dcl = m.object.Left()
	m.dct = m.object.Top()
	// 更新设计查看器的属性信息
	m.UpdateLoadPropertyInfo()
}

func (m *DesigningComponent) OnMouseUp(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if m.isDown {
		m.drag.Show()
	}
	m.isDown = false
}

// 更新设计查看器加载当前属性信息
func (m *DesigningComponent) UpdateLoadPropertyInfo() {
	inspector.LoadComponent(m)
}

// 创建设计按钮
func NewButtonDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := new(DesigningComponent)
	m.owner = designerForm
	designerForm.addDesignerComponent(m)
	comp := lcl.NewButton(designerForm.designerBox)
	comp.SetParent(designerForm.designerBox)
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetCaption(designerForm.GetComponentCaptionName("Button"))
	comp.SetShowHint(true)
	comp.SetOnMouseMove(m.OnMouseMove)
	comp.SetOnMouseDown(m.OnMouseDown)
	comp.SetOnMouseUp(m.OnMouseUp)
	comp.SetOnShowHint(func(sender lcl.IObject, hintInfo lcl.THintInfo) {
		//fmt.Printf("SetOnShowHint: %+v\n", hintInfo)
		fmt.Println("SetOnShowHint:", hintInfo.HintStr, hintInfo.HintPos)
	})
	m.drag = newDrag(designerForm.designerBox, DsAll)
	m.drag.SetRelation(comp)
	m.object = lcl.AsWinControl(comp)
	return m
}

// 创建设计编辑框
func NewEditDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := new(DesigningComponent)
	m.owner = designerForm
	designerForm.addDesignerComponent(m)
	comp := lcl.NewEdit(designerForm.designerBox)
	comp.SetParent(designerForm.designerBox)
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetCaption(designerForm.GetComponentCaptionName("Edit"))
	comp.SetText(comp.Caption())
	comp.SetShowHint(true)
	comp.SetOnMouseMove(m.OnMouseMove)
	comp.SetOnMouseDown(m.OnMouseDown)
	comp.SetOnMouseUp(m.OnMouseUp)
	m.drag = newDrag(designerForm.designerBox, DsAll)
	m.drag.SetRelation(comp)
	m.object = lcl.AsWinControl(comp)
	return m
}

func NewCheckBoxDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := new(DesigningComponent)
	m.owner = designerForm
	designerForm.addDesignerComponent(m)
	comp := lcl.NewCheckBox(designerForm.designerBox)
	comp.SetParent(designerForm.designerBox)
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetCaption(designerForm.GetComponentCaptionName("CheckBox"))
	comp.SetChecked(false)
	comp.SetShowHint(true)
	comp.SetOnMouseMove(m.OnMouseMove)
	comp.SetOnMouseDown(m.OnMouseDown)
	comp.SetOnMouseUp(m.OnMouseUp)
	comp.SetOnChange(func(sender lcl.IObject) {
		comp.SetChecked(false)
	})
	m.drag = newDrag(designerForm.designerBox, DsAll)
	m.drag.SetRelation(comp)
	m.object = lcl.AsWinControl(comp)
	return m
}
