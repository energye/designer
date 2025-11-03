// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package vtedit

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

var MainForm lcl.IForm

// Laz 虚拟树动态创建组件

// IBaseEditLink 基础接口，需被实现
type IBaseEditLink interface {
	lcl.IObject
	BeginEdit() bool
	CancelEdit() bool
	EndEdit() bool
	PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool
	GetBounds() types.TRect
	ProcessMessage(msg *types.TLMessage)
	SetBounds(R types.TRect)
	Destroy(sender lcl.IObject)
}

// TBaseEditLink 基础对象，被嵌套继承
type TBaseEditLink struct {
	lcl.TObject
	baseEditLink lcl.ICustomVTEditLink
	self         IBaseEditLink
	BindData     *TEditNodeData
	VTree        lcl.ILazVirtualStringTree
	Node         types.PVirtualNode
	Column       int32
}

func (m *TBaseEditLink) AsIVTEditLink() lcl.IVTEditLink {
	return lcl.AsVTEditLink(m.baseEditLink.AsIntfVTEditLink())
}

// NewEditLink 基础实现, 需实现 IBaseEditLink 接口
func NewEditLink(self IBaseEditLink) *TBaseEditLink {
	m := new(TBaseEditLink)
	m.self = self
	m.baseEditLink = lcl.NewCustomVTEditLink()
	m.baseEditLink.SetOnBeginEdit(m.OnBeginEdit)
	m.baseEditLink.SetOnCancelEdit(m.OnCancelEdit)
	m.baseEditLink.SetOnEndEdit(m.OnEndEdit)
	m.baseEditLink.SetOnPrepareEdit(m.OnPrepareEdit)
	m.baseEditLink.SetOnGetBounds(m.OnGetBounds)
	m.baseEditLink.SetOnProcessMessage(m.OnProcessMessage)
	m.baseEditLink.SetOnSetBounds(m.OnSetBounds)
	m.baseEditLink.SetOnDestroy(m.OnDestroy)
	return m
}

func (m *TBaseEditLink) OnBeginEdit() bool {
	if m.self != nil {
		return m.self.BeginEdit()
	}
	return false
}

func (m *TBaseEditLink) OnCancelEdit() bool {
	if m.self != nil {
		return m.self.CancelEdit()
	}
	return false
}

func (m *TBaseEditLink) OnEndEdit() bool {
	if m.self != nil {
		return m.self.EndEdit()
	}
	return false
}

func (m *TBaseEditLink) OnPrepareEdit(tree lcl.IBaseVirtualTree, node types.PVirtualNode, column int32) bool {
	if m.self != nil {
		return m.self.PrepareEdit(lcl.AsLazVirtualStringTree(tree), node, column)
	}
	return false
}

func (m *TBaseEditLink) OnGetBounds() (R types.TRect) {
	if m.self != nil {
		R = m.self.GetBounds()
	}
	return
}

func (m *TBaseEditLink) OnProcessMessage(msg *types.TLMessage) {
	if m.self != nil {
		m.self.ProcessMessage(msg)
	}
}

func (m *TBaseEditLink) OnSetBounds(R types.TRect) {
	if m.self != nil {
		m.self.SetBounds(R)
	}
}

func (m *TBaseEditLink) OnDestroy(sender lcl.IObject) {
	if m.self != nil {
		m.self.Destroy(sender)
	}
}
