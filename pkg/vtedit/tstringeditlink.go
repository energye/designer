package vtedit

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"log"
)

type IStringEditLink interface {
	lcl.IVTEditLink
	lcl.ICustomVTEditLink
}

type TStringEditLink struct {
	link lcl.ICustomVTEditLink
	edit lcl.IEdit
}

func NewStringEditLink() *TStringEditLink {
	m := new(TStringEditLink)
	m.link = lcl.NewCustomVTEditLink()
	m.link.SetOnBeginEdit(m.BeginEdit)
	m.link.SetOnCancelEdit(m.CancelEdit)
	m.link.SetOnEndEdit(m.EndEdit)
	m.link.SetOnPrepareEdit(m.PrepareEdit)
	m.link.SetOnGetBounds(m.GetBounds)
	m.link.SetOnProcessMessage(m.ProcessMessage)
	m.link.SetOnSetBounds(m.SetBounds)
	m.link.SetOnDestroy(m.Destroy)
	m.edit = lcl.NewEdit(nil)
	return m
}

func (m *TStringEditLink) BeginEdit() bool {
	log.Println("BeginEdit")
	return false
}

func (m *TStringEditLink) CancelEdit() bool {
	log.Println("CancelEdit")
	return false
}

func (m *TStringEditLink) EndEdit() bool {
	log.Println("EndEdit")
	return false
}

func (m *TStringEditLink) PrepareEdit(tree lcl.IBaseVirtualTree, node types.PVirtualNode, column int32) bool {
	log.Println("PrepareEdit")

	return false
}

func (m *TStringEditLink) GetBounds() types.TRect {
	log.Println("GetBounds")
	return types.TRect{}
}

func (m *TStringEditLink) ProcessMessage(msg *types.TLMessage) {
	log.Println("ProcessMessage")
	lcl.ControlHelper.WindowProc(m.edit, msg)
}

func (m *TStringEditLink) SetBounds(R types.TRect) {
	log.Println("SetBounds")
}

func (m *TStringEditLink) Destroy(sender lcl.IObject) {
	log.Println("Destroy")
}
