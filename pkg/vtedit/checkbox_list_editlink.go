package vtedit

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// CheckBoxList

type TCheckBoxListEditLink struct {
	*TBaseEditLink
	edit      lcl.IEdit
	alignment types.TAlignment
	stopping  bool
}

func NewCheckBoxListEditLink(bindData *TEditLinkNodeData) *TCheckBoxListEditLink {
	link := new(TCheckBoxListEditLink)
	link.TBaseEditLink = NewEditLink(link)
	link.BindData = bindData
	link.Create()
	return link
}

func (m *TCheckBoxListEditLink) Create() {
	logs.Debug("TCheckBoxListEditLink CreateEdit")
	m.edit = lcl.NewEdit(nil)
	m.edit.SetVisible(false)
	m.edit.SetBorderStyle(types.BsSingle)
	m.edit.SetAutoSize(false)
	m.edit.SetDoubleBuffered(true)
	m.edit.SetReadOnly(true)
	m.edit.SetText(m.BindData.StringValue)
}

func (m *TCheckBoxListEditLink) BeginEdit() bool {
	if !m.stopping {
		m.edit.Show()
		m.edit.SelectAll()
		m.edit.SetFocus()
	}
	return true
}

func (m *TCheckBoxListEditLink) CancelEdit() bool {
	logs.Debug("TCheckBoxListEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.edit.SetVisible(false)
		m.VTree.CancelEditNode()
	}
	return true
}

func (m *TCheckBoxListEditLink) EndEdit() bool {
	value := m.edit.Text()
	logs.Debug("TCheckBoxListEditLink EndEdit Modified:", m.edit.Modified(), "value:", value, "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		m.VTree.EndEditNode()
		m.edit.Hide()
	}
	return true
}

func (m *TCheckBoxListEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	logs.Debug("TCheckBoxListEditLink PrepareEdit")
	if m.edit == nil || !m.edit.IsValid() {
		m.Create()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column
	m.edit.Font().SetColor(colors.ClWindowText)
	m.edit.SetParent(m.VTree)
	m.edit.HandleNeeded()

	if column <= -1 {
		m.edit.SetBiDiMode(m.VTree.BiDiMode())
		m.alignment = m.VTree.Alignment()
	} else {
		columns := m.VTree.Header().Columns()
		m.edit.SetBiDiMode(columns.ItemsWithColumnIndexToVirtualTreeColumn(column).BiDiMode())
		m.alignment = columns.ItemsWithColumnIndexToVirtualTreeColumn(column).Alignment()
	}

	if m.edit.BiDiMode() != types.BdLeftToRight {
		switch m.alignment {
		case types.TaLeftJustify:
			m.alignment = types.TaRightJustify
		case types.TaRightJustify:
			m.alignment = types.TaLeftJustify
		}
	}
	return true
}

func (m *TCheckBoxListEditLink) GetBounds() (R types.TRect) {
	logs.Debug("TCheckBoxListEditLink GetBounds")
	return m.edit.BoundsRect()
}

func (m *TCheckBoxListEditLink) ProcessMessage(msg *types.TLMessage) {
	logs.Debug("TCheckBoxListEditLink ProcessMessage")
	lcl.ControlHelper.WindowProc(m.edit, msg)
}

func (m *TCheckBoxListEditLink) SetBounds(R types.TRect) {
	logs.Debug("TCheckBoxListEditLink SetBounds", R)
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left
	R.Top = columnRect.Top
	R.SetHeight(columnRect.Height())
	R.SetWidth(columnRect.Width())
	m.edit.SetBoundsRect(R)

}

func (m *TCheckBoxListEditLink) Destroy(sender lcl.IObject) {
	logs.Debug("TCheckBoxListEditLink Destroy")
	m.edit.Free()
}
