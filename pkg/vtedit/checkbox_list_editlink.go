package vtedit

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"log"
	"strconv"
)

// CheckBoxList

type TCheckBoxListEditLink struct {
	*TBaseEditLink
	checkbox  lcl.ICheckBox
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
	m.checkbox = lcl.NewCheckBox(nil)
	m.checkbox.SetVisible(false)
	m.checkbox.SetCaption("(false)")
	m.checkbox.SetDoubleBuffered(true)
	m.checkbox.SetOnChange(func(sender lcl.IObject) {
		m.checkbox.SetCaption("(" + strconv.FormatBool(m.checkbox.Checked()) + ")")
	})
}

func (m *TCheckBoxListEditLink) BeginEdit() bool {
	if !m.stopping {
		m.checkbox.SetVisible(true)
	}
	return true
}

func (m *TCheckBoxListEditLink) CancelEdit() bool {
	log.Println("TCheckBoxListEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.checkbox.SetVisible(false)
		m.VTree.CancelEditNode()
	}
	return true
}

func (m *TCheckBoxListEditLink) EndEdit() bool {
	value := m.checkbox.Checked()
	log.Println("TCheckBoxListEditLink EndEdit", "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		m.BindData.BoolValue = value
		m.VTree.EndEditNode()
		m.checkbox.SetVisible(false)
	}
	return true
}

func (m *TCheckBoxListEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	log.Println("TCheckBoxListEditLink PrepareEdit")
	if m.checkbox == nil || m.checkbox.IsValid() {
		m.Create()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column
	value := m.VTree.Text(m.Node, m.Column)
	if v, err := strconv.ParseBool(value); err == nil {
		m.value = v
	}
	m.checkbox.SetParent(m.VTree)
	m.checkbox.HandleNeeded()
	m.checkbox.SetChecked(m.value)
	return true
}

func (m *TCheckBoxListEditLink) GetBounds() (R types.TRect) {
	log.Println("TCheckBoxListEditLink GetBounds")
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R = columnRect
	return
}

func (m *TCheckBoxListEditLink) ProcessMessage(msg *types.TLMessage) {
	log.Println("TCheckBoxListEditLink ProcessMessage")
	//lcl.ControlHelper.WindowProc(m.checkbox, msg)
}

func (m *TCheckBoxListEditLink) SetBounds(R types.TRect) {
	log.Println("TCheckBoxListEditLink SetBounds", R)
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left + 5
	R.Top = columnRect.Top
	m.checkbox.SetBoundsRect(R)

}

func (m *TCheckBoxListEditLink) Destroy(sender lcl.IObject) {
	log.Println("TCheckBoxListEditLink Destroy")
	m.checkbox.Free()
}
