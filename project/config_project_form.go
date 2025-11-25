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
	"github.com/energye/designer/pkg/helperform"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/font"
	"github.com/energye/widget/wg"
	"os"
	"sync"
)

var (
	configProjectFormWidth  = int32(555)
	configProjectFormHeight = int32(515)
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
	closing   bool
	one       sync.Once
	box       lcl.IPanel
	selectDir lcl.ISelectDirectoryDialog

	font lcl.IFont

	appConfigTitle lcl.ILabel

	appTitleText lcl.ILabel
	appTitleEdit lcl.IEdit

	appIconBtn     *wg.TButton
	appIconData    []byte
	appIdText      lcl.ILabel
	appIdEdit      lcl.IEdit
	appDescText    lcl.ILabel
	appDescEdit    lcl.IEdit
	appVersionText lcl.ILabel
	appVersionEdit lcl.IEdit

	platformTitle lcl.ILabel

	platformTab            *wg.TTab
	platformTabPageWindows *wg.TPage
	platformTabPageMacOS   *wg.TPage
	platformTabPageLinux   *wg.TPage

	// windows manifest

	compatibilityOSText                  lcl.ILabel
	compatibilityOSBox                   lcl.IComboBox
	dpiText                              lcl.ILabel
	dpiBox                               lcl.IComboBox
	runLevelText                         lcl.ILabel
	runLevelBox                          lcl.IComboBox
	uiAccessCheckBox                     lcl.ICheckBox
	autoElevateBox                       lcl.ICheckBox
	disableThemingBox                    lcl.ICheckBox
	disableWindowFilteringBox            lcl.ICheckBox
	highResolutionScrollingAwareBox      lcl.ICheckBox
	ultraHighResolutionScrollingAwareBox lcl.ICheckBox
	longPathAwareBox                     lcl.ICheckBox
	printerDriverIsolationBox            lcl.ICheckBox
	gDIScalingBox                        lcl.ICheckBox
	segmentHeapBox                       lcl.ICheckBox
	useCommonControlsV6Box               lcl.ICheckBox

	// macos

	// linux

	// 操作按钮
	cancelBtn *wg.TButton
	saveBtn   *wg.TButton
}

func (m *TConfigProjectForm) FormCreate(sender lcl.IObject) {
	logs.Debug("TConfigProjectForm FormCreate")
	m.SetCaption("应用配置")
	m.SetWidth(configProjectFormWidth)
	m.SetHeight(configProjectFormHeight)
	m.SetVisible(false)
	m.SetDoubleBuffered(true)
	m.SetBorderIcons(types.NewSet(types.BiSystemMenu))
	m.WorkAreaCenter()
	m.box = lcl.NewPanel(m)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetAlign(types.AlClient)
	m.box.SetParent(m)
	m.SetOnShow(m.onShow)
	m.initComponents()

	//(&hook.TWindowHook{Form: m}).Hook()
}

func (m *TConfigProjectForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	m.closing = true
}

func (m *TConfigProjectForm) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
}

// 窗口显示事件
func (m *TConfigProjectForm) onShow(sender lcl.IObject) {
	logs.Debug("TConfigProjectForm Show")
	m.one.Do(func() {
		addSize := int32(60)
		br := m.BoundsRect()
		br.SetWidth(configProjectFormWidth)
		br.SetHeight(configProjectFormHeight + addSize)
		m.SetBoundsRect(br) // trigger WM_NCCALCSIZE hook msg
		constr := m.Constraints()
		constr.SetMaxWidth(configProjectFormWidth)
		constr.SetMaxHeight(configProjectFormHeight + addSize)
		constr.SetMinWidth(configProjectFormWidth)
		constr.SetMinHeight(configProjectFormHeight + addSize)
		m.WorkAreaCenter()
	})
}

// initComponents 初始化项目配置表单中的所有组件。
// 此函数负责创建并设置表单中各个 UI 元素的位置、样式和行为，
// 包括应用程序标题输入框、图标选择区域、平台配置选项卡以及保存/取消按钮等。
func (m *TConfigProjectForm) initComponents() {
	m.selectDir = lcl.NewSelectDirectoryDialog(m)

	left := int32(35)
	textLeft := int32(100)

	m.font = lcl.NewFont()
	m.font.SetName("微软雅黑")
	m.font.SetCharSet(font.CHINESEBIG5_CHARSET)

	m.appConfigTitle = lcl.NewLabel(m)
	m.appConfigTitle.SetLeft(10)
	m.appConfigTitle.SetTop(10)
	m.appConfigTitle.SetCaption("应用程序配置")
	m.appConfigTitle.SetFont(m.font)
	m.appConfigTitle.Font().SetSize(10)
	m.appConfigTitle.SetParent(m.box)

	baseTop := int32(40)
	{
		m.appTitleText = lcl.NewLabel(m)
		m.appTitleText.SetLeft(left)
		m.appTitleText.SetTop(baseTop)
		m.appTitleText.SetCaption("应用标题")
		m.appTitleText.SetFont(m.font)
		m.appTitleText.SetParent(m.box)

		m.appTitleEdit = lcl.NewEdit(m)
		m.appTitleEdit.SetBounds(textLeft, baseTop-5, 280, 30)
		m.appTitleEdit.SetFont(m.font)
		m.appTitleEdit.SetTextHint("my energy app")
		m.appTitleEdit.SetText(gProject.AppOption.Title)
		m.appTitleEdit.SetParent(m.box)
	}

	{
		m.appIdText = lcl.NewLabel(m)
		m.appIdText.SetLeft(left)
		m.appIdText.SetTop(baseTop + 40)
		m.appIdText.SetCaption("应用标识")
		m.appIdText.SetParent(m.box)
		m.appIdEdit = lcl.NewEdit(m)
		m.appIdEdit.SetBounds(textLeft, baseTop+35, 280, 30)
		m.appIdEdit.SetFont(m.font)
		m.appIdEdit.SetTextHint("company.product.app")
		m.appIdEdit.SetText(gProject.AppOption.Id)
		m.appIdEdit.SetParent(m.box)

		m.appDescText = lcl.NewLabel(m)
		m.appDescText.SetLeft(left)
		m.appDescText.SetTop(baseTop + 80)
		m.appDescText.SetCaption("应用描述")
		m.appDescText.SetParent(m.box)
		m.appDescEdit = lcl.NewEdit(m)
		m.appDescEdit.SetBounds(textLeft, baseTop+75, 280, 30)
		m.appDescEdit.SetFont(m.font)
		m.appDescEdit.SetTextHint("your application description.")
		m.appDescEdit.SetText(gProject.AppOption.Desc)
		m.appDescEdit.SetParent(m.box)

		m.appVersionText = lcl.NewLabel(m)
		m.appVersionText.SetLeft(left)
		m.appVersionText.SetTop(baseTop + 120)
		m.appVersionText.SetCaption("应用版本")
		m.appVersionText.SetParent(m.box)
		m.appVersionEdit = lcl.NewEdit(m)
		m.appVersionEdit.SetBounds(textLeft, baseTop+115, 280, 30)
		m.appVersionEdit.SetFont(m.font)
		m.appVersionEdit.SetTextHint("1.2.3.4")
		m.appVersionEdit.SetText(gProject.AppOption.Version)
		m.appVersionEdit.SetParent(m.box)

		m.appIconBtn = wg.NewButton(m)
		m.appIconBtn.SetIconFormBytes(resources.Images("button/upload_64x64.png"))
		m.appIconBtn.SetRadius(3)
		appIconRect := types.TRect{Left: m.Width() - 158, Top: baseTop - 5}
		appIconRect.SetWidth(145)
		appIconRect.SetHeight(145)
		m.appIconBtn.TextOffSetY = 50
		m.appIconBtn.SetCaption("点击加载应用图标")
		m.appIconBtn.SetHint("点击加载应用图标")
		m.appIconBtn.SetShowHint(true)
		m.appIconBtn.SetFont(m.font)
		m.appIconBtn.SetBoundsRect(appIconRect)
		m.appIconBtn.SetBorderColor(wg.BbdNone, colors.RGBToColor(91, 155, 213))
		m.appIconBtn.SetBorderWidth(wg.BbdNone, 1)
		m.appIconBtn.SetColor(0xF3F4F6)
		m.appIconBtn.SetParent(m.box)
		m.appIconBtn.SetOnClick(m.appIconBtnClick)
	}

	{
		m.platformTitle = lcl.NewLabel(m)
		m.platformTitle.SetLeft(10)
		m.platformTitle.SetTop(baseTop + 155)
		m.platformTitle.SetCaption("平台配置")
		m.platformTitle.SetFont(m.font)
		m.platformTitle.Font().SetSize(10)
		m.platformTitle.SetParent(m.box)
	}

	{

		tabBtnColor := colors.TColor(0xF3F4F6)
		m.platformTab = wg.NewTab(m)
		m.platformTab.Margin = 10
		tabBR := types.TRect{Left: 0, Top: m.platformTitle.Top() + 25}
		tabBR.SetWidth(m.Width())
		tabBR.SetHeight(m.Height() - (tabBR.Top - 10))
		m.platformTab.SetBoundsRect(tabBR)
		m.platformTab.SetColor(tabBtnColor)
		m.platformTab.EnableScrollButton(false)
		m.platformTab.SetParent(m.box)
		m.platformTab.SetOnChange(func(sender lcl.IObject) {
			for _, page := range m.platformTab.Pages() {
				if page.Active() {
					//page.Button().SetBorderDirections(types.NewSet(wg.BbdBottom))
				} else {
					//page.Button().SetBorderDirections(0)
				}
			}
		})

		btnColor := colors.RGBToColor(173, 216, 230)
		// 设置标签按钮样式
		setTabPageStyle := func(page *wg.TPage) {
			page.SetTop(40)
			page.SetHeight(m.platformTab.Height() - 40)
			page.SetColor(m.platformTab.Color()) // 设置背景色
			page.Button().SetHeight(35)
			page.Button().SetLeft(10)
			page.Button().RoundedCorner = types.NewSet(wg.RcLeftTop, wg.RcRightTop, wg.RcLeftBottom, wg.RcRightBottom)
			page.Button().Font().SetColor(colors.ClBlack)
			page.Button().SetBorderColor(wg.BbdNone, wg.LightenColor(colors.ClGray, 0.5))
			page.Button().SetColor(tabBtnColor)
			page.Button().SetRadius(35)
			page.SetDefaultColor(tabBtnColor)
			page.Button().SetDownColor(wg.DarkenColor(btnColor, 0.15), wg.DarkenColor(btnColor, 0.15))
			page.Button().SetEnterColor(wg.DarkenColor(btnColor, 0.1), wg.DarkenColor(btnColor, 0.1))
			page.SetActiveColor(btnColor)
		}

		m.platformTabPageWindows = m.platformTab.NewPage()
		m.platformTabPageWindows.SetCaption("　Windows　")
		m.platformTabPageWindows.Button().SetIconFavoriteFormBytes(resources.Images("button/windows.png"))
		setTabPageStyle(m.platformTabPageWindows)
		m.initWindowsOptions()

		m.platformTabPageMacOS = m.platformTab.NewPage()
		m.platformTabPageMacOS.SetCaption("　MacOS　")
		m.platformTabPageMacOS.Button().SetIconFavoriteFormBytes(resources.Images("button/macos.png"))
		setTabPageStyle(m.platformTabPageMacOS)
		m.initMacOSOptions()

		m.platformTabPageLinux = m.platformTab.NewPage()
		m.platformTabPageLinux.SetCaption("　Linux　")
		m.platformTabPageLinux.Button().SetIconFavoriteFormBytes(resources.Images("button/linux.png"))
		setTabPageStyle(m.platformTabPageLinux)
		m.initLinuxOptions()

		m.platformTabPageWindows.SetActive(true)

	}

	{
		m.cancelBtn = wg.NewButton(m)
		m.cancelBtn.SetText("取　消")
		m.cancelBtn.SetFont(m.font)
		m.cancelBtn.Font().SetColor(colors.ClWhite)
		m.cancelBtn.SetRadius(3)
		cancelBtnRect := types.TRect{Left: 315, Top: 530}
		cancelBtnRect.SetWidth(100)
		cancelBtnRect.SetHeight(35)
		m.cancelBtn.SetBoundsRect(cancelBtnRect)
		m.cancelBtn.SetColor(colors.RGBToColor(255, 127, 127))
		m.cancelBtn.SetParent(m.box)
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
		m.saveBtn.SetParent(m.box)
		m.saveBtn.SetOnClick(m.saveClick)
	}
}

func (m *TConfigProjectForm) closeClick(sender lcl.IObject) {
	m.Close()
}

func (m *TConfigProjectForm) saveClick(sender lcl.IObject) {
	m.saveWindows()
}

// 应用程序图标
func (m *TConfigProjectForm) appIconBtnClick(sender lcl.IObject) {
	priceForm := helperform.NewGraphicPropertyEditor(func(imageInfo helperform.ImageInfo) {
		if !imageInfo.OK {
			return
		}
		go func() {
			var (
				data []byte
				err  error
			)
			if imageInfo.FilePath == "" {
				data = imageInfo.Data
			} else {
				data, err = os.ReadFile(imageInfo.FilePath)
				if err != nil {
					logs.Error("图标加载 PNG ReadFile:", err.Error())
					return
				}
			}
			//imageFormat, err := tool.DetectImageFormatByte(data)
			//if err != nil {
			//	logs.Error("图标加载 PNG DetectImageFormatByte:", err.Error())
			//	return
			//}
			//if !tool.Equal(imageFormat, "png") {
			//	// TODO 非 png 需要转换为 png
			//}
			if data == nil {
				return
			}
			previewData := data
			if imageInfo.Rect.Width() > 128 || imageInfo.Rect.Height() > 128 {
				previewData = tool.Scale(data, 128, 128)
			}
			// 预览
			//mem := lcl.NewMemoryStream()
			//lcl.StreamHelper.WriteBuffer(mem, previewData)
			//mem.SetPosition(0)
			lcl.RunOnMainThreadAsync(func(id uint32) {
				//m.appIconPreview.Picture().LoadFromStream(mem)
				m.appIconBtn.SetIconFormBytes(previewData)
				m.appIconBtn.SetCaption("")
				//mem.Free()
			})
			// 缩放到
			// windows: 256x256
			// macos: 1024x1024
			// linux: 256x256
			saveData := data
			if imageInfo.Rect.Width() > 256 || imageInfo.Rect.Height() > 256 {
				saveData = tool.Scale(data, 256, 256)
			}
			m.appIconData = saveData
		}()

	})
	priceForm.SetWidth(450)
	priceForm.SetHeight(325)
	priceForm.WorkAreaCenter()
	priceForm.ShowModal()
	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.appIconBtn.Leave(sender)
	})
}
