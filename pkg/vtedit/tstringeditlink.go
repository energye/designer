package vtedit

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"log"
)

type IStringEditLink interface {
	lcl.ICustomVTEditLink
}

type TStringEditLink struct {
	lcl.ICustomVTEditLink
	edit lcl.IEdit
}

func NewStringEditLink() IStringEditLink {
	m := new(TStringEditLink)
	m.ICustomVTEditLink = lcl.NewCustomVTEditLink()
	m.ICustomVTEditLink.SetOnBeginEdit(m.BeginEdit)
	m.ICustomVTEditLink.SetOnCancelEdit(m.CancelEdit)
	m.ICustomVTEditLink.SetOnEndEdit(m.EndEdit)
	m.ICustomVTEditLink.SetOnPrepareEdit(m.PrepareEdit)
	m.ICustomVTEditLink.SetOnGetBounds(m.GetBounds)
	m.ICustomVTEditLink.SetOnProcessMessage(m.ProcessMessage)
	m.ICustomVTEditLink.SetOnSetBounds(m.SetBounds)
	m.ICustomVTEditLink.SetOnDestroy(m.Destroy)
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
