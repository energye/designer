package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/message"
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
	CtForm    ComponentType = iota // 窗体
	CtWrapper                      // 包装的控件
	CtOther                        // 其它 除窗体的所有控件
)

// 设计组件
type DesigningComponent struct {
	ownerFormTab      *FormTab                // 所属设计表单面板
	originObject      any                     // 原始组件对象
	object            lcl.IWinControl         // 组件 WinControl 对象, 转换后的父类
	objectWrapper     lcl.IImage              // 组件 包裹 对象
	drag              *drag                   // 拖拽控制
	dx, dy            int32                   // 拖拽控制
	dcl, dct          int32                   // 拖拽控制
	isDown            bool                    // 拖拽控制
	propertyList      []*vtedit.TEditNodeData // 组件属性
	eventList         []*vtedit.TEditNodeData // 组件事件
	isDesigner        bool                    // 组件是否正在设计
	componentType     ComponentType           // 控件类型
	node              lcl.ITreeNode           // 组件树节点对象
	id                int                     // id 标识
	parent            *DesigningComponent     // 所属父节点
	child             []*DesigningComponent   // 拥有的子节点列表
	compPropTreeState ComponentPropTreeState  // 组件树状态
}

// 组件属性树状态
type ComponentPropTreeState struct {
	selectPropName string             // 当前选中(编辑)的属性名
	selectNode     types.PVirtualNode // 根据选中的属性名获得的节点对象
}

// 返回当前控件实例指针
func (m *DesigningComponent) Instance() uintptr {
	return m.object.Instance()
}

// 设计组件鼠标移动
func (m *DesigningComponent) OnMouseMove(sender lcl.IObject, shift types.TShiftState, X int32, Y int32) {
	br := m.object.BoundsRect()
	hint := fmt.Sprintf(`%v
	Left: %v Top: %v
	Width: %v Height: %v`, m.TreeName(), br.Left, br.Top, br.Width(), br.Height())
	m.object.SetHint(hint)
	if m.isDown {
		m.drag.Hide()
		point := m.object.ClientToParent(types.TPoint{X: X, Y: Y}, m.ownerFormTab.designerBox.object)
		x := point.X - m.dx
		y := point.Y - m.dy
		m.object.SetBounds(m.dcl+x, m.dct+y, br.Width(), br.Height())

		msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", m.dcl+x, m.dct+y, br.Width(), br.Height())
		message.Follow(msgContent)
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

		br := m.object.BoundsRect()
		msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", br.Left, br.Top, br.Width(), br.Height())
		message.Follow(msgContent)
		lcl.Mouse.SetCapture(m.object.Handle())
	}
}

// 设计组件鼠标抬起事件
func (m *DesigningComponent) OnMouseUp(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if m.isDown {
		m.drag.Show()
	}
	m.isDown = false
	message.FollowHide()
	lcl.Mouse.SetCapture(0)
}

func (m *DesigningComponent) SetObject(object any) {
	m.object = lcl.AsWinControl(object)
	m.originObject = object
	SetDesignMode(m.object)
}

// 加载组件属性到设计器
func (m *DesigningComponent) LoadPropertyToInspector() {
	// 加载到设计器
	inspector.LoadComponentProps(m)
}

// 设置组件父子关系
func (m *DesigningComponent) SetParent(parent *DesigningComponent) {
	// 设置父组件
	m.object.SetParent(parent.object)
	m.drag.SetParent(parent.object)
	m.parent = parent
	// 添加子组件
	parent.child = append(parent.child, m)
}

// 返回组件类名
func (m *DesigningComponent) Name() string {
	return m.object.Name()
}

// 返回组件树节点名
func (m *DesigningComponent) TreeName() string {
	if m.componentType == CtForm {
		return fmt.Sprintf("%v: %v", m.object.Name(), "TForm")
	} else {
		return fmt.Sprintf("%v: %v", m.object.Name(), m.object.ToString())
	}
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

// 设置组件模式为设计模式
func SetDesignMode(component lcl.IComponent) {
	lcl.DesigningComponent().SetComponentDesignMode(component, true)
	lcl.DesigningComponent().SetComponentDesignInstanceMode(component, true)
	lcl.DesigningComponent().SetComponentInlineMode(component, true)
	lcl.DesigningComponent().SetWidgetSetDesigning(component)
}

// 创建设计按钮
func NewButtonDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := new(DesigningComponent)
	m.componentType = CtOther
	m.ownerFormTab = designerForm
	comp := lcl.NewButton(designerForm.designerBox.object)
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetName(designerForm.GetComponentCaptionName("Button"))
	comp.SetCaption(comp.Name())
	comp.SetShowHint(true)
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
	comp := lcl.NewEdit(designerForm.designerBox.object)
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetName(designerForm.GetComponentCaptionName("Edit"))
	comp.SetText(comp.Name())
	comp.SetCaption(comp.Name())
	comp.SetShowHint(true)
	m.drag = newDrag(designerForm.designerBox.object, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

func NewCheckBoxDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := new(DesigningComponent)
	m.componentType = CtOther
	m.ownerFormTab = designerForm
	comp := lcl.NewCheckBox(designerForm.designerBox.object)
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetName(designerForm.GetComponentCaptionName("CheckBox"))
	comp.SetCaption(comp.Caption())
	comp.SetChecked(false)
	comp.SetShowHint(true)
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
	comp := lcl.NewPanel(designerForm.designerBox.object)
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetCaption(comp.Caption())
	comp.SetName(designerForm.GetComponentCaptionName("Panel"))
	comp.SetShowHint(true)
	m.drag = newDrag(designerForm.designerBox.object, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

func NewMainMenuDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := new(DesigningComponent)
	m.componentType = CtWrapper
	m.ownerFormTab = designerForm
	comp := lcl.NewMainMenu(designerForm.designerBox.object)
	comp.SetName(designerForm.GetComponentCaptionName("TMainMenu"))

	compWrapper := lcl.NewImage(designerForm.designerBox.object)
	compWrapper.SetLeft(x)
	compWrapper.SetTop(y)
	compWrapper.SetCursor(types.CrSize)
	compWrapper.SetShowHint(true)
	compWrapper.SetName(designerForm.GetComponentCaptionName("TMainMenuWrapper"))
	compWrapper.SetCaption(compWrapper.Name())
	m.drag = newDrag(designerForm.designerBox.object, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)

	return m
}
