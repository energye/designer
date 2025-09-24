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
	checkbox  lcl.ICheckBox
	bounds    types.TRect
	value     bool
	alignment types.TAlignment
	stopping  bool
}

func NewCheckBoxEditLink(bindData *TNodeData) *TCheckBoxEditLink {
	m := new(TCheckBoxEditLink)
	m.TBaseEditLink = NewEditLink(m)
	m.BindData = bindData
	m.Create()
	return m
}

func (m *TCheckBoxEditLink) Create() {
	log.Println("TCheckBoxEditLink CreateEdit")
	m.checkbox = lcl.NewCheckBox(nil)
	m.checkbox.SetVisible(false)
	m.checkbox.SetCaption("(false)")
	m.checkbox.SetDoubleBuffered(true)
	m.checkbox.SetOnChange(func(sender lcl.IObject) {
		m.checkbox.SetCaption("(" + strconv.FormatBool(m.checkbox.Checked()) + ")")
	})
}

func (m *TCheckBoxEditLink) BeginEdit() bool {
	if !m.stopping {
		m.checkbox.SetVisible(true)
	}
	return true
}

func (m *TCheckBoxEditLink) CancelEdit() bool {
	log.Println("TCheckBoxEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.checkbox.SetVisible(false)
		m.VTree.CancelEditNode()
	}
	return true
}

func (m *TCheckBoxEditLink) EndEdit() bool {
	value := m.checkbox.Checked()
	log.Println("TCheckBoxEditLink EndEdit", "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		m.BindData.BoolValue = value
		m.VTree.EndEditNode()
		m.checkbox.SetVisible(false)
	}
	return true
}

func (m *TCheckBoxEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	log.Println("TCheckBoxEditLink PrepareEdit")
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

func (m *TCheckBoxEditLink) GetBounds() (R types.TRect) {
	log.Println("TCheckBoxEditLink GetBounds")
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R = columnRect
	return
}

func (m *TCheckBoxEditLink) ProcessMessage(msg *types.TLMessage) {
	log.Println("TCheckBoxEditLink ProcessMessage")
	lcl.ControlHelper.WindowProc(m.checkbox, msg)
}

func (m *TCheckBoxEditLink) SetBounds(R types.TRect) {
	log.Println("TCheckBoxEditLink SetBounds", R)
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left + 5
	R.Top = columnRect.Top + 3
	//R.SetHeight(columnRect.Height())
	//R.SetWidth(columnRect.Width())
	m.checkbox.SetBoundsRect(R)

}

func (m *TCheckBoxEditLink) Destroy(sender lcl.IObject) {
	log.Println("TCheckBoxEditLink Destroy")
	m.checkbox.Free()
}
