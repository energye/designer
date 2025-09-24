package vtedit

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"log"
)

// CheckBoxList

type TCheckBoxListEditLink struct {
	*TBaseEditLink
	edit      lcl.IEdit
	bounds    types.TRect
	value     bool
	alignment types.TAlignment
	stopping  bool
}

func NewCheckBoxListEditLink(bindData *TNodeData) *TCheckBoxListEditLink {
	m := new(TCheckBoxListEditLink)
	m.TBaseEditLink = NewEditLink(m)
	m.BindData = bindData
	m.Create()
	return m
}

func (m *TCheckBoxListEditLink) Create() {
	log.Println("TCheckBoxListEditLink CreateEdit")
	m.edit = lcl.NewEdit(nil)
	m.edit.SetVisible(false)
	m.edit.SetBorderStyle(types.BsSingle)
	m.edit.SetAutoSize(false)
	m.edit.SetDoubleBuffered(true)
	m.edit.SetReadOnly(true)
}

func (m *TCheckBoxListEditLink) AddChild() {

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
	log.Println("TCheckBoxListEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.edit.SetVisible(false)
		m.VTree.CancelEditNode()
	}
	return true
}

func (m *TCheckBoxListEditLink) EndEdit() bool {
	value := m.edit.Text()
	log.Println("TCheckBoxListEditLink EndEdit Modified:", m.edit.Modified(), "value:", value, "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		m.VTree.EndEditNode()
		m.edit.Hide()
	}
	return true
}

func (m *TCheckBoxListEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	log.Println("TCheckBoxListEditLink PrepareEdit")
	if m.edit == nil || !m.edit.IsValid() {
		m.Create()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column
	m.edit.Font().SetColor(colors.ClWindowText)
	m.edit.SetParent(m.VTree)
	m.edit.HandleNeeded()

	m.AddChild()

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
	log.Println("TCheckBoxListEditLink GetBounds")
	return m.edit.BoundsRect()
}

func (m *TCheckBoxListEditLink) ProcessMessage(msg *types.TLMessage) {
	log.Println("TCheckBoxListEditLink ProcessMessage")
	lcl.ControlHelper.WindowProc(m.edit, msg)
}

func (m *TCheckBoxListEditLink) SetBounds(R types.TRect) {
	log.Println("TCheckBoxListEditLink SetBounds", R)
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left
	R.Top = columnRect.Top
	R.SetHeight(columnRect.Height())
	R.SetWidth(columnRect.Width())
	m.edit.SetBoundsRect(R)

}

func (m *TCheckBoxListEditLink) Destroy(sender lcl.IObject) {
	log.Println("TCheckBoxListEditLink Destroy")
	m.edit.Free()
}
