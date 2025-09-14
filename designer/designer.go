package designer

// 创建设计器布局
func (m *TAppWindow) createDesignerLayout() {
	// 顶部菜单
	m.createMenu()
	// 工具栏
	m.createTopToolbar()
	// 底部布局
	m.createBottomBox()
}
