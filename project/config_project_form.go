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
	"github.com/energye/designer/project/bean"
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

	appIconBtn       *wg.TButton
	appRemoveIconBtn *wg.TButton
	appIconData      bean.TAppIcon // 应用配置窗口打开时默认从项目配置读取并设置
	appIdText        lcl.ILabel
	appIdEdit        lcl.IEdit
	appDescText      lcl.ILabel
	appDescEdit      lcl.IEdit
	appVersionText   lcl.ILabel
	appVersionEdit   lcl.IEdit
	appCopyrightText lcl.ILabel
	appCopyrightEdit lcl.IEdit

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
	m.box.SetColor(colors.ClWhite)
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
		// 初始时设置图标
		m.appIconData = gProject.AppOption.Icon
		go func() {
			// 预览
			if gProject.AppOption.Icon.Data != nil {
				// 预览图标过大, 需要绽放
				tempIconData := gProject.AppOption.Icon.Data
				if gProject.AppOption.Icon.W > 128 || gProject.AppOption.Icon.H > 128 {
					tempIconData = tool.Scale(tempIconData, 128, 128)
				}
				lcl.RunOnMainThreadAsync(func(id uint32) {
					m.appIconBtnLoadData(tempIconData, "")
				})
			}
		}()
	})
}

// initComponents 初始化项目配置表单中的所有组件。
// 此函数负责创建并设置表单中各个 UI 元素的位置、样式和行为，
// 包括应用程序标题输入框、图标选择区域、平台配置选项卡以及保存/取消按钮等。
func (m *TConfigProjectForm) initComponents() {
	m.selectDir = lcl.NewSelectDirectoryDialog(m)

	left := int32(10)
	textLeft := int32(45)

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
		m.appTitleText.SetCaption("标题")
		m.appTitleText.SetFont(m.font)
		m.appTitleText.SetParent(m.box)

		m.appTitleEdit = lcl.NewEdit(m)
		m.appTitleEdit.SetBounds(textLeft, baseTop-5, 340, 30)
		m.appTitleEdit.SetFont(m.font)
		m.appTitleEdit.SetTextHint("my energy app")
		m.appTitleEdit.SetText(gProject.AppOption.Title)
		m.appTitleEdit.SetParent(m.box)
	}

	{
		m.appIdText = lcl.NewLabel(m)
		m.appIdText.SetLeft(left)
		m.appIdText.SetTop(baseTop + 40)
		m.appIdText.SetCaption("标识")
		m.appIdText.SetParent(m.box)
		m.appIdEdit = lcl.NewEdit(m)
		m.appIdEdit.SetBounds(textLeft, baseTop+35, 340, 30)
		m.appIdEdit.SetFont(m.font)
		m.appIdEdit.SetTextHint("company.product.app")
		m.appIdEdit.SetText(gProject.AppOption.Id)
		m.appIdEdit.SetParent(m.box)

		m.appDescText = lcl.NewLabel(m)
		m.appDescText.SetLeft(left)
		m.appDescText.SetTop(baseTop + 80)
		m.appDescText.SetCaption("描述")
		m.appDescText.SetParent(m.box)
		m.appDescEdit = lcl.NewEdit(m)
		m.appDescEdit.SetBounds(textLeft, baseTop+75, 340, 30)
		m.appDescEdit.SetFont(m.font)
		m.appDescEdit.SetTextHint("your application description.")
		m.appDescEdit.SetText(gProject.AppOption.Desc)
		m.appDescEdit.SetParent(m.box)

		m.appVersionText = lcl.NewLabel(m)
		m.appVersionText.SetLeft(left)
		m.appVersionText.SetTop(baseTop + 120)
		m.appVersionText.SetCaption("版本")
		m.appVersionText.SetParent(m.box)
		m.appVersionEdit = lcl.NewEdit(m)
		m.appVersionEdit.SetBounds(textLeft, baseTop+115, 100, 30)
		m.appVersionEdit.SetFont(m.font)
		m.appVersionEdit.SetTextHint("1.2.3.4")
		m.appVersionEdit.SetText(gProject.AppOption.Version)
		m.appVersionEdit.SetParent(m.box)

		m.appCopyrightText = lcl.NewLabel(m)
		m.appCopyrightText.SetLeft(m.appVersionEdit.Left() + m.appVersionEdit.Width() + left)
		m.appCopyrightText.SetTop(baseTop + 120)
		m.appCopyrightText.SetCaption("版权")
		m.appCopyrightText.SetParent(m.box)
		m.appCopyrightEdit = lcl.NewEdit(m)
		m.appCopyrightEdit.SetBounds(m.appCopyrightText.Left()+35, baseTop+115, 195, 30)
		m.appCopyrightEdit.SetFont(m.font)
		m.appCopyrightEdit.SetTextHint("Copyright (C)")
		m.appCopyrightEdit.SetText(gProject.AppOption.Copyright)
		m.appCopyrightEdit.SetParent(m.box)

		m.appIconBtn = wg.NewButton(m)
		m.appIconBtn.SetIconFormBytes(resources.Images("button/upload_64x64.png"))
		//m.appIconBtn.SetIconCloseFormBytes(resources.Images("button/remove_16x16.png"))
		//m.appIconBtn.SetIconCloseHighlightFormBytes(resources.Images("button/remove_16x16_highlight.png"))
		m.appIconBtn.SetRadius(3)
		appIconRect := types.TRect{Left: m.Width() - 158, Top: baseTop - 5}
		appIconRect.SetWidth(145)
		appIconRect.SetHeight(145)
		m.appIconBtn.TextOffSetY = 50
		m.appIconBtn.SetBoundsRect(appIconRect)
		m.appIconBtn.SetCaption("点击加载应用图标")
		m.appIconBtn.SetHint("点击加载应用图标")
		m.appIconBtn.SetShowHint(true)
		m.appIconBtn.SetFont(m.font)
		m.appIconBtn.SetBorderColor(wg.BbdNone, colors.RGBToColor(91, 155, 213))
		m.appIconBtn.SetBorderWidth(wg.BbdNone, 1)
		m.appIconBtn.SetColor(0xF3F4F6)
		m.appIconBtn.SetParent(m.box)
		m.appIconBtn.SetOnMouseUp(m.appIconBtnClick)

		m.appRemoveIconBtn = wg.NewButton(m)
		m.appRemoveIconBtn.SetRadius(3)
		m.appRemoveIconBtn.SetIconFormBytes(resources.Images("button/remove_16x16_highlight.png"))
		appRemoveIconRect := types.TRect{Left: appIconRect.Left + appIconRect.Width() - 24, Top: appIconRect.Top}
		appRemoveIconRect.SetWidth(24)
		appRemoveIconRect.SetHeight(24)
		m.appRemoveIconBtn.SetBoundsRect(appRemoveIconRect)
		m.appRemoveIconBtn.SetColor(colors.RGBToColor(255, 127, 127))
		m.appRemoveIconBtn.SetAlpha(200)
		m.appRemoveIconBtn.SetVisible(false)
		m.appRemoveIconBtn.SetParent(m.box)
		m.appRemoveIconBtn.SetOnMouseUp(m.appRemoveIconBtnClick)
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

		tabColor := colors.ClWhite //colors.TColor(0xF3F4F6)
		btnColor := colors.RGBToColor(173, 216, 230)

		m.platformTab = wg.NewTab(m)
		m.platformTab.Margin = 10
		tabBR := types.TRect{Left: 0, Top: m.platformTitle.Top() + 25}
		tabBR.SetWidth(m.Width())
		tabBR.SetHeight(m.Height() - (tabBR.Top - 10))
		m.platformTab.SetBoundsRect(tabBR)
		m.platformTab.SetColor(colors.ClWhite)
		m.platformTab.EnableScrollButton(false)
		m.platformTab.SetParent(m.box)
		m.platformTab.SetOnChange(func(sender lcl.IObject) {
			for _, page := range m.platformTab.Pages() {
				if page.Active() {
					page.Button().SetBorderDirections(0)
				} else {
					page.Button().SetBorderDirections(types.NewSet(wg.BbdBottom, wg.BbdLeft, wg.BbdTop, wg.BbdRight))
				}
			}
		})
		// 设置标签按钮样式
		setTabPageStyle := func(page *wg.TPage) {
			page.SetTop(40)
			page.SetHeight(m.platformTab.Height() - 40)
			page.SetColor(m.platformTab.Color()) // 设置背景色
			page.Button().SetHeight(35)
			page.Button().SetLeft(10)
			page.Button().RoundedCorner = types.NewSet(wg.RcLeftTop, wg.RcRightTop, wg.RcLeftBottom, wg.RcRightBottom)
			page.Button().Font().SetColor(colors.ClBlack)
			page.Button().SetBorderColor(wg.BbdNone, wg.DarkenColor(tabColor, 0.1))
			page.Button().SetRadius(35)
			page.Button().SetColor(tabColor)
			page.Button().SetDownColor(wg.DarkenColor(btnColor, 0.15), wg.DarkenColor(btnColor, 0.15))
			page.Button().SetEnterColor(wg.DarkenColor(btnColor, 0.1), wg.DarkenColor(btnColor, 0.1))
			page.SetDefaultColor(tabColor)
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
		m.cancelBtn.SetText("关　闭")
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

func (m *TConfigProjectForm) validateInputs() bool {
	return false
}

func (m *TConfigProjectForm) saveClick(sender lcl.IObject) {
	if m.validateInputs() {
		logs.Error("保存-验证输入失败")
		return
	}
	m.cancelBtn.SetDisable(true)
	m.saveBtn.SetDisable(true)
	// 保存项目配置
	m.saveProjectConfig()
	go func() {
		defer func() {
			if !m.closing {
				// 恢复按钮
				m.cancelBtn.SetDisable(false)
				m.saveBtn.SetDisable(false)
			}
		}()
		// 更新 windows 配置并生成程序信息
		saveOrUpdateWindowsManifest()
		// 更新图标
		updateWindowICON()
	}()
}

// appIconBtnLoadData 加载应用图标按钮的数据和文本
// 该函数用于设置应用图标按钮显示的图标和文本内容，并根据是否有图标数据来控制移除图标按钮的显示状态
//   - data: 图标图像数据字节切片，如果为nil则使用默认的上传图标
//   - text: 要设置给图标按钮的文本标题
func (m *TConfigProjectForm) appIconBtnLoadData(data []byte, text string) {
	if data == nil {
		m.appRemoveIconBtn.Hide()
		data = resources.Images("button/upload_64x64.png")
		m.appIconData = bean.TAppIcon{}
	} else {
		m.appRemoveIconBtn.Show()
	}
	m.appIconBtn.SetIconFormBytes(data)
	m.appIconBtn.SetCaption(text)
}

// appRemoveIconBtnClick 处理应用图标移除按钮的点击事件
// 当用户点击移除图标按钮时，该函数会被触发，用于清除当前的应用图标并重置为默认状态
func (m *TConfigProjectForm) appRemoveIconBtnClick(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	m.appIconBtnLoadData(nil, "点击加载应用图标")
}

// appIconBtnClick 是应用图标按钮的点击事件处理函数。
// 当用户点击图标设置按钮时，打开图形属性编辑器用于选择或上传新的应用图标。
func (m *TConfigProjectForm) appIconBtnClick(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
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
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.appIconBtnLoadData(previewData, "")
			})
			// 缩放到
			// windows: 256x256
			// macos: 1024x1024
			// linux: 256x256
			saveData := data
			w, h := imageInfo.Rect.Width(), imageInfo.Rect.Height()
			if w > 1024 || h > 1024 {
				saveData = tool.Scale(data, 1024, 1024)
				w, h = 1204, 1204
			}
			m.appIconData = bean.TAppIcon{Data: saveData, W: w, H: h}
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

func (m *TConfigProjectForm) AppTitle() string {
	title := m.appTitleEdit.Text()
	if title == "" {
		title = gProject.AppOption.Title
	}
	return title
}

func (m *TConfigProjectForm) AppCopyright() string {
	copyright := m.appCopyrightEdit.Text()
	if copyright == "" {
		copyright = gProject.AppOption.Copyright
	}
	return copyright
}

func (m *TConfigProjectForm) AppId() string {
	id := m.appIdEdit.Text()
	if id == "" {
		id = gProject.AppOption.Id
	}
	return id
}

func (m *TConfigProjectForm) AppDesc() string {
	desc := m.appDescEdit.Text()
	if desc == "" {
		desc = gProject.AppOption.Desc
	}
	return desc
}

func (m *TConfigProjectForm) AppVersion() string {
	appVersion := m.appVersionEdit.Text()
	if appVersion == "" {
		appVersion = gProject.AppOption.Version
	}
	return appVersion
}
