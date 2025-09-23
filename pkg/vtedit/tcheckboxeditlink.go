package vtedit

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// CheckBox

type TCheckBoxEditLink struct {
	*TBaseEditLink
	edit       lcl.ICheckBox
	textBounds types.TRect
	text       string
	alignment  types.TAlignment
	stopping   bool
}

func NewCheckBoxEditLink() *TCheckBoxEditLink {
	m := new(TCheckBoxEditLink)
	m.TBaseEditLink = NewEditLink(m)
	return m
}

func (m *TCheckBoxEditLink) BeginEdit() bool {
	return false
}

func (m *TCheckBoxEditLink) CancelEdit() bool {
	return false
}

func (m *TCheckBoxEditLink) EndEdit() bool {
	return false
}

func (m *TCheckBoxEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	m.VTree = tree
	m.Node = node
	m.Column = column
	return false
}

func (m *TCheckBoxEditLink) GetBounds() (R types.TRect) {
	return
}

func (m *TCheckBoxEditLink) ProcessMessage(msg *types.TLMessage) {
}

func (m *TCheckBoxEditLink) SetBounds(R types.TRect) {
}

func (m *TCheckBoxEditLink) Destroy(sender lcl.IObject) {
}
