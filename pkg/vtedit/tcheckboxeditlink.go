package vtedit

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"log"
	"strconv"
)

// CheckBox

type TCheckBoxEditLink struct {
	*TBaseEditLink
	parent    lcl.IWinControl
	edit      lcl.ICheckBox
	bounds    types.TRect
	value     bool
	alignment types.TAlignment
	stopping  bool
}

func NewCheckBoxEditLink(parent lcl.IWinControl) *TCheckBoxEditLink {
	m := new(TCheckBoxEditLink)
	m.TBaseEditLink = NewEditLink(m)
	m.parent = parent
	m.CreateEdit()
	return m
}

func (m *TCheckBoxEditLink) CreateEdit() {
	log.Println("TCheckBoxEditLink CreateEdit")
	m.edit = lcl.NewCheckBox(m.parent)
	m.edit.SetParent(m.parent)
	m.edit.SetVisible(false)
	m.edit.SetCaption("(False)")
	m.edit.SetDoubleBuffered(true)
	m.edit.SetParentColor(false)
}
func (m *TCheckBoxEditLink) BeginEdit() bool {
	if !m.stopping {
		m.edit.SetVisible(true)
	}
	return true
}

func (m *TCheckBoxEditLink) CancelEdit() bool {
	log.Println("TCheckBoxEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.edit.SetVisible(false)
		m.VTree.CancelEditNode()
	}
	return true
}

func (m *TCheckBoxEditLink) EndEdit() bool {
	value := m.edit.Checked()
	log.Println("TCheckBoxEditLink EndEdit", "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		if m.OnNewData != nil {
			m.OnNewData(m.Node, m.Column, strconv.FormatBool(value))
		}
		m.VTree.EndEditNode()
		m.edit.SetVisible(false)
	}
	return true
}

func (m *TCheckBoxEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	log.Println("TCheckBoxEditLink PrepareEdit")
	if m.edit == nil || m.edit.IsValid() {
		m.CreateEdit()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column
	value := m.VTree.Text(m.Node, m.Column)
	if v, err := strconv.ParseBool(value); err == nil {
		m.value = v
	}
	// 节点的初始大小、字体和文本。
	log.Println("  PrepareEdit GetTextInfo:", m.bounds, m.value)
	//m.edit.SetParent(m.VTree)
	//m.edit.HandleNeeded()
	m.edit.SetChecked(m.value)
	return true
}

func (m *TCheckBoxEditLink) GetBounds() (R types.TRect) {
	log.Println("TCheckBoxEditLink GetBounds")
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R = columnRect
	return
}

func (m *TCheckBoxEditLink) ProcessMessage(msg *types.TLMessage) {
	log.Println("TCheckBoxEditLink ProcessMessage")
	//lcl.ControlHelper.WindowProc(m.edit, msg)
}

func (m *TCheckBoxEditLink) SetBounds(R types.TRect) {
	log.Println("TCheckBoxEditLink SetBounds", R)
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left + 5
	R.Top = columnRect.Top
	//R.SetHeight(columnRect.Height())
	//R.SetWidth(columnRect.Width())
	m.edit.SetBoundsRect(R)

}

func (m *TCheckBoxEditLink) Destroy(sender lcl.IObject) {
	log.Println("TCheckBoxEditLink Destroy")
	m.edit.Free()
}
