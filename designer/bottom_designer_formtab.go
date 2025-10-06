package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"unsafe"
)

// 设计器面板

// 设计表单的 tab
type FormTab struct {
	id                   int                 // 索引, 关联 forms key: index
	name                 string              // 窗体名称
	scroll               lcl.IScrollBox      // 外 滚动条
	isDesigner           bool                // 是否正在设计
	sheet                lcl.ITabSheet       // tab sheet
	designerBox          *DesigningComponent // 设计器, 模拟 TForm, 也是组件树的根节点
	form                 *DesigningComponent // 设计器的窗体, 用于获取属性列表
	isDown, isUp, isMove bool                // 鼠标事件
	componentName        map[string]int      // 组件分类名, 同类组件ID序号
	//componentList        []*DesigningComponent // 设计器的组件列表
	tree lcl.ITreeView // 组件树
}

func (m *FormTab) IsDuplicateName(name string) bool {
	if m.designerBox.Name() == name {
		return true
	}
	var iterable func(comp *DesigningComponent)
	iterable = func(comp *DesigningComponent) {

	}
	for _, comp := range m.designerBox.child {
		iterable(comp)
	}
	return false
}

// 数据指针转设计组件
func (m *FormTab) DataToDesigningComponent(data uintptr) *DesigningComponent {
	dc := (*DesigningComponent)(unsafe.Pointer(data))
	return dc
}

func (m *FormTab) TreeOnGetSelectedIndex(sender lcl.IObject, node lcl.ITreeNode) {
	data := node.Data()
	component := m.DataToDesigningComponent(data)
	if component != nil {
		component.ownerFormTab.hideAllDrag() // 隐藏所有 drag
		component.drag.Show()                // 显示当前设计组件 drag
		go lcl.RunOnMainThreadAsync(func(id uint32) {
			component.LoadPropertyToInspector()
		})
	}
	logs.Info("Inspector-component-tree OnGetSelectedIndex name:", node.Text(), "id:", component.id)
}

//func (m *FormTab) addDesignerComponent(component *DesigningComponent) {
//	m.componentList = append(m.componentList, component)
//}

// 隐藏所有控件的 drag
func (m *FormTab) hideAllDrag() {
	var iterable func(comp *DesigningComponent)
	iterable = func(comp *DesigningComponent) {
		comp.drag.Hide()
		for _, comp := range comp.child {
			iterable(comp)
		}
	}
	iterable(m.designerBox)
}

// 放置设计组件到设计面板或父组件容器
func (m *FormTab) placeComponent(owner *DesigningComponent, x, y int32) bool {
	// 放置设计组件
	if toolbar.selectComponent != nil && !config.ContainerDenyList.IsDeny(owner.object.ToString()) {
		logs.Debug("选中设计组件:", toolbar.selectComponent.index, toolbar.selectComponent.name)
		m.designerBox.drag.Hide()
		// 获取注册的组件创建函数
		if create := GetRegisterComponent(toolbar.selectComponent.name); create != nil {
			// 创建设计组件
			newComp := create(m, x, y)
			newComp.SetParent(owner)
			// 隐藏所有 drag
			newComp.ownerFormTab.hideAllDrag()
			// 显示当前设计组件 drag
			newComp.drag.Show()
			// 1. 加载属性到设计器
			// 此步骤会初始化并填充设计组件实例
			newComp.LoadPropertyToInspector()
			// 2. 添加到组件树
			go lcl.RunOnMainThreadAsync(func(id uint32) {
				//m.componentTree.Load(component)
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

// 窗体设计界面 鼠标按下, 放置设计控件, 加载控件属性
func (m *FormTab) designerOnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
	// 创建组件
	logs.Debug("鼠标点击设计器")
	if !m.placeComponent(m.designerBox, x, y) {
		m.hideAllDrag()
		m.designerBox.drag.Show()
		logs.Debug("加载窗体属性")
		// 加载属性列表到设计器组件属性
		inspector.LoadComponentProps(m.designerBox)
		// 设置选中状态到设计器组件树
		m.designerBox.SetSelected()
	}

}

func (m *FormTab) onHide(sender lcl.IObject) {
	logs.Debug("Designer PageControl FormTab Hide", m.name)
	// 非设计状态
	m.isDesigner = false
	m.tree.SetVisible(false)
}

func (m *FormTab) onShow(sender lcl.IObject) {
	logs.Debug("Designer PageControl FormTab Show", m.name)
	// 设计状态
	m.isDesigner = true
	m.tree.SetVisible(true)

	// 加载设计组件
	// 默认窗体表单
	defaultComp := m.designerBox
	var iterable func(comp *DesigningComponent) bool
	iterable = func(comp *DesigningComponent) bool {
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
	iterable(m.designerBox)

	logs.Debug("Current Designer Component", "Name:", m.name)
	// 加载组件属性
	inspector.LoadComponentProps(defaultComp)
}

// 窗体设计界面 鼠标抬起
func (m *FormTab) designerOnMouseUp(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {

}

// 窗体设计界面 鼠标移动
func (m *FormTab) designerOnMouseMove(sender lcl.IObject, shift types.TShiftState, x, y int32) {
	//width, height := m.designerBox.Width(), m.designerBox.Height()
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

func (m *FormTab) designerOnPaint(sender lcl.IObject) {
	// 绘制网格
	m.drawGrid()
}

// 绘制风格线
func (m *FormTab) drawGrid() {
	gridSize := 9 // 小刻度
	designerBox := m.designerBox.object.(lcl.IPanel)
	canvas := designerBox.Canvas()
	canvas.PenToPen().SetColor(colors.ClBlack)
	width, height := designerBox.Width(), designerBox.Height()
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
	m.designerBox.id = nextTreeDataId()
	newNode := items.AddChild(nil, m.form.TreeName())
	newNode.SetImageIndex(m.form.IconIndex())    // 显示图标索引
	newNode.SetSelectedIndex(m.form.IconIndex()) // 选中图标索引
	newNode.SetSelected(true)
	newNode.SetData(m.designerBox.instance())
	m.designerBox.node = newNode
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
	if child.componentType == CtOther {
		m.tree.BeginUpdate()
		defer m.tree.EndUpdate()
		items := m.tree.Items()
		// 控件 子节点
		child.id = nextTreeDataId()
		//child.parent = parent
		node := items.AddChild(parent.node, child.TreeName())
		child.node = node
		node.SetImageIndex(child.IconIndex())    // 显示图标索引
		node.SetSelectedIndex(child.IconIndex()) // 选中图标索引
		node.SetSelected(true)
		node.SetData(child.instance())
		//parent.child = append(parent.child, child)
	} else {
		logs.Error("添加组件节点失败, 子节点非组件节点")
	}
}
