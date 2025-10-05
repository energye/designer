package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"sort"
	"strings"
)

// 组件设计创建管理

// 组件类型
type ComponentType int32

const (
	CtForm  ComponentType = iota // 窗体
	CtOther                      // 其它 除窗体的所有控件
)

// 设计组件
type DesigningComponent struct {
	ownerFormTab  *FormTab                // 所属设计表单面板
	originObject  any                     // 原始组件对象
	object        lcl.IWinControl         // 组件 WinControl 对象, 转换后的父类
	drag          *drag                   // 拖拽控制
	dx, dy        int32                   // 拖拽控制
	dcl, dct      int32                   // 拖拽控制
	isDown        bool                    // 拖拽控制
	propertyList  []*vtedit.TEditNodeData // 组件属性
	eventList     []*vtedit.TEditNodeData // 组件事件
	isDesigner    bool                    // 组件是否正在设计
	componentType ComponentType           // 控件类型
	node          lcl.ITreeNode           // 组件树节点对象
	id            int                     // id 标识
	parent        *DesigningComponent     // 所属父节点
	child         []*DesigningComponent   // 拥有的子节点列表
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
		point := m.object.ClientToParent(types.TPoint{X: X, Y: Y}, m.ownerFormTab.designerBox.object)
		x := point.X - m.dx
		y := point.Y - m.dy
		m.object.SetBounds(m.dcl+x, m.dct+y, br.Width(), br.Height())
	}
}

// 设计组件鼠标按下事件
func (m *DesigningComponent) OnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	logs.Debug("OnMouseDown 设计组件", m.object.ToString())
	if !m.ownerFormTab.placeComponent(m, X, Y) {
		m.isDown = true
		point := m.object.ClientToParent(types.TPoint{X: X, Y: Y}, m.ownerFormTab.designerBox.object)
		m.dx, m.dy = point.X, point.Y
		m.dcl = m.object.Left()
		m.dct = m.object.Top()
		// 更新设计查看器的属性信息
		m.ownerFormTab.hideAllDrag() // 隐藏所有 drag
		m.drag.Show()                // 显示当前设计组件 drag
		go lcl.RunOnMainThreadAsync(func(id uint32) {
			m.LoadPropertyToInspector()
		})
		// 更新设计查看器的组件树信息
		go lcl.RunOnMainThreadAsync(func(id uint32) {
			// 设置选中状态
			m.SetSelected()
		})
	}
}

// 设计组件鼠标抬起事件
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
	// 加载到设计器
	inspector.LoadComponentProps(m)
}

func (m *DesigningComponent) SetParent(parent *DesigningComponent) {
	m.object.SetParent(parent.object)
	m.drag.SetParent(parent.object)
	m.parent = parent
	parent.child = append(parent.child, m)
}

// 返回组件类名
func (m *DesigningComponent) Name() string {
	return m.object.Name()
}

// 返回组件树节点名
func (m *DesigningComponent) TreeName() string {
	return fmt.Sprintf("%v: %v", m.object.Name(), m.object.ToString())
}

// 返回组件树节点使用的图标索引
func (m *DesigningComponent) IconIndex() int32 {
	return CompTreeIcon(m.object.ToString())
}

// 获取当前组件对象属性
func (m *DesigningComponent) GetProps() {
	// 属性列表为空时获取属性列表
	if m.propertyList == nil {
		properties := lcl.DesigningComponent().GetComponentProperties(m.object)
		logs.Debug("LoadComponent Count:", len(properties))
		// 拆分 属性和事件
		var (
			eventList    []*vtedit.TEditNodeData
			propertyList []*vtedit.TEditNodeData
		)
		for _, prop := range properties {
			newProp := prop
			newEditLinkNodeData := vtedit.NewEditLinkNodeData(&newProp)
			newEditNodeData := &vtedit.TEditNodeData{EditNodeData: newEditLinkNodeData, OriginNodeData: newEditLinkNodeData.Clone(), AffiliatedComponent: m}
			if newProp.Kind == "tkMethod" {
				// tkMethod 事件函数
				eventList = append(eventList, newEditNodeData)
			} else {
				// 其它侧为属性
				propertyList = append(propertyList, newEditNodeData)
			}
			//logs.Debug("  ", toJSON(prop))
		}
		// 排序
		sort.Slice(eventList, func(i, j int) bool {
			return strings.ToLower(eventList[i].EditNodeData.Name) < strings.ToLower(eventList[j].EditNodeData.Name)
		})
		sort.Slice(propertyList, func(i, j int) bool {
			return strings.ToLower(propertyList[i].EditNodeData.Name) < strings.ToLower(propertyList[j].EditNodeData.Name)
		})
		m.eventList = eventList
		m.propertyList = propertyList
	}
}

// 创建设计窗体-隐藏
func NewFormDesigner(designerForm *FormTab) *DesigningComponent {
	m := new(DesigningComponent)
	m.componentType = CtForm
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
	m.componentType = CtOther
	m.ownerFormTab = designerForm
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
	m.componentType = CtOther
	m.ownerFormTab = designerForm
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
	m.componentType = CtOther
	m.ownerFormTab = designerForm
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
	m.componentType = CtOther
	m.ownerFormTab = designerForm
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
