package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

type ContentLayoutInspector struct {
	searchEdit lcl.ITreeFilterEdit // 组件搜索框
	topBox     lcl.IPanel
	title      lcl.ILabel
}

func initContentLayoutInspector(owner *ContentLayout) *ContentLayoutInspector {
	m := &ContentLayoutInspector{}
	m.searchEdit = lcl.NewTreeFilterEdit(owner.inspectorPanel)
	m.searchEdit.SetTextHint("搜索属性")
	m.searchEdit.SetAlign(types.AlTop)
	m.searchEdit.SetAutoSelect(false)
	borderSpacing := m.searchEdit.BorderSpacing()
	borderSpacing.SetLeft(3)
	borderSpacing.SetRight(3)
	borderSpacing.SetTop(3)
	borderSpacing.SetBottom(3)
	m.searchEdit.SetParent(owner.inspectorPanel)

	m.topBox = lcl.NewPanel(owner.inspectorPanel)
	m.topBox.SetBorderStyleToBorderStyle(types.BsNone)
	m.topBox.SetBevelOuter(types.BvNone)
	m.topBox.SetAlign(types.AlTop)
	m.topBox.SetHeight(30)
	borderSpacing = m.topBox.BorderSpacing()
	m.topBox.SetParent(owner.inspectorPanel)

	title := lcl.NewLabel(m.topBox)
	title.SetCaption("属性检查器")
	title.SetLeft(5)
	title.SetTop(5)
	font := title.Font()
	font.SetSize(10)
	font.SetStyle(types.NewSet(types.FsBold))
	title.SetParent(m.topBox)

	return m
}
