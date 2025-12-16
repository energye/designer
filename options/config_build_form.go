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

package options

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/widget/wg"
)

var (
	buildFormWidth  = int32(555)
	buildFormHeight = int32(555)
)

func NewBuildForm() *TBuildForm {
	newEngForm := lcl.NewEngForm(nil)
	newForm := &TBuildForm{TEngForm: newEngForm.(*lcl.TEngForm)}
	newForm.FormCreate(newEngForm)
	newForm.SetOnCloseQuery(newForm.OnCloseQuery)
	newForm.SetOnClose(newForm.OnClose)
	return newForm
}

type TBuildForm struct {
	*lcl.TEngForm
	closing   bool
	font      lcl.IFont
	selectDir lcl.ISelectDirectoryDialog

	// 操作按钮
	cancelBtn *wg.TButton
	saveBtn   *wg.TButton
}

func (m *TBuildForm) FormCreate(sender lcl.IObject) {
	logs.Debug("TBuildForm FormCreate")
	fontSize := int32(12)
	if tool.IsLinux {
		fontSize = 10
	}
	m.SetCaption("构建配置")
	m.SetWidth(buildFormWidth)
	m.SetHeight(buildFormHeight)
	constr := m.Constraints()
	constr.SetMaxWidth(buildFormWidth)
	constr.SetMaxHeight(buildFormHeight)
	constr.SetMinWidth(buildFormWidth)
	constr.SetMinHeight(buildFormHeight)
	m.SetVisible(false)
	m.SetDoubleBuffered(true)
	m.SetBorderIcons(types.NewSet(types.BiSystemMenu))
	m.WorkAreaCenter()
	m.font = lcl.NewFont()
	m.font.SetName("微软雅黑")
	m.font.SetSize(fontSize)
	m.SetColor(colors.ClWhite)

	m.selectDir = lcl.NewSelectDirectoryDialog(m)
	{

	}

	{
		m.cancelBtn = wg.NewButton(m)
		m.cancelBtn.SetText("关　闭")
		m.cancelBtn.SetFont(m.font)
		m.cancelBtn.Font().SetColor(colors.ClWhite)
		m.cancelBtn.SetRadius(3)
		cancelBtnRect := types.TRect{Left: 315, Top: buildFormHeight - 45}
		cancelBtnRect.SetWidth(100)
		cancelBtnRect.SetHeight(35)
		m.cancelBtn.SetBoundsRect(cancelBtnRect)
		m.cancelBtn.SetColor(colors.RGBToColor(255, 127, 127))
		m.cancelBtn.SetParent(m)
		m.cancelBtn.SetOnClick(m.closeClick)

		m.saveBtn = wg.NewButton(m)
		m.saveBtn.SetText("保　存")
		m.saveBtn.SetFont(m.font)
		m.saveBtn.Font().SetColor(colors.ClWhite)
		m.saveBtn.SetRadius(3)
		saveBtnRect := types.TRect{Left: cancelBtnRect.Left + cancelBtnRect.Width() + 30, Top: cancelBtnRect.Top}
		saveBtnRect.SetWidth(100)
		saveBtnRect.SetHeight(35)
		m.saveBtn.SetBoundsRect(saveBtnRect)
		m.saveBtn.SetColor(colors.RGBToColor(46, 204, 113))
		m.saveBtn.SetParent(m)
		m.saveBtn.SetOnClick(m.saveClick)
	}
	//(&hook.TWindowHook{Form: m}).Hook()
}

func (m *TBuildForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	m.closing = true
}

func (m *TBuildForm) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	*closeAction = types.CaFree
}

func (m *TBuildForm) closeClick(sender lcl.IObject) {
	m.Close()
}

func (m *TBuildForm) saveClick(sender lcl.IObject) {

}
