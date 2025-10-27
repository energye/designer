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

type SyncLock struct {
	Point bool
	Size  bool
}

// 设计组件
type DesigningComponent struct {
	formTab       *FormTab                // 所属设计窗体
	id            int                     // id 标识
	originObject  any                     // 原始组件对象
	object        lcl.IWinControl         // 组件 对象 可视
	objectNon     lcl.IComponent          // 组件 对象 非可视
	objectNonWrap *NonVisualComponentWrap // 组件 对象 非可视, 呈现控制
	parent        *DesigningComponent     // 所属父节点
	child         []*DesigningComponent   // 拥有的子节点列表
	drag          *drag                   // 拖拽控制
	dx, dy        int32                   // 拖拽控制
	dcl, dct      int32                   // 拖拽控制
	isDown        bool                    // 拖拽控制
	propertyList  []*vtedit.TEditNodeData // 数据 组件属性
	eventList     []*vtedit.TEditNodeData // 数据 组件事件
	isDesigner    bool                    // 组件是否正在设计
	componentType ComponentType           // 组件类型
	node          lcl.ITreeNode           // 查看器 组件树节点对象
	//compPropTreeState ComponentPropTreeState    // 查看器 组件属性树状态
	page         lcl.IPageControl          // 查看器 属性页和事件页
	pageProperty lcl.ITabSheet             // 查看器 属性页
	pageEvent    lcl.ITabSheet             // 查看器 事件页
	propertyTree lcl.ILazVirtualStringTree // 查看器 组件属性树
	eventTree    lcl.ILazVirtualStringTree // 查看器 组件事件树
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

// 创建组件属性页
func (m *DesigningComponent) createComponentPropertyPage() {
	m.page = lcl.NewPageControl(inspector.componentProperty.box)
	m.page.SetParent(inspector.componentProperty.box)
	m.page.SetTabStop(true)
	m.page.SetTop(32)
	m.page.SetWidth(leftBoxWidth)
	m.page.SetHeight(inspector.componentProperty.box.Height() - m.page.Top())
	m.page.SetAlign(types.AlCustom)
	m.page.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkBottom, types.AkRight))

	m.pageProperty = lcl.NewTabSheet(m.page)
	m.pageProperty.SetParent(m.page)
	m.pageProperty.SetCaption("  属性  ")
	m.pageProperty.SetAlign(types.AlClient)
	m.pageProperty.SetBorderWidth(0)

	m.pageEvent = lcl.NewTabSheet(m.page)
	m.pageEvent.SetParent(m.page)
	m.pageEvent.SetCaption("  事件  ")
	m.pageEvent.SetAlign(types.AlClient)

	m.propertyTree = lcl.NewLazVirtualStringTree(m.pageProperty)
	vstConfig(m.propertyTree)
	m.propertyTree.SetParent(m.pageProperty)

	m.eventTree = lcl.NewLazVirtualStringTree(m.pageEvent)
	vstConfig(m.eventTree)
	m.eventTree.SetParent(m.pageEvent)

	m.initComponentPropertyTreeEvent()
}

// 创建可视组件
func newVisualComponent(designerForm *FormTab) *DesigningComponent {
	m := new(DesigningComponent)
	m.componentType = CtVisual
	m.formTab = designerForm

	m.createComponentPropertyPage()
	return m
}

// 创建非可视组件
func newNonVisualComponent(formTab *FormTab, x, y int32) *DesigningComponent {
	m := new(DesigningComponent)
	m.componentType = CtNonVisual
	m.formTab = formTab
	objectWrap := NewNonVisualComponentWrap(formTab.formRoot.object, m)
	objectWrap.SetLeftTop(x, y)
	m.objectNonWrap = objectWrap

	m.createComponentPropertyPage()
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
		point := m.ClientToParent(types.TPoint{X: X, Y: Y}, m.formTab.formRoot.object)
		x := point.X - m.dx
		y := point.Y - m.dy
		m.SetBounds(m.dcl+x, m.dct+y, br.Width(), br.Height())

		msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", m.dcl+x, m.dct+y, br.Width(), br.Height())
		message.Follow(msgContent)
		go m.UpdateNodeDataPoint(br.Left, br.Top)
	}
}

// 设计组件鼠标按下事件
func (m *DesigningComponent) OnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	logs.Debug("OnMouseDown 设计组件", m.ClassName())
	if !m.formTab.placeComponent(m, X, Y) {
		m.isDown = true
		point := m.ClientToParent(types.TPoint{X: X, Y: Y}, m.formTab.formRoot.object)
		m.dx, m.dy = point.X, point.Y
		br := m.BoundsRect()
		m.dcl = br.Left
		m.dct = br.Top
		// 更新设计查看器的属性信息
		m.formTab.hideAllDrag() // 隐藏所有 drag
		m.drag.Show()           // 显示当前设计组件 drag
		//go lcl.RunOnMainThreadAsync(func(id uint32) {
		//	m.LoadPropertyToInspector()
		//})
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
	br := m.BoundsRect()
	go m.UpdateNodeDataPoint(br.Left, br.Top)
	m.isDown = false
	message.FollowHide()
	lcl.Mouse.SetCapture(0)
	m.DragEnd()
}

// 更新节点数据, Left Top
func (m *DesigningComponent) UpdateNodeDataPoint(x, y int32) {
	var (
		top  *vtedit.TEditNodeData
		left *vtedit.TEditNodeData
	)
	for _, prop := range m.propertyList {
		switch prop.Name() {
		case "Left":
			top = prop
		case "Top":
			left = prop
		}
		if top != nil && left != nil {
			break
		}
	}
	if top != nil && left != nil {
		lcl.RunOnMainThreadAsync(func(id uint32) {
			top.SetEditValue(x)
			m.propertyTree.InvalidateNode(top.AffiliatedNode)
			left.SetEditValue(y)
			m.propertyTree.InvalidateNode(left.AffiliatedNode)
		})
	}
}

// 更新节点数据, Width Height
func (m *DesigningComponent) UpdateNodeDataSize(w, h int32) {
	var (
		width  *vtedit.TEditNodeData
		height *vtedit.TEditNodeData
	)
	for _, prop := range m.propertyList {
		switch prop.Name() {
		case "Width":
			width = prop
		case "Height":
			height = prop
		}
		if width != nil && height != nil {
			break
		}
	}
	if width != nil && height != nil {
		lcl.RunOnMainThreadAsync(func(id uint32) {
			width.SetEditValue(w)
			m.propertyTree.InvalidateNode(width.AffiliatedNode)
			height.SetEditValue(h)
			m.propertyTree.InvalidateNode(height.AffiliatedNode)
		})
	}
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
	// 设置父组件
	if m.componentType == CtNonVisual {
		m.objectNonWrap.SetParent(parent.object)
	} else {
		m.object.SetParent(parent.object)
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
