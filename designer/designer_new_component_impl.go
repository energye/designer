// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package designer

import (
	"fmt"
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/vtedit"
	uiBean "github.com/energye/designer/uigen/bean"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"sort"
	"strings"
)

// 设计组件实现

// 设计组件
type TDesigningComponent struct {
	FormTab        *FormTab                  // 所属设计窗体
	id             int                       // id 标识
	originObject   any                       // 原始组件对象
	object         lcl.IWinControl           // 组件 对象 可视
	objectNon      lcl.IComponent            // 组件 对象 非可视
	objectNonWrap  *TNonVisualComponentWrap  // 组件 对象 非可视, 呈现控制
	parent         *TDesigningComponent      // 所属父节点
	Child          []*TDesigningComponent    // 拥有的子节点列表
	drag           *drag                     // 拖拽控制
	PropertyList   []*vtedit.TEditNodeData   // 数据 组件属性
	EventList      []*vtedit.TEditNodeData   // 数据 组件事件
	isDesigner     bool                      // 组件是否正在设计
	ComponentType  consts.ComponentType      // 组件类型
	node           lcl.ITreeNode             // 查看器 组件树节点对象
	page           lcl.IPageControl          // 查看器 属性页和事件页
	pageProperty   lcl.ITabSheet             // 查看器 属性页
	pageEvent      lcl.ITabSheet             // 查看器 事件页
	propertyTree   lcl.ILazVirtualStringTree // 查看器 组件属性树
	eventTree      lcl.ILazVirtualStringTree // 查看器 组件事件树
	isLoadProperty bool                      // 是否加载完成属性到属性列表
	isLoadEvent    bool                      // 是否加载完成属性到事件列表
	mod            string                    // 组件所属依赖模块, lcl, cef, wv
}

// 组件属性树状态
type ComponentPropTreeState struct {
	selectPropName string             // 当前选中(编辑)的属性名
	selectNode     types.PVirtualNode // 根据选中的属性名获得的节点对象
}

// 设置组件模式为设计模式
func SetDesignMode(component lcl.IComponent) {
	lcl.DesigningComponent().SetComponentDesignMode(component, true)
	lcl.DesigningComponent().SetComponentInlineMode(component, true)
	lcl.DesigningComponent().SetComponentDesignInstanceMode(component, true)
	lcl.DesigningComponent().SetWidgetSetDesigning(component)
}

func SetComponentDesignMode(component uintptr) {
	api.SetComponentDesignMode(component, true)
	api.SetComponentInlineMode(component, true)
	api.SetComponentDesignInstanceMode(component, true)
	api.SetWidgetSetDesigning(component)
}

func (m *TDesigningComponent) Free() {
	//if m.ComponentType == consts.CtNonVisual {
	//	//m.objectNonWrap.Free()
	//} else {
	//	//m.object.Free()
	//}

	m.FormTab = nil
	m.parent = nil
	for _, child := range m.Child {
		child.Free()
	}
	m.Child = nil
	if m.drag != nil {
		m.drag.Free()
	}
	for _, data := range m.PropertyList {
		data.Free()
	}
	for _, data := range m.EventList {
		data.Free()
	}
	m.PropertyList = nil
	m.EventList = nil
	if m.page != nil && m.page.IsValid() {
		m.page.Free()
	}
}

// 初始化组件属性页控件
// 该函数确保组件属性页只被初始化一次，创建属性页、事件页以及对应的树形控件
func (m *TDesigningComponent) mustComponentPropertyPage() {
	if m.page != nil {
		return
	}
	m.page = lcl.NewPageControl(inspector.componentProperty.propComponentProp)
	m.page.SetTabStop(true)
	m.page.SetAlign(types.AlClient)
	SetComponentDefaultColor(m.page)
	m.page.SetVisible(false)
	m.page.SetParent(inspector.componentProperty.propComponentProp)

	m.pageProperty = lcl.NewTabSheet(m.page)
	m.pageProperty.SetParent(m.page)
	m.pageProperty.SetCaption("  属性  ")
	m.pageProperty.SetAlign(types.AlClient)
	SetComponentDefaultColor(m.pageProperty)
	m.pageProperty.SetBorderWidth(0)

	m.pageEvent = lcl.NewTabSheet(m.page)
	m.pageEvent.SetParent(m.page)
	m.pageEvent.SetCaption("  事件  ")
	SetComponentDefaultColor(m.pageEvent)
	m.pageEvent.SetAlign(types.AlClient)

	m.propertyTree = lcl.NewLazVirtualStringTree(m.pageProperty)
	vstConfig(m.propertyTree)
	m.propertyTree.SetParent(m.pageProperty)

	m.eventTree = lcl.NewLazVirtualStringTree(m.pageEvent)
	vstConfig(m.eventTree)
	m.eventTree.SetParent(m.pageEvent)
	// 初始化组件属性树事件
	m.initComponentPropertyTreeEvent()
	// 初始化组件事件树事件
	m.initComponentEventTreeEvent()
}

// 创建可视组件
func newVisualComponent(designerForm *FormTab) *TDesigningComponent {
	m := new(TDesigningComponent)
	m.ComponentType = consts.CtVisual
	m.FormTab = designerForm
	return m
}

// 创建非可视组件
func newNonVisualComponent(formTab *FormTab, x, y int32) *TDesigningComponent {
	m := new(TDesigningComponent)
	m.ComponentType = consts.CtNonVisual
	m.FormTab = formTab
	objectWrap := NewNonVisualComponentWrap(formTab.FormRoot.object, m)
	objectWrap.SetLeftTop(x, y)
	m.objectNonWrap = objectWrap
	return m
}

// 设置基础通用属性
func setBaseProp(comp lcl.IControl, x, y int32) {
	if x != 0 {
		comp.SetLeft(x)
	}
	if y != 0 {
		comp.SetTop(y)
	}
	if api.IsObjectInstanceOf(comp.Instance(), lcl.TCustomEditClass()) {
		customEdit := lcl.AsCustomEdit(comp)
		customEdit.SetReadOnly(true)
	}
}

// 返回当前组件实例指针
func (m *TDesigningComponent) Instance() uintptr {
	return m.Component().Instance()
}
func (m *TDesigningComponent) SetHint(hint string) {
	m.Control().SetHint(hint)
}

func (m *TDesigningComponent) ClassName() string {
	return m.Component().ToString()
}

func (m *TDesigningComponent) BoundsRect() types.TRect {
	return m.Control().BoundsRect()
}

func (m *TDesigningComponent) SetBounds(x, y, w, h int32) {
	control := m.Control()
	if m.ComponentType == consts.CtNonVisual {
		control.SetBounds(x, y, nonWrapW, nonWrapH)
	} else {
		control.SetBounds(x, y, w, h)
	}
}

func (m *TDesigningComponent) ClientToParent(point types.TPoint, parent lcl.IWinControl) types.TPoint {
	control := m.Control()
	return control.ClientToParent(point, parent)
}

// 设计组件鼠标移动
func (m *TDesigningComponent) OnMouseMove(sender lcl.IObject, shift types.TShiftState, X int32, Y int32) {
	m.drag.OnMouseMove(m, shift, X, Y)
}

// 设计组件鼠标按下事件
func (m *TDesigningComponent) OnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if button == types.MbRight {
		m.SetSelected()
	} else {
		m.drag.OnMouseDown(m, button, shift, X, Y)
	}
}

// 设计组件鼠标抬起事件
func (m *TDesigningComponent) OnMouseUp(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if button == types.MbRight && m.ComponentType != consts.CtForm {
		cursorPos := lcl.Mouse.CursorPos()
		m.FormTab.componentMenu.treePopupMenu.PopUpWithIntX2(cursorPos.X, cursorPos.Y)
	} else {
		m.drag.OnMouseUp(m, button, shift, X, Y)
	}
}

// 更新节点数据, Left Top
func (m *TDesigningComponent) UpdateNodeDataPoint(x, y int32) {
	if m.page == nil {
		return
	}
	var (
		top  *vtedit.TEditNodeData
		left *vtedit.TEditNodeData
	)
	for _, prop := range m.PropertyList {
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
		// 更新坐标
		triggerUIGeneration(m, nil, event.CodeGenUI)
		lcl.RunOnMainThreadAsync(func(id uint32) {
			// 同步更新属性列表
			top.SetEditValue(x)
			m.propertyTree.InvalidateNode(top.AffiliatedNode)
			left.SetEditValue(y)
			m.propertyTree.InvalidateNode(left.AffiliatedNode)
		})
	}
}

// 更新节点数据, Width Height
func (m *TDesigningComponent) UpdateNodeDataSize(w, h int32) {
	var (
		width  *vtedit.TEditNodeData
		height *vtedit.TEditNodeData
	)
	for _, prop := range m.PropertyList {
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
		// 更新宽高
		triggerUIGeneration(m, nil, event.CodeGenUI)
		lcl.RunOnMainThreadAsync(func(id uint32) {
			width.SetEditValue(w)
			m.propertyTree.InvalidateNode(width.AffiliatedNode)
			height.SetEditValue(h)
			m.propertyTree.InvalidateNode(height.AffiliatedNode)
		})
	}
}

// 设置对象实例
func (m *TDesigningComponent) SetObject(object any) {
	if m.ComponentType == consts.CtNonVisual {
		m.objectNon = lcl.AsComponent(object)
	} else {
		m.object = lcl.AsWinControl(object)
	}
	m.originObject = object
}

// 加载组件属性到设计器
func (m *TDesigningComponent) LoadPropertyToInspector() {
	// 加载到设计器
	if m == nil {
		logs.Error("加载组件属性/事件失败, 设计组件为空")
		return
	}
	// 属性列表为空时获取属性列表
	m.GetProps()
	// 加载属性列表
	m.loadPropertyList()
	// 加载事件列表
	m.loadEventList()
	logs.Debug("加载组件属性完成", m.ClassName())

}

// 设置组件父子关系
// 只设计组件才会调用此函数
func (m *TDesigningComponent) SetParent(parent *TDesigningComponent) {
	// 设置父组件
	control := parent.object
	if m.ComponentType == consts.CtNonVisual {
		m.objectNonWrap.SetParent(control)
	} else {
		m.object.SetParent(control)
	}
	m.parent = parent
	// 添加子组件
	parent.Child = append(parent.Child, m)
}

// 返回组件类名
func (m *TDesigningComponent) Name() string {
	name := m.Component().Name()
	return name
}

func (m *TDesigningComponent) SetName(name string) {
	m.Component().SetName(name)
}

// 返回组件树节点名
func (m *TDesigningComponent) TreeName() string {
	return fmt.Sprintf("%v: %v", m.Name(), m.ClassName())
}

// 返回组件树节点使用的图标索引
func (m *TDesigningComponent) IconIndex() int32 {
	name := m.ClassName() + ".png"
	return imageComponents.ImageIndex(name)
}

func (m *TDesigningComponent) Control() lcl.IControl {
	if m.ComponentType == consts.CtNonVisual {
		return m.objectNonWrap.icon
	} else if m.ComponentType == consts.CtForm {
		return m.originObject.(*TDesignerForm)
	}
	return m.object
}

func (m *TDesigningComponent) Component() lcl.IComponent {
	if m.ComponentType == consts.CtNonVisual {
		return m.objectNon
	} else if m.ComponentType == consts.CtForm {
		return m.originObject.(*TDesignerForm)
	}
	return m.object
}

func (m *TDesigningComponent) SetVisible(v bool) {
	if control := m.Control(); control != nil {
		control.SetVisible(v)
	}
}

// 获取当前组件对象属性
func (m *TDesigningComponent) GetProps() {
	// 属性列表为空时获取属性列表
	if m.PropertyList == nil {
		methods := tool.GetObjectMethodNames(m.originObject)
		if methods == nil {
			logs.Error("获取当前组件对象属性错误, 获取对象方法列表为空, 组件名:", m.Name())
		}
		properties := lcl.DesigningComponent().GetComponentProperties(m.Component())
		logs.Debug("LoadComponent Count:", len(properties), "ClassName:", m.ClassName())
		// 拆分 属性和事件
		var (
			eventList       []*vtedit.TEditNodeData
			propertyList    []*vtedit.TEditNodeData
			configCompProp  = config.ComponentProperty
			tempPropertyMap = make(map[string]struct{}) // 用于下面判断
		)
		for _, prop := range properties {
			metadata := prop
			tool.FixPropInfo(methods, &metadata)

			if customMetadatas := configCompProp.GetCustomPropertyList(metadata.Name); customMetadatas != nil {
				if len(customMetadatas) == 1 {
					// 数组只有一个元素时, 将它做为元数据使用, 替换当前属性
					customMetadata := customMetadatas[0]
					// Options 默认值设置
					if customMetadata.Options == "" {
						// 未配置options时, 使用源数据默认值
						customMetadata.Options = metadata.Options
					}
					metadata = customMetadata
				}
			}

			// 创建节点数据
			newEditLinkNodeData := vtedit.NewEditLinkNodeData(&metadata)
			// 节点数据包装
			newEditNodeData := &vtedit.TEditNodeData{
				EditNodeData:        newEditLinkNodeData,
				OriginNodeData:      newEditLinkNodeData.Clone(), // 克隆一份原始数据
				AffiliatedComponent: m,
			}

			if metadata.Kind == consts.TkMethod {
				// tkMethod 事件函数
				eventList = append(eventList, newEditNodeData)
			} else {
				// 其它侧为属性
				propertyList = append(propertyList, newEditNodeData)
			}
			// 构建节点数据
			newEditNodeData.Build()
			tempPropertyMap[newEditNodeData.Name()] = struct{}{}
		}

		// 通用属性 包含属性, 如果重复忽略配置的
		for _, prop := range configCompProp.Include() {
			if _, ok := tempPropertyMap[prop.Name]; ok {
				// 忽略配置属性
				continue
			}
			// 创建节点数据
			customProperty := vtedit.NewEditLinkNodeData(&prop)
			// 节点数据包装
			newEditNodeData := &vtedit.TEditNodeData{EditNodeData: customProperty,
				OriginNodeData:      customProperty.Clone(), // 克隆一份原始数据
				AffiliatedComponent: m}
			propertyList = append(propertyList, newEditNodeData) // 添加到组件属性
			// 构建节点数据
			newEditNodeData.Build()
		}

		// 排序
		sort.Slice(eventList, func(i, j int) bool {
			return strings.ToLower(eventList[i].EditNodeData.Name) < strings.ToLower(eventList[j].EditNodeData.Name)
		})
		sort.Slice(propertyList, func(i, j int) bool {
			return strings.ToLower(propertyList[i].EditNodeData.Name) < strings.ToLower(propertyList[j].EditNodeData.Name)
		})
		m.EventList = eventList
		m.PropertyList = propertyList
	}
}

// 拖拽开始调用
func (m *TDesigningComponent) DragBegin() {
	if m.ComponentType == consts.CtNonVisual {
		m.objectNonWrap.TextFollowHide()
	}
}

// 拖拽结束调用，或创建后调用
func (m *TDesigningComponent) DragEnd() {
	if m.ComponentType == consts.CtNonVisual {
		m.objectNonWrap.TextFollowShow()
	}
}

func (m *TDesigningComponent) Parent() *TDesigningComponent {
	return m.parent
}

// 删除当前节点
func (m *TDesigningComponent) Remove() {
	// 父组件删除
	if m.parent != nil {
		idx := m.Index()
		if idx != -1 {
			m.parent.Child = append(m.parent.Child[:idx], m.parent.Child[idx+1:]...)
		}
	}
	// 组件树删除
	m.node.Delete()
	// 设计窗体删除
	m.Free()
}

// MoveTo 将当前设计组件移动到指定的目标组件位置
//
// 参数:
//
//	destination - 目标设计组件，表示要移动到的位置
//	mode - 节点附加模式，控制移动的具体行为和位置关系
func (m *TDesigningComponent) MoveTo(destination *TDesigningComponent, mode types.TNodeAttachMode) {
	// 从当前父组件中完全移除自己
	removeSelf := func() {
		if m.parent != nil {
			idx := m.Index()
			if idx != -1 {
				m.parent.Child = append(m.parent.Child[:idx], m.parent.Child[idx+1:]...)
			}
		}
	}

	switch mode {
	case types.NaAddChild: // 添加到最后
		if destination.parent != nil {
			removeSelf()
			m.parent = destination.parent
			destination.parent.Child = append(destination.parent.Child, m)
		}
	case types.NaAddFirst: // 添加到最前面
		if destination.parent != nil {
			removeSelf()
			m.parent = destination.parent
			destination.parent.Child = append([]*TDesigningComponent{m}, destination.parent.Child...)
		}
	case types.NaInsert: // 插入到目标节点前面
		if destination.parent != nil {
			targetIndex := destination.Index()
			if targetIndex != -1 && destination != m {
				removeSelf()
				m.parent = destination.parent
				m.parent.Child = append(
					m.parent.Child[:targetIndex],
					append([]*TDesigningComponent{m}, m.parent.Child[targetIndex:]...)...,
				)
			}
		}
	case types.NaInsertBehind: // 插入到目标节点后面
		if destination.parent != nil {
			targetIndex := destination.Index()
			if targetIndex != -1 && destination != m {
				removeSelf()
				m.parent = destination.parent
				insertIndex := targetIndex + 1
				if insertIndex > len(m.parent.Child) {
					insertIndex = len(m.parent.Child)
				}
				m.parent.Child = append(
					m.parent.Child[:insertIndex],
					append([]*TDesigningComponent{m}, m.parent.Child[insertIndex:]...)...,
				)
			}
		}
	}
}

func (m *TDesigningComponent) LastChild() *TDesigningComponent {
	size := len(m.Child)
	if size > 0 {
		return m.Child[size-1]
	}
	return nil
}

func (m *TDesigningComponent) FirstChild() *TDesigningComponent {
	size := len(m.Child)
	if size > 0 {
		return m.Child[0]
	}
	return nil
}

func (m *TDesigningComponent) Index() int {
	if m.parent != nil {
		for i, comp := range m.parent.Child {
			if m == comp {
				return i
			}
		}
	}
	return -1
}

func (m *TDesigningComponent) NextSibling() *TDesigningComponent {
	if m.parent != nil {
		idx := m.Index()
		if idx != -1 && idx < len(m.parent.Child)-1 {
			return m.parent.Child[idx+1]
		}
	}
	return nil
}

func (m *TDesigningComponent) PrevSibling() *TDesigningComponent {
	if m.parent != nil {
		idx := m.Index()
		if idx != -1 && idx > 0 {
			return m.parent.Child[idx-1]
		}
	}
	return nil
}

// 查找属性节点, 根据属性名路径查找属性节点数据
// namePaths: 属性名路径, [Font, Style] [Header, Font, Style]
func (m *TDesigningComponent) FindNodeDataByNamePaths(property uiBean.TProperty) (result *vtedit.TEditNodeData) {
	namePaths := tool.Split(property.Name, ".") // 属性名路径 Font.Style
	if len(namePaths) == 0 || m.PropertyList == nil || m.EventList == nil {
		return nil
	}
	propName := namePaths[0]
	// 查找当前属性节点数据
	if property.Type == consts.PdtMethod {
		for _, data := range m.EventList {
			if data.Name() == propName {
				result = data
				break
			}
		}
	} else {
		for _, data := range m.PropertyList {
			if data.Name() == propName {
				result = data
				break
			}
		}
	}
	if result != nil && len(namePaths) > 1 {
		var namePathsIndex = 1
		var iterator func(node *vtedit.TEditNodeData)
		iterator = func(node *vtedit.TEditNodeData) {
			for _, data := range node.Child {
				if namePathsIndex < len(namePaths) && data.Name() == namePaths[namePathsIndex] {
					result = data
					namePathsIndex++
				}
				if data.Type() == consts.PdtClass {
					iterator(data)
				}
			}
		}
		iterator(result)
	}
	return
}
