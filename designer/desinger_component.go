package designer

import (
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 设计中的组件

type DesigningComponent struct {
	owner  *FormTab
	object lcl.IWinControl
	drag   *drag
}

func (m *DesigningComponent) OnMouseMove(sender lcl.IObject, shift types.TShiftState, X int32, Y int32) {
	hint := fmt.Sprintf(`%v: %v
	Left: %v Top: %v
	Width: %v Height: %v`, m.object.Caption(), m.object.ToString(), m.object.Left(), m.object.Top(), m.object.Width(), m.object.Height())
	m.object.SetHint(hint)
}

func (m *DesigningComponent) OnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	m.owner.hideAllDrag()
	m.drag.Show()
}

func (m *DesigningComponent) OnMouseUp(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {

}

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
