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

package project

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"sync"
)

// NewConfigProjectForm 创建一个新的项目创建表单实例
// 该函数初始化一个 TConfigProjectForm 结构体，并通过 lcl.Application.NewForm 方法将其注册为应用程序窗体
func NewConfigProjectForm() *TConfigProjectForm {
	newForm := &TConfigProjectForm{}
	lcl.Application.NewForm(newForm)
	return newForm
}

type TConfigProjectForm struct {
	lcl.TEngForm
	oldWndPrc   uintptr
	closing     bool
	goVersionOK bool
	one         sync.Once
	box         lcl.IPanel
	selectDir   lcl.ISelectDirectoryDialog
}

func (m *TConfigProjectForm) FormCreate(sender lcl.IObject) {
	logs.Debug("TConfigProjectForm FormCreate")
	m.SetCaption("应用配置")
	m.SetWidth(formWidth)
	m.SetHeight(formHeight)
	m.SetVisible(false)
	m.SetDoubleBuffered(true)
	constr := m.Constraints()
	constr.SetMaxWidth(formWidth)
	constr.SetMaxHeight(formHeight)
	constr.SetMinWidth(formWidth)
	constr.SetMinHeight(formHeight)
	m.SetFormStyle(types.FsStayOnTop)
	m.SetShowInTaskBar(types.StNever)
	//m.SetColor(bgColor)
	m.SetBorderIcons(types.NewSet(types.BiSystemMenu))
	m.WorkAreaCenter()
	m.box = lcl.NewPanel(m)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetAlign(types.AlClient)
	m.box.SetParent(m)
	m.SetOnShow(m.onShow)
	m.initComponents()
	//m._HookWndProcMessage()
}

func (m *TConfigProjectForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	m.closing = true
}

func (m *TConfigProjectForm) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
}

func (m *TConfigProjectForm) initComponents() {
}

// 窗口显示事件
func (m *TConfigProjectForm) onShow(sender lcl.IObject) {
	logs.Debug("TConfigProjectForm Show")

}
