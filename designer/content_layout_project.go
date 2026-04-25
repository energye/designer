package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

type ContentLayoutProject struct {
	topBox lcl.IPanel
	title  lcl.ILabel
	box    lcl.IPanel
	// 项目栏 tree, 2个根节点 componentRoot assetsRoot
	tree lcl.ITreeView
	// 1: 组件根节点(所有窗体和组件)
	componentRoot lcl.ITreeNode
	//  2: 资源目录和文件根节点（所有代码和文件）
	assetsRoot lcl.ITreeNode
}

func initContentLayoutProject(owner *ContentLayout) *ContentLayoutProject {
	m := &ContentLayoutProject{}
	m.topBox = lcl.NewPanel(owner.projectPanel)
	m.topBox.SetBorderStyleToBorderStyle(types.BsNone)
	m.topBox.SetBevelOuter(types.BvNone)
	m.topBox.SetAlign(types.AlTop)
	m.topBox.SetHeight(30)
	m.topBox.SetParent(owner.projectPanel)

	title := lcl.NewLabel(m.topBox)
	title.SetCaption("项目管理器")
	title.SetLeft(5)
	title.SetTop(5)
	font := title.Font()
	font.SetSize(10)
	font.SetStyle(types.NewSet(types.FsBold))
	title.SetParent(m.topBox)

	m.box = lcl.NewPanel(owner.projectPanel)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetDoubleBuffered(true)
	m.box.SetAlign(types.AlClient)
	m.box.SetParent(owner.projectPanel)

	m.tree = lcl.NewTreeView(owner.projectPanel)
	m.tree.SetAutoExpand(true)
	m.tree.SetReadOnly(true)
	m.tree.SetDoubleBuffered(true)
	m.tree.SetTreeLineColor(colors.RGBToColor(128, 128, 128))
	m.tree.SetTreeLinePenStyle(types.PsSolid)
	m.tree.SetAlign(types.AlClient)
	m.tree.SetVisible(false)
	m.tree.SetBorderStyleToBorderStyle(types.BsNone)
	m.tree.Font().SetHeight(-10)
	//form.tree.SetImages(imageComponents.ImageList50())
	m.tree.SetImages(imageComponents.ImageList100())
	//m.tree.SetMultiSelect(true) // 多选控制
	//m.tree.SetOnGetSelectedIndex(form.TreeOnGetSelectedIndex)
	//m.tree.SetOnMouseDown(form.TreeOnMouseDown)
	//m.tree.SetOnContextPopup(form.TreeOnContextPopup)
	//m.CreateComponentMenu()
	//m.tree.SetPopupMenu(form.componentMenu.treePopupMenu)
	//m.tree.SetParent(m.box)

	return m
}
