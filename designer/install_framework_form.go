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

package designer

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/frameworks"
	"github.com/energye/designer/resources/frameworks/lib"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/widget/wg"
	"os"
	"path/filepath"
)

type TInstallFrameworkForm struct {
	*lcl.TEngForm
	installMsgLab lcl.ILabel
	selectDir     lcl.ISelectDirectoryDialog
	selectDirEdit lcl.ILabeledEdit
	cancelBtn     *wg.TButton
	okBtn         *wg.TButton
}

func NewInstallFrameworkForm() *TInstallFrameworkForm {
	newEngForm := lcl.NewEngForm(nil)
	newForm := &TInstallFrameworkForm{TEngForm: newEngForm.(*lcl.TEngForm)}
	newForm.FormCreate(newEngForm)
	newForm.SetOnCloseQuery(newForm.OnCloseQuery)
	newForm.SetOnClose(newForm.OnClose)
	return newForm
}

func (m *TInstallFrameworkForm) FormCreate(sender lcl.IObject) {
	//m.SetCaption("Framework Install Wizard")
	m.SetCaption("框架安装向导")
	m.SetWidth(460)
	m.SetHeight(269)
	constr := m.Constraints()
	constr.SetMaxWidth(m.Width())
	constr.SetMaxHeight(m.Height())
	constr.SetMinWidth(m.Width())
	constr.SetMinHeight(m.Height())
	m.SetBorderStyleToFormBorderStyle(types.BsSingle)
	m.SetBorderIcons(types.NewSet())
	m.WorkAreaCenter()

	m.selectDir = lcl.NewSelectDirectoryDialog(m)

	box := lcl.NewPanel(m)
	box.SetAlign(types.AlClient)
	box.SetParent(m)

	titleLab := lcl.NewLabel(m)
	titleLab.SetCaption("请设置框架安装目录")
	titleLabFont := titleLab.Font()
	titleLabFont.SetSize(14)
	titleLabFont.SetStyle(types.NewSet(types.FsBold))
	titleLab.SetLeft(150)
	titleLab.SetTop(20)
	titleLab.SetParent(box)

	// 软件运行所需组件将安装到此目录
	//仅需设置一次，后续自动使用
	titleCaption := "所需依赖组件库将安装到此目录\n仅需设置一次，后续自动使用"
	titleSubLab := lcl.NewLabel(m)
	titleSubLab.SetCaption(titleCaption)
	titleSubLabFont := titleSubLab.Font()
	titleSubLabFont.SetSize(12)
	titleSubLabFont.SetStyle(types.NewSet(types.FsBold))
	titleSubLabFont.SetColor(colors.ClGray)
	titleSubLab.SetLeft(125)
	titleSubLab.SetTop(60)
	titleSubLab.SetParent(box)

	m.selectDirEdit = lcl.NewLabeledEdit(m)
	m.selectDirEdit.SetLeft(10)
	m.selectDirEdit.SetTop(135)
	m.selectDirEdit.SetWidth(360)
	m.selectDirEdit.EditLabel().SetCaption("框架安装路径：")
	selectDirEditLabFont := m.selectDirEdit.EditLabel().Font()
	selectDirEditLabFont.SetSize(12)
	selectDirEditLabFont.SetStyle(types.NewSet(types.FsBold))
	selectDirEditLabFont.SetColor(0x696969)
	m.selectDirEdit.SetShowHint(true)
	m.selectDirEdit.SetTextHint("请选择有效的空文件夹或新建目录")
	m.selectDirEdit.SetCaption(config.Config.FrameworkDir)
	m.selectDirEdit.SetParent(box)

	selectDirBtn := lcl.NewButton(m)
	selectDirBtn.SetCaption("浏 览...")
	selectDirBtn.SetTop(m.selectDirEdit.Top())
	selectDirBtn.SetLeft(m.selectDirEdit.Width() + 15)
	selectDirBtn.SetParent(box)
	selectDirBtn.SetOnClick(m.selectDirClick)

	m.installMsgLab = lcl.NewLabel(m)
	m.installMsgLab.SetTop(m.selectDirEdit.Top() + 35)
	m.installMsgLab.SetLeft(10)
	m.installMsgLab.SetCaption("---")
	m.installMsgLab.SetParent(box)

	cancelRect := types.TRect{Left: 280, Top: 230}
	cancelRect.SetWidth(80)
	cancelRect.SetHeight(30)
	m.cancelBtn = wg.NewButton(m)
	m.cancelBtn.SetText("关　闭")
	m.cancelBtn.SetBoundsRect(cancelRect)
	m.cancelBtn.SetCursor(types.CrHandPoint)
	m.cancelBtn.Font().SetColor(colors.ClWhite)
	m.cancelBtn.Font().SetStyle(types.NewSet(types.FsBold))
	m.cancelBtn.SetColor(colors.RGBToColor(255, 127, 127))
	m.cancelBtn.SetParent(box)
	m.cancelBtn.SetOnClick(m.cancel)

	okBtnRect := types.TRect{Left: cancelRect.Left + cancelRect.Width() + 10, Top: cancelRect.Top}
	okBtnRect.SetWidth(80)
	okBtnRect.SetHeight(30)
	m.okBtn = wg.NewButton(m)
	m.okBtn.SetText("确定并安装")
	m.okBtn.SetBoundsRect(okBtnRect)
	m.okBtn.SetCursor(types.CrHandPoint)
	m.okBtn.Font().SetColor(colors.ClWhite)
	m.okBtn.Font().SetStyle(types.NewSet(types.FsBold))
	m.okBtn.SetColor(colors.RGBToColor(46, 204, 113))
	m.okBtn.SetParent(box)
	m.okBtn.SetOnClick(m.ok)
}

func (m *TInstallFrameworkForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
}

func (m *TInstallFrameworkForm) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	*closeAction = types.CaFree
}

func (m *TInstallFrameworkForm) cancel(sender lcl.IObject) {
	if m.cancelBtn.Disable() {
		return
	}
	m.Close()
}

func (m *TInstallFrameworkForm) ok(sender lcl.IObject) {
	if m.okBtn.Disable() {
		return
	}
	if m.checkInstallDir() {
		installDir := m.selectDirEdit.Text()
		m.cancelBtn.SetDisable(true)
		m.okBtn.SetDisable(true)
		m.setInstallMsgFont(colors.ClBlack)
		go func() {
			config.UpdateFrameworkDir(installDir)
			config.UpdateConfig()
			// 释放 lib runtime 库文件
			runtimeDir := filepath.Join(installDir, "runtime")
			_ = os.MkdirAll(runtimeDir, os.ModePerm)
			lib.ExtractLibrary(runtimeDir)
			// 提取所有启用的框架
			frameworks.ExtractFrameworks(func(message string) {
				lcl.RunOnMainThreadAsync(func(id uint32) {
					m.installMsgLab.SetCaption(message)
				})
			}, func() {
				// 提取完成
				lcl.RunOnMainThreadAsync(func(id uint32) {
					m.setInstallMsgFont(colors.ClGreen)
					m.installMsgLab.SetCaption("框架安装完成")
				})
				SetEnableFuncComponent(true)
				m.cancelBtn.SetDisable(false)
				m.cancelBtn.SetEnabled(true)
				m.okBtn.SetDisable(false)
				m.okBtn.SetEnabled(true)
			})
		}()
	}
}

// 检查安装目录是否有效
func (m *TInstallFrameworkForm) checkInstallDir() bool {
	installDir := m.selectDirEdit.Text()
	if !tool.IsExist(installDir) {
		m.setInstallMsgFont(colors.ClRed)
		m.installMsgLab.SetCaption("无效目录: " + installDir)
		return false
	}
	m.setInstallMsgFont(colors.ClBlack)
	m.installMsgLab.SetCaption("---")
	return true
}

// 选择目录事件
func (m *TInstallFrameworkForm) selectDirClick(sender lcl.IObject) {
	m.selectDir.SetTitle("选择框架安装目录")
	if m.selectDir.Execute() {
		installDir := m.selectDir.FileName()
		m.selectDirEdit.SetText(installDir)
		m.checkInstallDir()
	}
}

// 设置安装提示消息 label 的 font
func (m *TInstallFrameworkForm) setInstallMsgFont(color colors.TColor) {
	m.installMsgLab.Font().SetColor(color)
}
