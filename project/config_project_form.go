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

	appIcon        lcl.ILabel
	appIconBox     lcl.IScrollBox
	appIconPreview lcl.IImage
	appIconBtn     *wg.TButton
	appIconData    []byte

	platformTitle lcl.ILabel

	platformTab            *wg.TTab
	platformTabPageWindows *wg.TPage
	platformTabPageMacOS   *wg.TPage
	platformTabPageLinux   *wg.TPage

	// windows manifest
	appNameText                          lcl.ILabel
	appNameEdit                          lcl.IEdit
	appDescText                          lcl.ILabel
	appDescEdit                          lcl.IEdit
	appVersionText                       lcl.ILabel
	appVersionEdit                       lcl.IEdit
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
		m.appTitleEdit.SetLeft(textLeft)
		m.appTitleEdit.SetTop(baseTop - 5)
		m.appTitleEdit.SetWidth(440)
		m.appTitleEdit.SetFont(m.font)
		m.appTitleEdit.SetParent(m.box)
		m.appTitleEdit.SetTextHint("my energy app")
	}

	{
		m.appIcon = lcl.NewLabel(m)
		m.appIcon.SetLeft(left)
		m.appIcon.SetTop(baseTop + 40)
		m.appIcon.SetCaption("应用图标")
		m.appIcon.SetFont(m.font)
		m.appIcon.SetParent(m.box)

		m.appIconBox = lcl.NewScrollBox(m)
		m.appIconBox.SetLeft(textLeft)
		m.appIconBox.SetTop(baseTop + 40)
		m.appIconBox.SetWidth(128)
		m.appIconBox.SetHeight(128)
		m.appIconBox.SetAutoScroll(false)
		m.appIconBox.SetBorderStyleToBorderStyle(types.BsSingle)
		m.appIconBox.HorzScrollBar().SetTracking(true)
		m.appIconBox.HorzScrollBar().SetVisible(true)
		m.appIconBox.VertScrollBar().SetTracking(true)
		m.appIconBox.VertScrollBar().SetVisible(true)
		m.appIconBox.SetParent(m.box)

		m.appIconPreview = lcl.NewImage(m)
		m.appIconPreview.SetAlign(types.AlClient)
		m.appIconPreview.SetAutoSize(true)
		m.appIconPreview.SetCenter(true)
		m.appIconPreview.SetParent(m.appIconBox)
		//m.appIconPreview.SetOnPaintBackground(m.appIconPreviewPaintBackground)

		m.appIconBtn = wg.NewButton(m)
		m.appIconBtn.SetIconFormBytes(resources.Images("button/image.png"))
		m.appIconBtn.SetRadius(3)
		appIconRect := types.TRect{Left: m.appIconBox.Left() + m.appIconBox.Width() + 15, Top: baseTop + 40}
		appIconRect.SetWidth(48)
		appIconRect.SetHeight(48)
		m.appIconBtn.SetBoundsRect(appIconRect)
		m.appIconBtn.SetBorderColor(wg.BbdNone, colors.RGBToColor(91, 155, 213))
		m.appIconBtn.SetBorderWidth(wg.BbdNone, 1)
		m.appIconBtn.SetColor(colors.RGBToColor(135, 206, 235))
		m.appIconBtn.SetParent(m.box)
		m.appIconBtn.SetOnClick(m.appIconBtnClick)
	}

	{
		m.platformTitle = lcl.NewLabel(m)
		m.platformTitle.SetLeft(10)
		m.platformTitle.SetTop(baseTop + 175)
		m.platformTitle.SetCaption("平台配置")
		m.platformTitle.SetFont(m.font)
		m.platformTitle.Font().SetSize(10)
		m.platformTitle.SetParent(m.box)
	}

	{
		setPlatformTitle := func(title string) {
			m.platformTitle.SetCaption(title)
		}
		tabBtnColor := colors.TColor(0xF3F4F6)
		m.platformTab = wg.NewTab(m)
		tabBR := types.TRect{Left: 0, Top: m.platformTitle.Top() + 25}
		tabBR.SetWidth(m.Width())
		tabBR.SetHeight(m.Height() - (tabBR.Top - 10))
		m.platformTab.SetBoundsRect(tabBR)
		m.platformTab.SetColor(tabBtnColor)
		m.platformTab.EnableScrollButton(false)
		m.platformTab.SetParent(m.box)
		m.platformTab.SetOnChange(func(sender lcl.IObject) {
			for i, page := range m.platformTab.Pages() {
				if page.Active() {
					page.Button().SetBorderDirections(types.NewSet(wg.BbdBottom))
					switch i {
					case 0:
						setPlatformTitle("平台配置 - 适用于 Windows 系统 (Manifest 资源)")
					case 1:
						setPlatformTitle("平台配置 - 适用于 MacOS 系统")
					case 2:
						setPlatformTitle("平台配置 - 适用于 Linux 系统")
					}
				} else {
					page.Button().SetBorderDirections(0)
				}
			}
		})

		// 设置标签按钮样式
		setTabPageStyle := func(page *wg.TPage) {
			page.SetColor(m.platformTab.Color()) // 设置背景色
			page.Button().Font().SetColor(colors.ClBlack)
			page.Button().SetBorderDirections(types.NewSet(wg.BbdBottom))
			page.Button().SetBorderWidth(wg.BbdBottom, 2)
			page.Button().SetBorderColor(wg.BbdBottom, colors.TColor(0xdb9612))
			page.Button().SetColor(tabBtnColor)
			page.SetDefaultColor(tabBtnColor)
			page.SetActiveColor(wg.DarkenColor(tabBtnColor, 0.15))
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
		m.cancelBtn.SetText("关 闭")
		m.cancelBtn.SetFont(m.font)
		m.cancelBtn.Font().SetColor(colors.ClWhite)
		m.cancelBtn.Font().SetStyle(types.NewSet(types.FsBold))
		m.cancelBtn.SetRadius(3)
		cancelBtnRect := types.TRect{Left: 315, Top: 530}
		cancelBtnRect.SetWidth(100)
		cancelBtnRect.SetHeight(35)
		m.cancelBtn.SetBoundsRect(cancelBtnRect)
		m.cancelBtn.SetColor(colors.RGBToColor(255, 127, 127))
		m.cancelBtn.SetParent(m.box)
		m.cancelBtn.SetOnClick(m.closeClick)

		m.saveBtn = wg.NewButton(m)
		m.saveBtn.SetText("保 存")
		m.saveBtn.SetFont(m.font)
		m.saveBtn.Font().SetStyle(types.NewSet(types.FsBold))
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

			previewData := data
			if imageInfo.Rect.Width() > 128 || imageInfo.Rect.Height() > 128 {
				previewData = tool.Scale(data, 128, 128)
			}
			// 预览
			mem := lcl.NewMemoryStream()
			lcl.StreamHelper.WriteBuffer(mem, previewData)
			mem.SetPosition(0)
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.appIconPreview.Picture().LoadFromStream(mem)
				mem.Free()
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
}
