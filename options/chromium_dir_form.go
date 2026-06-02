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
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"os"
	"path/filepath"
)

var (
	chromiumDirFormWidth  = int32(500)
	chromiumDirFormHeight = int32(180)
)

// NewChromiumDirForm 创建 CEF 框架目录设置窗口
func NewChromiumDirForm() *TChromiumDirForm {
	newEngForm := lcl.NewEngForm(nil)
	newForm := &TChromiumDirForm{TEngForm: *newEngForm.(*lcl.TEngForm)}
	newForm.FormCreate(newForm)
	newForm.SetOnCloseQuery(newForm.OnCloseQuery)
	newForm.SetOnClose(newForm.OnClose)
	return newForm
}

type TChromiumDirForm struct {
	lcl.TEngForm
	closing   bool
	selectDir lcl.ISelectDirectoryDialog

	// 目录设置
	dirText lcl.ILabel
	dirEdit lcl.ILabeledEdit
	dirBtn  *wg.TButton

	// 操作按钮
	defaultBtn *wg.TButton
	cancelBtn  *wg.TButton
	confirmBtn *wg.TButton
}

func (m *TChromiumDirForm) FormCreate(sender lcl.IObject) {
	logs.Debug("TChromiumDirForm FormCreate")
	m.SetName("ChromiumDirForm")
	m.SetCaption("CEF 框架目录设置")
	m.SetWidth(chromiumDirFormWidth)
	m.SetHeight(chromiumDirFormHeight)
	m.SetVisible(false)
	m.SetDoubleBuffered(true)
	m.SetBorderIcons(types.NewSet(types.BiSystemMenu))
	m.SetColor(colors.ClWhite)
	SetWindowCenterByMainWindow(m)

	m.selectDir = lcl.NewSelectDirectoryDialog(m)
	m.selectDir.SetName("ChromiumDirFormSelectDir")
	m.selectDir.SetTitle("选择 CEF 框架安装目录")

	gTop := int32(0)
	nextTop := func(top int32) int32 {
		gTop += top
		return gTop
	}

	// 说明文字
	{
		m.dirText = lcl.NewLabel(m)
		m.dirText.SetLeft(15)
		m.dirText.SetTop(nextTop(15))
		m.dirText.SetCaption("设置 CEF 框架目录，请选择安装目录或使用默认目录。")
		m.dirText.SetParent(m)
	}

	// 目录输入
	{
		m.dirEdit = lcl.NewLabeledEdit(m)
		m.dirEdit.SetName("ChromiumDirFormDirEdit")
		m.dirEdit.SetLeft(80)
		m.dirEdit.SetTop(nextTop(30))
		m.dirEdit.SetWidth(350)
		m.dirEdit.SetDoubleBuffered(true)
		m.dirEdit.SetTextHint(config.Config.Chromium.DefaultDir())
		m.dirEdit.SetLabelPosition(types.LpLeft)
		m.dirEdit.SetText("")
		dirLabelText := m.dirEdit.EditLabel()
		dirLabelText.SetCaption("安装目录:")
		m.dirEdit.SetParent(m)

		m.dirBtn = wg.NewButton(m)
		m.dirBtn.SetIconFormBytes(resources.Images("menu/menu_project_open.png"))
		m.dirBtn.SetRadius(3)
		cusRect := types.TRect{Left: m.dirEdit.Left() + m.dirEdit.Width() + 5, Top: m.dirEdit.Top()}
		cusRect.SetWidth(35)
		if tool.IsLinux {
			cusRect.SetHeight(35)
		} else {
			cusRect.SetHeight(25)
		}
		m.dirBtn.SetBoundsRect(cusRect)
		m.dirBtn.SetColor(grayBtnColor)
		m.dirBtn.SetBorderColor(wg.BbdNone, grayBtnColor)
		m.dirBtn.SetCursor(types.CrHandPoint)
		m.dirBtn.SetParent(m)
		m.dirBtn.SetOnClick(m.dirBtnClick)
	}

	// 操作按钮
	{
		defaultBtnRect := types.TRect{Left: 15, Top: nextTop(35)}
		defaultBtnRect.SetWidth(100)
		defaultBtnRect.SetHeight(25)
		m.defaultBtn = wg.NewButton(m)
		m.defaultBtn.SetName("ChromiumDirFormDefaultBtn")
		m.defaultBtn.SetText("使用默认目录")
		m.defaultBtn.Font().SetSize(8)
		m.defaultBtn.SetRadius(3)
		m.defaultBtn.SetBoundsRect(defaultBtnRect)
		m.defaultBtn.SetColor(grayBtnColor)
		m.defaultBtn.SetCursor(types.CrHandPoint)
		m.defaultBtn.SetParent(m)
		m.defaultBtn.SetOnClick(m.defaultBtnClick)

		cancelBtnRect := types.TRect{Left: 325, Top: defaultBtnRect.Top}
		cancelBtnRect.SetWidth(60)
		cancelBtnRect.SetHeight(25)
		m.cancelBtn = wg.NewButton(m)
		m.cancelBtn.SetName("ChromiumDirFormCancelBtn")
		m.cancelBtn.SetText("取 消")
		m.cancelBtn.Font().SetSize(8)
		m.cancelBtn.SetRadius(3)
		m.cancelBtn.SetBoundsRect(cancelBtnRect)
		m.cancelBtn.SetColor(grayBtnColor)
		m.cancelBtn.SetCursor(types.CrHandPoint)
		m.cancelBtn.SetParent(m)
		m.cancelBtn.SetOnClick(m.cancelBtnClick)

		confirmBtnRect := types.TRect{Left: cancelBtnRect.Left + cancelBtnRect.Width() + 20, Top: cancelBtnRect.Top}
		confirmBtnRect.SetWidth(60)
		confirmBtnRect.SetHeight(25)
		m.confirmBtn = wg.NewButton(m)
		m.confirmBtn.SetName("ChromiumDirFormConfirmBtn")
		m.confirmBtn.SetText("确 定")
		m.confirmBtn.Font().SetSize(8)
		m.confirmBtn.Font().SetColor(colors.ClWhite)
		m.confirmBtn.SetRadius(3)
		m.confirmBtn.SetBoundsRect(confirmBtnRect)
		m.confirmBtn.SetColor(blueBtnColor)
		m.confirmBtn.SetCursor(types.CrHandPoint)
		m.confirmBtn.SetParent(m)
		m.confirmBtn.SetOnClick(m.confirmBtnClick)
	}
}

func (m *TChromiumDirForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	m.closing = true
}

func (m *TChromiumDirForm) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	*closeAction = types.CaFree
}

// cancelBtnClick 取消
func (m *TChromiumDirForm) cancelBtnClick(sender lcl.IObject) {
	m.Close()
}

// dirBtnClick 浏览目录
func (m *TChromiumDirForm) dirBtnClick(sender lcl.IObject) {
	if m.selectDir.Execute() {
		m.dirEdit.SetText(m.selectDir.FileName())
	}
}

// defaultBtnClick 使用默认目录
func (m *TChromiumDirForm) defaultBtnClick(sender lcl.IObject) {
	defaultDir := config.Config.Chromium.DefaultDir()
	_ = os.MkdirAll(defaultDir, os.ModePerm)
	m.dirEdit.SetText(defaultDir)
	event.ConsoleWriteInfo("CEF framework directory set to default:", defaultDir)
}

// confirmBtnClick 确认自定义目录
func (m *TChromiumDirForm) confirmBtnClick(sender lcl.IObject) {
	dir := m.dirEdit.Text()
	if dir == "" {
		m.defaultBtnClick(sender)
		return
	}
	// 转为绝对路径
	absDir, err := filepath.Abs(dir)
	if err != nil {
		event.ConsoleWriteError("Invalid directory path:", err.Error())
		return
	}
	// 目录不存在则创建
	if !tool.IsExist(absDir) {
		if err = os.MkdirAll(absDir, os.ModePerm); err != nil {
			event.ConsoleWriteError("Failed to create directory:", err.Error())
			return
		}
	}
	config.Config.Chromium.Dir = absDir
	config.UpdateConfig()
	event.ConsoleWriteInfo("CEF framework directory set to:", absDir)
	m.Close()
}
