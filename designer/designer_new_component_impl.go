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

// 组件实现函数

// 组件类型
type ComponentType int32

const (
	CtForm      ComponentType = iota // 窗体
	CtNonVisual                      // 非可视组件
	CtVisual                         // 可视组件
)

// 设计组件
type DesigningComponent struct {
	ownerFormTab      *FormTab                // 所属设计表单面板
	originObject      any                     // 原始组件对象
	designerBox       lcl.IWinControl         // 设计窗体组件 对象 可视
	object            lcl.IWinControl         // 组件 对象 可视
	objectNon         lcl.IComponent          // 组件 对象 非可视
	objectNonWrap     *NonVisualComponentWrap // 组件 对象 非可视, 呈现控制
	drag              *drag                   // 拖拽控制
	dx, dy            int32                   // 拖拽控制
	dcl, dct          int32                   // 拖拽控制
	isDown            bool                    // 拖拽控制
	propertyList      []*vtedit.TEditNodeData // 组件属性
	eventList         []*vtedit.TEditNodeData // 组件事件
	isDesigner        bool                    // 组件是否正在设计
	componentType     ComponentType           // 组件类型
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

// 设置组件模式为设计模式
func SetDesignMode(component lcl.IComponent) {
	lcl.DesigningComponent().SetComponentDesignMode(component, true)
	lcl.DesigningComponent().SetComponentDesignInstanceMode(component, true)
	lcl.DesigningComponent().SetComponentInlineMode(component, true)
	lcl.DesigningComponent().SetWidgetSetDesigning(component)
}

// 创建可视组件
func newVisualComponent(designerForm *FormTab) *DesigningComponent {
	m := new(DesigningComponent)
	m.componentType = CtVisual
	m.ownerFormTab = designerForm
	return m
}

// 创建非可视组件
func newNonVisualComponent(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := new(DesigningComponent)
	m.componentType = CtNonVisual
	m.ownerFormTab = designerForm
	objectWrap := NewNonVisualComponentWrap(designerForm.designerBox.object, m)
	objectWrap.SetLeftTop(x, y)
	m.objectNonWrap = objectWrap
	return m
}

// 设置基础通用属性
func setBaseProp(comp lcl.IControl, x, y int32) {
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetCaption(comp.Name())
	comp.SetShowHint(true)
}

// 返回当前组件实例指针
func (m *DesigningComponent) Instance() uintptr {
	if m.componentType == CtNonVisual {
		return m.objectNonWrap.Instance()
	} else {
		return m.object.Instance()
	}
}
func (m *DesigningComponent) SetHint(hint string) {
	if m.componentType == CtNonVisual {
		m.objectNonWrap.SetHint(hint)
	} else {
		m.object.SetHint(hint)
	}
}

func (m *DesigningComponent) ClassName() string {
	if m.componentType == CtNonVisual {
		return m.objectNon.ToString()
	} else {
		return m.object.ToString()
	}
}

func (m *DesigningComponent) BoundsRect() types.TRect {
	if m.componentType == CtNonVisual {
		return m.objectNonWrap.BoundsRect()
	} else {
		return m.object.BoundsRect()
	}
}

func (m *DesigningComponent) SetBounds(x, y, w, h int32) {
	if m.componentType == CtNonVisual {
		m.objectNonWrap.SetLeftTop(x, y)
	} else {
		m.object.SetBounds(x, y, w, h)
	}
}

func (m *DesigningComponent) ClientToParent(point types.TPoint, parent lcl.IWinControl) types.TPoint {
	if m.componentType == CtNonVisual {
		return m.objectNonWrap.ClientToParent(point, parent)
	}
	return m.object.ClientToParent(point, parent)
}

// 设计组件鼠标移动
func (m *DesigningComponent) OnMouseMove(sender lcl.IObject, shift types.TShiftState, X int32, Y int32) {
	br := m.BoundsRect()
	hint := fmt.Sprintf(`%v
	Left: %v Top: %v
	Width: %v Height: %v`, m.TreeName(), br.Left, br.Top, br.Width(), br.Height())
	m.SetHint(hint)
	if m.isDown {
		m.drag.Hide()
		point := m.ClientToParent(types.TPoint{X: X, Y: Y}, m.ownerFormTab.designerBox.object)
		x := point.X - m.dx
		y := point.Y - m.dy
		m.SetBounds(m.dcl+x, m.dct+y, br.Width(), br.Height())

		msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", m.dcl+x, m.dct+y, br.Width(), br.Height())
		message.Follow(msgContent)

	}
}

// 设计组件鼠标按下事件
func (m *DesigningComponent) OnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	logs.Debug("OnMouseDown 设计组件", m.ClassName())
	if !m.ownerFormTab.placeComponent(m, X, Y) {
		m.isDown = true
		point := m.ClientToParent(types.TPoint{X: X, Y: Y}, m.ownerFormTab.designerBox.object)
		m.dx, m.dy = point.X, point.Y
		br := m.BoundsRect()
		m.dcl = br.Left
		m.dct = br.Top
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
		msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", br.Left, br.Top, br.Width(), br.Height())
		message.Follow(msgContent)
		if m.object != nil {
			lcl.Mouse.SetCapture(m.object.Handle())
		}
		m.DragBegin()
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
	m.DragEnd()
}

// 设置对象实例
func (m *DesigningComponent) SetObject(object any) {
	if m.componentType == CtNonVisual {
		m.objectNon = lcl.AsComponent(object)
		//SetDesignMode(m.objectNonWrap.wrap)
		SetDesignMode(m.objectNonWrap.icon) // 使用icon
	} else {
		m.object = lcl.AsWinControl(object)
		SetDesignMode(m.object)
	}
	m.originObject = object
}

// 加载组件属性到设计器
func (m *DesigningComponent) LoadPropertyToInspector() {
	// 加载到设计器
	inspector.LoadComponentProps(m)
}

// 设置组件父子关系
func (m *DesigningComponent) SetParent(parent *DesigningComponent) {
	var parentObject lcl.IWinControl
	if parent.componentType == CtForm {
		parentObject = parent.designerBox
	} else {
		parentObject = parent.object
	}
	// 设置父组件
	if m.componentType == CtNonVisual {
		m.objectNonWrap.SetParent(parentObject)
	} else {
		m.object.SetParent(parentObject)
	}
	m.parent = parent
	// 添加子组件
	parent.child = append(parent.child, m)
}

// 返回组件类名
func (m *DesigningComponent) Name() string {
	if m.componentType == CtNonVisual {
		return m.objectNon.Name()
	} else {
		return m.object.Name()
	}
}

// 返回组件树节点名
func (m *DesigningComponent) TreeName() string {
	return fmt.Sprintf("%v: %v", m.Name(), m.ClassName())
}

// 返回组件树节点使用的图标索引
func (m *DesigningComponent) IconIndex() int32 {
	name := m.ClassName() + ".png"
	return imageComponents.ImageIndex(name)
}

// 返回真实对象
func (m *DesigningComponent) Object() lcl.IObject {
	if m.componentType == CtNonVisual {
		return m.objectNon
	}
	return m.object
}

// 获取当前组件对象属性
func (m *DesigningComponent) GetProps() {
	// 属性列表为空时获取属性列表
	if m.propertyList == nil {

		properties := lcl.DesigningComponent().GetComponentProperties(m.Object())
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
			newEditNodeData.Build()
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

// 拖拽开始调用
func (m *DesigningComponent) DragBegin() {
	if m.componentType == CtNonVisual {
		m.objectNonWrap.TextFollowHide()
	}
}

// 拖拽结束调用，或创建后调用
func (m *DesigningComponent) DragEnd() {
	if m.componentType == CtNonVisual {
		m.objectNonWrap.TextFollowShow()
	}
}
