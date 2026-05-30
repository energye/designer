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
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/helperform"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/designer/resources/metadata"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/font"
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
	newEngForm := lcl.NewEngForm(nil)
	newForm := &TConfigProjectForm{TEngForm: *newEngForm.(*lcl.TEngForm)}
	//lcl.Application.NewForm(newForm) // 不使用原因：go debug 模式有问题
	newForm.FormCreate(newEngForm)
	newForm.SetOnCloseQuery(newForm.OnCloseQuery)
	newForm.SetOnClose(newForm.OnClose)
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

	appIconBtn       *wg.TButton
	appRemoveIconBtn *wg.TButton
	appIconData      bean.TAppIcon // 应用配置窗口打开时默认从项目配置读取并设置
	appTitleEdit     lcl.ILabeledEdit
	appIdEdit        lcl.ILabeledEdit
	appDescEdit      lcl.ILabeledEdit
	appVersionEdit   lcl.ILabeledEdit
	appCopyrightEdit lcl.ILabeledEdit

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
	CFBundleNameEdit           lcl.ILabeledEdit
	CFBundleLocalizationsEdit  lcl.ILabeledEdit
	LSUIElementText            lcl.ILabel
	LSUIElementBox             lcl.IComboBox
	LSMinimumSystemVersionText lcl.ILabel
	LSMinimumSystemVersionBox  lcl.IComboBox

	// linux

	// 操作按钮
	cancelBtn *wg.TButton
	saveBtn   *wg.TButton

	statusBar lcl.IStatusBar
}

func (m *TConfigProjectForm) FormCreate(sender lcl.IObject) {
	logs.Debug("TConfigProjectForm FormCreate")
	m.SetName("ConfigProjectForm")
	m.SetCaption(metadata.GI18n.Dict("ConfigProjectForm.Caption"))
	m.SetWidth(configProjectFormWidth)
	m.SetHeight(configProjectFormHeight)
	m.SetVisible(false)
	m.SetDoubleBuffered(true)
	m.SetBorderIcons(types.NewSet(types.BiSystemMenu))
	m.box = lcl.NewPanel(m)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetAlign(types.AlClient)
	m.box.SetColor(colors.ClWhite)
	m.box.SetParent(m)

	m.statusBar = lcl.NewStatusBar(m)
	m.statusBar.SetParent(m)
	m.statusBar.SetAutoHint(true)

	m.SetOnShow(m.onShow)
	m.initComponents()
	SetWindowCenterByMainWindow(m)

	//(&hook.TWindowHook{Form: m}).Hook()
}

func (m *TConfigProjectForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	m.closing = true
}

func (m *TConfigProjectForm) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	*closeAction = types.CaFree
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
		//constr := m.Constraints()
		//constr.SetMaxWidth(configProjectFormWidth)
		//constr.SetMaxHeight(configProjectFormHeight + addSize)
		//constr.SetMinWidth(configProjectFormWidth)
		//constr.SetMinHeight(configProjectFormHeight + addSize)
		// 初始时设置图标
		m.appIconData = bean.GProject.AppOption.Icon
		go func() {
			// 预览
			if bean.GProject.AppOption.Icon.Data != nil {
				// 预览图标过大, 需要绽放
				tempIconData := bean.GProject.AppOption.Icon.Data
				if bean.GProject.AppOption.Icon.W > 128 || bean.GProject.AppOption.Icon.H > 128 {
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
	textLeft := int32(70)

	m.font = lcl.NewFont()
	m.font.SetName("微软雅黑")
	m.font.SetCharSet(font.CHINESEBIG5_CHARSET)

	m.appConfigTitle = lcl.NewLabel(m)
	m.appConfigTitle.SetName("ConfigProjectFormAppConfigTitle")
	m.appConfigTitle.SetLeft(10)
	m.appConfigTitle.SetTop(10)
	m.appConfigTitle.SetCaption(metadata.GI18n.Dict("ConfigProjectFormAppConfigTitle.Caption"))
	m.appConfigTitle.SetFont(m.font)
	m.appConfigTitle.Font().SetSize(10)
	m.appConfigTitle.SetParent(m.box)

	baseTop := int32(40)
	nextTop := func(v int32) int32 {
		baseTop += v
		return baseTop
	}
	{
		m.appTitleEdit = lcl.NewLabeledEdit(m)
		m.appTitleEdit.SetName("ConfigProjectFormAppTitleEdit")
		m.appTitleEdit.EditLabel().SetCaption(metadata.GI18n.Dict("ConfigProjectFormAppTitleEdit.EditLabel.Caption"))
		m.appTitleEdit.SetBounds(textLeft, nextTop(0), 320, 30)
		m.appTitleEdit.SetFont(m.font)
		m.appTitleEdit.SetTextHint("MyEnergyApp")
		m.appTitleEdit.SetText(bean.GProject.AppOption.Title)
		m.appTitleEdit.SetLabelPosition(types.LpLeft)
		m.appTitleEdit.SetParent(m.box)
	}

	{
		m.appIdEdit = lcl.NewLabeledEdit(m)
		m.appIdEdit.SetName("ConfigProjectFormAppIdEdit")
		m.appIdEdit.EditLabel().SetCaption(metadata.GI18n.Dict("ConfigProjectFormAppIdEdit.EditLabel.Caption"))
		m.appIdEdit.SetBounds(textLeft, nextTop(35), 320, 30)
		m.appIdEdit.SetFont(m.font)
		m.appIdEdit.SetTextHint("company.product.app")
		m.appIdEdit.SetText(bean.GProject.AppOption.Id)
		m.appIdEdit.SetLabelPosition(types.LpLeft)
		m.appIdEdit.SetParent(m.box)

		m.appDescEdit = lcl.NewLabeledEdit(m)
		m.appDescEdit.SetName("ConfigProjectFormAppDescEdit")
		m.appDescEdit.EditLabel().SetCaption(metadata.GI18n.Dict("ConfigProjectFormAppDescEdit.EditLabel.Caption"))
		m.appDescEdit.SetBounds(textLeft, nextTop(35), 320, 30)
		m.appDescEdit.SetFont(m.font)
		m.appDescEdit.SetTextHint("your application description.")
		m.appDescEdit.SetText(bean.GProject.AppOption.Desc)
		m.appDescEdit.SetLabelPosition(types.LpLeft)
		m.appDescEdit.SetParent(m.box)

		m.appVersionEdit = lcl.NewLabeledEdit(m)
		m.appVersionEdit.SetName("ConfigProjectFormAppVersionEdit")
		m.appVersionEdit.EditLabel().SetCaption(metadata.GI18n.Dict("ConfigProjectFormAppVersionEdit.EditLabel.Caption"))
		m.appVersionEdit.SetBounds(textLeft, nextTop(35), 100, 30)
		m.appVersionEdit.SetFont(m.font)
		m.appVersionEdit.SetTextHint("1.2.3.4")
		m.appVersionEdit.SetText(bean.GProject.AppOption.Version)
		m.appVersionEdit.SetLabelPosition(types.LpLeft)
		m.appVersionEdit.SetParent(m.box)

		m.appCopyrightEdit = lcl.NewLabeledEdit(m)
		m.appCopyrightEdit.SetName("ConfigProjectFormAppCopyrightEdit")
		m.appCopyrightEdit.EditLabel().SetCaption(metadata.GI18n.Dict("ConfigProjectFormAppCopyrightEdit.EditLabel.Caption"))
		m.appCopyrightEdit.SetBounds(m.appVersionEdit.Left()+m.appVersionEdit.Width()+left+35,
			m.appVersionEdit.Top(), 175, 30)
		m.appCopyrightEdit.SetFont(m.font)
		m.appCopyrightEdit.SetTextHint("Copyright (C)")
		m.appCopyrightEdit.SetText(bean.GProject.AppOption.Copyright)
		m.appCopyrightEdit.SetLabelPosition(types.LpLeft)
		m.appCopyrightEdit.SetParent(m.box)

		m.appIconBtn = wg.NewButton(m)
		m.appIconBtn.SetName("ConfigProjectFormAppIconBtn")
		m.appIconBtn.SetIconFormBytes(resources.Images("button/upload_64x64.png"))
		//m.appIconBtn.SetIconCloseFormBytes(resources.Images("button/remove_16x16.png"))
		//m.appIconBtn.SetIconCloseHighlightFormBytes(resources.Images("button/remove_16x16_highlight.png"))
		m.appIconBtn.SetRadius(3)
		appIconRect := types.TRect{Left: m.Width() - 158, Top: 35}
		appIconRect.SetWidth(145)
		appIconRect.SetHeight(145)
		m.appIconBtn.TextOffSetY = 50
		m.appIconBtn.SetBoundsRect(appIconRect)
		m.appIconBtn.SetCaption(metadata.GI18n.Dict("ConfigProjectFormAppIconBtn.Caption"))
		m.appIconBtn.SetHint(metadata.GI18n.Dict("ConfigProjectFormAppIconBtn.Hint"))
		m.appIconBtn.SetShowHint(true)
		m.appIconBtn.SetFont(m.font)
		m.appIconBtn.SetBorderColor(wg.BbdNone, colors.RGBToColor(91, 155, 213))
		m.appIconBtn.SetBorderWidth(wg.BbdNone, 1)
		m.appIconBtn.SetColor(0xF3F4F6)
		m.appIconBtn.SetCursor(types.CrHandPoint)
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
		m.platformTitle.SetName("ConfigProjectFormPlatformTitle")
		m.platformTitle.SetLeft(10)
		m.platformTitle.SetTop(nextTop(30))
		m.platformTitle.SetCaption(metadata.GI18n.Dict("ConfigProjectFormPlatformTitle.Caption"))
		m.platformTitle.SetFont(m.font)
		m.platformTitle.Font().SetSize(10)
		m.platformTitle.SetParent(m.box)
	}

	{
		type Button struct {
			iconDefault []byte
			iconActive  []byte
		}
		buttons := tool.NewHashMap[string, *Button]()
		buttons.Add("Windows", &Button{
			iconDefault: resources.Images("button/windows_16x16.png"),
			iconActive:  resources.Images("button/windows_white_16x16.png"),
		})
		buttons.Add("MacOS", &Button{
			iconDefault: resources.Images("button/macos_16x16.png"),
			iconActive:  resources.Images("button/macos_white_16x16.png"),
		})
		buttons.Add("Linux", &Button{
			iconDefault: resources.Images("button/linux_16x16.png"),
			iconActive:  resources.Images("button/linux_white_16x16.png"),
		})

		m.platformTab = wg.NewTab(m)
		m.platformTab.Margin = 2
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
					page.Button().SetColor(tabActiveBgColor)
					page.Button().Font().SetColor(tabActiveTextColor)
					page.Button().SetBorderColor(wg.BbdNone, tabActiveBorderColor)
				} else {
					page.Button().SetColor(tabNoActiveBgColor)
					page.Button().Font().SetColor(tabNoActiveTextColor)
					page.Button().SetBorderColor(wg.BbdNone, tabNoActiveBorderColor)
				}
			}
		})
		// 设置标签按钮样式
		setTabPageStyle := func(page *wg.TPage) {
			page.SetTop(40)
			page.SetHeight(m.platformTab.Height() - 40)
			page.SetColor(m.platformTab.Color()) // 设置背景色
			page.Button().SetWidth(95)
			page.Button().SetHeight(30)
			page.Button().SetLeft(10)
			page.Button().SetRadius(0)
			page.Button().SetCursor(types.CrHandPoint)
		}

		m.platformTabPageWindows = m.platformTab.NewPage()
		m.platformTabPageWindows.SetCaption("Windows")
		m.platformTabPageWindows.Button().SetIconFavoriteFormBytes(buttons.Get("Windows").iconDefault)
		setTabPageStyle(m.platformTabPageWindows)
		m.initWindowsOptions()

		m.platformTabPageMacOS = m.platformTab.NewPage()
		m.platformTabPageMacOS.SetCaption("MacOS")
		m.platformTabPageMacOS.Button().SetIconFavoriteFormBytes(buttons.Get("MacOS").iconDefault)
		setTabPageStyle(m.platformTabPageMacOS)
		m.initMacOSOptions()

		m.platformTabPageLinux = m.platformTab.NewPage()
		m.platformTabPageLinux.SetCaption("Linux")
		m.platformTabPageLinux.Button().SetIconFavoriteFormBytes(buttons.Get("Linux").iconDefault)
		setTabPageStyle(m.platformTabPageLinux)
		m.initLinuxOptions()

		if tool.IsWindows {
			m.platformTabPageWindows.SetActive(true)
		} else if tool.IsDarwin {
			m.platformTabPageMacOS.SetActive(true)
		} else if tool.IsLinux {
			m.platformTabPageLinux.SetActive(true)
		}
	}

	{
		cancelBtnRect := types.TRect{Left: 400, Top: 525}
		cancelBtnRect.SetWidth(60)
		cancelBtnRect.SetHeight(25)
		m.cancelBtn = wg.NewButton(m)
		m.cancelBtn.SetName("ConfigProjectFormCancelBtn")
		m.cancelBtn.SetText(metadata.GI18n.Dict("ConfigProjectFormCancelBtn.Caption"))
		m.cancelBtn.Font().SetSize(8)
		m.cancelBtn.SetRadius(3)
		m.cancelBtn.SetBoundsRect(cancelBtnRect)
		m.cancelBtn.SetColor(grayBtnColor)
		m.cancelBtn.SetCursor(types.CrHandPoint)
		m.cancelBtn.SetParent(m.box)
		m.cancelBtn.SetOnClick(m.closeClick)

		saveBtnRect := types.TRect{Left: cancelBtnRect.Left + cancelBtnRect.Width() + 20, Top: cancelBtnRect.Top}
		saveBtnRect.SetWidth(60)
		saveBtnRect.SetHeight(25)
		m.saveBtn = wg.NewButton(m)
		m.saveBtn.SetName("ConfigProjectFormSaveBtn")
		m.saveBtn.SetText(metadata.GI18n.Dict("ConfigProjectFormSaveBtn.Caption"))
		m.saveBtn.Font().SetSize(8)
		m.saveBtn.Font().SetColor(colors.ClWhite)
		m.saveBtn.SetRadius(3)
		m.saveBtn.SetBoundsRect(saveBtnRect)
		m.saveBtn.SetColor(blueBtnColor)
		m.saveBtn.SetCursor(types.CrHandPoint)
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
	event.ConsoleWriteInfo("Project Configuration - Save")
	if m.validateInputs() {
		event.ConsoleWriteError("Input validation failed on save")
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
		// 创建 本地语言文件
		createAppLocalizations()
		// 更新图标
		updateWindowICON()
		event.ConsoleWriteInfo("Project Configuration - Save Completed")
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
		title = bean.GProject.AppOption.Title
	}
	return title
}

func (m *TConfigProjectForm) AppCopyright() string {
	copyright := m.appCopyrightEdit.Text()
	if copyright == "" {
		copyright = bean.GProject.AppOption.Copyright
	}
	return copyright
}

func (m *TConfigProjectForm) AppId() string {
	id := m.appIdEdit.Text()
	if id == "" {
		id = bean.GProject.AppOption.Id
	}
	return id
}

func (m *TConfigProjectForm) AppDesc() string {
	desc := m.appDescEdit.Text()
	if desc == "" {
		desc = bean.GProject.AppOption.Desc
	}
	return desc
}

func (m *TConfigProjectForm) AppVersion() string {
	appVersion := m.appVersionEdit.Text()
	if appVersion == "" {
		appVersion = bean.GProject.AppOption.Version
	}
	return appVersion
}

func (m *TConfigProjectForm) AppBundleName() string {
	bundleName := m.CFBundleNameEdit.Text()
	if bundleName == "" {
		bundleName = bean.GProject.Name
	}
	return bundleName
}

// AppBundleExecutable 从构建配置里获取
//
//	macOS 应用的主可执行文件名称
func (m *TConfigProjectForm) AppBundleExecutable() string {
	return bean.GProject.BuildOption.BuildFileName
}

func (m *TConfigProjectForm) AppBundleLocalizations() []string {
	locals := m.CFBundleLocalizationsEdit.Text()
	bundleLocalizations := tool.Split(locals, ",")
	if len(bundleLocalizations) == 0 {
		bundleLocalizations = []string{bean.GProject.AppOption.Lang}
	}
	return bundleLocalizations
}

func (m *TConfigProjectForm) AppLSUIElement() bool {
	LSUIElement := m.LSUIElementBox.ItemIndex() != int32(bean.MacOSUIElementListNo)
	return LSUIElement
}

func (m *TConfigProjectForm) AppLSMinimumSystemVersion() string {
	switch bean.LSMinimumSystemVersion(m.LSMinimumSystemVersionBox.ItemIndex()) {
	case bean.LSMinimumSystemVersion_10_15:
		return "10.15"
	case bean.LSMinimumSystemVersion_11_0:
		return "11.0"
	}
	return "10.15"
}
