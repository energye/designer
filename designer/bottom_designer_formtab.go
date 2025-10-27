package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 设计器面板

// 设计表单的 tab
type FormTab struct {
	id                   int                 // 索引, 关联 forms key: index
	name                 string              // 窗体名称
	scroll               lcl.IScrollBox      // 外 滚动条
	isDesigner           bool                // 当前窗体Form是否正在设计
	sheet                lcl.ITabSheet       // tab sheet
	isDown, isUp, isMove bool                // 鼠标事件
	componentName        map[string]int      // 组件分类名, 同类组件ID序号
	treePopupMenu        lcl.IPopupMenu      // 组件树右键菜单
	formDesigner         *TEngFormDesigner   // 设计器处理器
	formRoot             *DesigningComponent // 设计器, 窗体 Form, 组件树的根节点
	tree                 lcl.ITreeView       // 组件树
}

func (m *FormTab) IsDuplicateName(currComp *DesigningComponent, name string) bool {
	if m.formRoot != currComp && m.formRoot.Name() == name {
		return true
	}
	var iterable func(comp *DesigningComponent) bool
	iterable = func(comp *DesigningComponent) bool {
		if comp != currComp && comp.Name() == name {
			return true
		}
		for _, comp := range comp.child {
			if iterable(comp) {
				return true
			}
		}
		return false
	}
	return iterable(m.formRoot)
}

// 添加设计组件到组件列表
func (m *FormTab) AddComponentToList(component *DesigningComponent) {
	m.formDesigner.AddComponentToList(component)
}

// 返回设计组件
func (m *FormTab) GetComponentFormList(instance uintptr) *DesigningComponent {
	return m.formDesigner.GetComponentFormList(instance)
}

// 删除一个设计组件
func (m *FormTab) RemoveComponentFormList(instance uintptr) {
	m.formDesigner.RemoveComponentFormList(instance)
}

// 切换组件设计
func (m *FormTab) switchComponentEditing(targetComp *DesigningComponent) {
	// 隐藏之前设计的组件
	// 隐藏之前设计组件的属性和事件列表
	var iterable func(comp *DesigningComponent)
	iterable = func(comp *DesigningComponent) {
		if comp.isDesigner {
			comp.drag.Hide()
			comp.page.Hide()
		}
		for _, comp := range comp.child {
			iterable(comp)
		}
	}
	iterable(m.formRoot)

	// 显示当前设计组件 drag
	targetComp.drag.Show()
	// 显示当前设计组件的属性和事件列表
	targetComp.page.Show()
}

//
//func (m *FormTab) show

// 放置设计组件到设计面板或父组件容器
func (m *FormTab) placeComponent(owner *DesigningComponent, x, y int32) bool {
	// 放置设计组件
	isAcceptsControl := false
	if owner.object != nil {
		isAcceptsControl = owner.object.ControlStyle().In(types.CsAcceptsControls)
	}
	if toolbar.selectComponent != nil && isAcceptsControl {
		logs.Debug("选中设计组件:", toolbar.selectComponent.index, toolbar.selectComponent.name)
		m.formRoot.drag.Hide()
		// 获取注册的组件创建函数
		if create := GetRegisterComponent(toolbar.selectComponent.name); create != nil {
			// 创建设计组件
			newComp := create(m, x, y)
			newComp.SetParent(owner)
			newComp.formTab.switchComponentEditing(newComp)
			newComp.DragEnd()
			// 1. 加载属性到设计器
			// 此步骤会初始化并填充设计组件实例
			newComp.LoadPropertyToInspector()
			// 2. 添加到组件树
			go lcl.RunOnMainThreadAsync(func(id uint32) {
				owner.AddChild(newComp)
			})
		} else {
			logs.Warn("选中设计组件", toolbar.selectComponent.name, "未实现或未注册")
		}
		// 重置工具栏选项卡上的组件工具按钮按下
		toolbar.ResetTabComponentDown()
		return true
	}
	return false
}

// 窗体设计界面 鼠标移动
func (m *FormTab) designerOnMouseMove(sender lcl.IObject, shift types.TShiftState, x, y int32) {
	br := m.formRoot.BoundsRect()
	hint := fmt.Sprintf(`%v: TForm
	Left: %v Top: %v
	Width: %v Height: %v`, m.formRoot.Name(), br.Left, br.Top, br.Width(), br.Height())
	m.formRoot.SetHint(hint)
}

// 窗体设计界面 鼠标按下, 放置设计组件, 加载组件属性
func (m *FormTab) designerOnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
	// 创建组件
	logs.Debug("鼠标点击设计器")
	if !m.placeComponent(m.formRoot, x, y) {
		m.switchComponentEditing(m.formRoot)
		logs.Debug("加载窗体属性")
		// 设置选中状态到设计器组件树
		m.formRoot.SetSelected()
		//lcl.Mouse.SetCapture(m.formRoot.object.Handle())
	}
}

// 窗体设计界面 鼠标抬起
func (m *FormTab) designerOnMouseUp(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
	//lcl.Mouse.SetCapture(0)
}

func (m *FormTab) onHide(sender lcl.IObject) {
	logs.Debug("Designer PageControl FormTab Hide")
	// 非设计状态
	m.isDesigner = false
	m.tree.SetVisible(false)
}

func (m *FormTab) onShow(sender lcl.IObject) {
	logs.Debug("Designer PageControl FormTab Show")
	// 设计状态
	m.isDesigner = true
	m.tree.SetVisible(true)

	// 加载设计组件
	// 默认窗体表单
	defaultComp := m.formRoot
	var iterable func(comp *DesigningComponent) bool
	iterable = func(comp *DesigningComponent) bool {
		if comp == nil {
			return false
		}
		// 如果有当前设计面板有正在设计的组件
		// 加载正在设计的组件
		if comp.isDesigner {
			defaultComp = comp
			return true
		}
		for _, comp := range comp.child {
			if iterable(comp) {
				return true
			}
		}
		return false
	}
	iterable(m.formRoot)

	logs.Debug("Current Designer Component")
	// 加载组件属性
	defaultComp.LoadPropertyToInspector()
}

// 获取组件名 Caption
func (m *FormTab) GetComponentCaptionName(component string) string {
	if c, ok := m.componentName[component]; ok {
		m.componentName[component] = c + 1
	} else {
		m.componentName[component] = 1
	}
	component = fmt.Sprintf("%v%d", component, m.componentName[component])
	return component
}

func (m *FormTab) designerOnPaint(control lcl.ICustomControl) {
	//control.SetOnPaint(func(sender lcl.IObject) {
	//	// 绘制网格
	//	m.drawGrid(control)
	//})
}

// 绘制风格线
func (m *FormTab) drawGrid(control lcl.ICustomControl) {
	//logs.Debug("drawGrid")
	gridSize := 9 // 小刻度
	formRoot := control
	canvas := formRoot.Canvas()
	canvas.PenToPen().SetColor(colors.ClBlack)
	width, height := formRoot.Width(), formRoot.Height()
	for i := 1; i < int(width)/gridSize; i++ {
		x := int32(i * gridSize)
		for j := 1; j < int(height)/gridSize; j++ {
			y := int32(j * gridSize)
			canvas.SetPixels(x, y, colors.ClBlack)
		}
	}
}

// 添加窗体表单根节点
func (m *FormTab) AddFormNode() {
	// 窗体 根节点
	m.tree.BeginUpdate()
	defer m.tree.EndUpdate()
	items := m.tree.Items()
	m.formRoot.id = nextTreeDataId()
	newNode := items.AddChild(nil, m.formRoot.TreeName())
	newNode.SetImageIndex(m.formRoot.IconIndex())    // 显示图标索引
	newNode.SetSelectedIndex(m.formRoot.IconIndex()) // 选中图标索引
	newNode.SetSelected(true)
	newNode.SetData(m.formRoot.instance())
	m.formRoot.node = newNode
	// 添加到设计组件列表
	m.AddComponentToList(m.formRoot)
}

// 添加组件节点
func (m *FormTab) AddComponentNode(parent, child *DesigningComponent) {
	if parent == nil {
		logs.Error("添加组件节点失败, 父节点为空")
		return
	} else if child == nil {
		logs.Error("添加组件节点失败, 子节点为空")
		return
	}
	if child.componentType == CtVisual || child.componentType == CtNonVisual {
		m.tree.BeginUpdate()
		defer m.tree.EndUpdate()
		items := m.tree.Items()
		// 组件 子节点
		child.id = nextTreeDataId()
		node := items.AddChild(parent.node, child.TreeName())
		child.node = node
		node.SetImageIndex(child.IconIndex())    // 显示图标索引
		node.SetSelectedIndex(child.IconIndex()) // 选中图标索引
		node.SetSelected(true)
		node.SetData(child.instance())
		// 添加到设计组件列表
		m.AddComponentToList(child)
	} else {
		logs.Error("添加组件节点失败, 子节点非组件节点")
	}
}
