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
	"context"
	"fmt"
	cmdcef "github.com/energye/designer/cmd/cef"
	"github.com/energye/designer/resources/metadata"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

var (
	chromiumDirFormWidth  = int32(500)
	chromiumDirFormHeight = int32(200)
)

// 下载状态
type downloadState int32

const (
	downloadIdle       downloadState = iota // 空闲 - 按钮显示"确 定"
	downloadRunning                         // 下载中 - 按钮显示"停 止"
	downloadPaused                          // 已暂停 - 按钮显示"继 续"
	downloadExtracting                      // 解压中 - 按钮显示"停 止"
	downloadCompleted                       // 已完成
)

// runChromiumDirConfig 运行 CEF 框架配置窗口
func runChromiumDirConfig() {
	lcl.RunOnMainThreadAsync(func(id uint32) {
		form := NewChromiumDirForm()
		form.ShowModal()
	})
}

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
	dlState   downloadState
	dlCancel  context.CancelFunc
	dlMu      sync.Mutex
	dlVersion string // 当前正在下载的版本, 用于暂停后检测版本变更
	Confirmed bool   // 用户是否点击了"完成"确认按钮
	Version   string // 已安装完成的 CEF 版本

	// 目录设置
	dirText lcl.ILabel
	dirEdit lcl.ILabeledEdit
	dirBtn  *wg.TButton

	// 系统/架构/版本选择
	osText       lcl.ILabel
	osBox        lcl.IComboBox
	archText     lcl.ILabel
	archBox      lcl.IComboBox
	versionText  lcl.ILabel
	versionBox   lcl.IComboBox
	versionList  []string        // 下拉框索引对应的真实版本号
	installedSet map[string]bool // 已安装版本集合

	// 下载进度
	progressBar lcl.IProgressBar
	statusLabel lcl.ILabel

	// 操作按钮
	defaultBtn *wg.TButton
	cancelBtn  *wg.TButton
	confirmBtn *wg.TButton
}

func (m *TChromiumDirForm) FormCreate(sender lcl.IObject) {
	logs.Debug("TChromiumDirForm FormCreate")
	m.SetName("ChromiumDirForm")
	m.SetCaption(metadata.GI18n.Dict("ChromiumDirForm.Caption"))
	m.SetWidth(chromiumDirFormWidth)
	m.SetHeight(chromiumDirFormHeight)
	m.SetVisible(false)
	m.SetDoubleBuffered(true)
	m.SetBorderIcons(types.NewSet(types.BiSystemMenu))
	m.SetColor(colors.ClWhite)
	SetWindowCenterByMainWindow(m)

	m.selectDir = lcl.NewSelectDirectoryDialog(m)
	m.selectDir.SetName("ChromiumDirFormSelectDir")
	m.selectDir.SetTitle(metadata.GI18n.Dict("ChromiumDirFormSelectDir.Title"))

	gTop := int32(0)
	nextTop := func(top int32) int32 {
		gTop += top
		return gTop
	}

	m.setupDirSection(nextTop)
	m.setupVersionSection(nextTop)
	m.setupProgressSection(nextTop)
	m.setupActionButtons(nextTop)

	if m.isVersionInstalled() {
		m.confirmBtn.SetText(metadata.GI18n.Dict("ChromiumDirFormConfirmBtn.TextUse"))
	}
}

// setupDirSection 创建目录输入区域 (说明文字 + 输入框 + 浏览按钮)
func (m *TChromiumDirForm) setupDirSection(nextTop func(int32) int32) {
	m.dirText = lcl.NewLabel(m)
	m.dirText.SetCaption("ChromiumDirFormDirText")
	m.dirText.SetLeft(15)
	m.dirText.SetTop(nextTop(15))
	m.dirText.SetCaption(metadata.GI18n.Dict("ChromiumDirFormDirText.Caption"))
	m.dirText.SetParent(m)

	m.dirEdit = lcl.NewLabeledEdit(m)
	m.dirEdit.SetName("ChromiumDirFormDirEdit")
	m.dirEdit.SetLeft(100)
	m.dirEdit.SetTop(nextTop(25))
	m.dirEdit.SetWidth(330)
	m.dirEdit.SetDoubleBuffered(true)
	m.dirEdit.SetTextHint(config.Config.Chromium.DefaultDir())
	m.dirEdit.SetLabelPosition(types.LpLeft)
	if config.Config.Chromium.Dir != "" {
		m.dirEdit.SetText(config.Config.Chromium.Dir)
	}
	m.dirEdit.EditLabel().SetCaption(metadata.GI18n.Dict("ChromiumDirFormDirEdit.EditLabel.Caption"))
	m.dirEdit.SetParent(m)

	dirBtnRect := types.TRect{Left: m.dirEdit.Left() + m.dirEdit.Width() + 5, Top: m.dirEdit.Top()}
	dirBtnRect.SetWidth(35)
	if tool.IsLinux {
		dirBtnRect.SetHeight(35)
	} else {
		dirBtnRect.SetHeight(25)
	}
	m.dirBtn = wg.NewButton(m)
	m.dirBtn.SetIconFormBytes(resources.Images("menu/menu_project_open.png"))
	m.dirBtn.SetRadius(3)
	m.dirBtn.SetBoundsRect(dirBtnRect)
	m.dirBtn.SetColor(grayBtnColor)
	m.dirBtn.SetBorderColor(wg.BbdNone, grayBtnColor)
	m.dirBtn.SetCursor(types.CrHandPoint)
	m.dirBtn.SetParent(m)
	m.dirBtn.SetOnClick(m.dirBtnClick)
}

// setupVersionSection 创建系统/架构/版本选择区域 (同一行: OS + ARCH + CEF Version)
func (m *TChromiumDirForm) setupVersionSection(nextTop func(int32) int32) {
	rowTop := nextTop(40)

	// Version 标签 + 下拉框
	m.versionText = lcl.NewLabel(m)
	m.versionText.SetName("ChromiumDirFormVersionText")
	m.versionText.SetLeft(0)
	m.versionText.SetTop(rowTop)
	m.versionText.SetAutoSize(false)
	m.versionText.SetWidth(96)
	m.versionText.SetAlignment(types.TaRightJustify)
	m.versionText.SetCaption(metadata.GI18n.Dict("ChromiumDirFormVersionText.Caption"))
	m.versionText.SetParent(m)

	m.versionBox = lcl.NewComboBox(m)
	m.versionBox.SetName("ChromiumDirFormVersionBox")
	m.versionBox.SetLeft(100)
	m.versionBox.SetTop(rowTop)
	m.versionBox.SetWidth(100)
	m.versionBox.SetReadOnly(true)
	m.versionBox.AnchorSideTop().SetControl(m.versionText)
	m.versionBox.AnchorSideTop().SetSide(types.AsrCenter)
	m.versionBox.SetStyle(types.CsDropDownList)
	m.versionBox.SetBorderStyle(types.BsSingle)
	m.versionBox.SetOnChange(m.onVersionChange)
	m.versionBox.SetParent(m)

	// OS 标签 + 下拉框
	m.osText = lcl.NewLabel(m)
	m.osText.SetName("ChromiumDirFormOSText")
	m.osText.SetLeft(m.versionBox.Left() + m.versionBox.Width() + 10)
	m.osText.SetTop(rowTop)
	m.osText.SetAutoSize(false)
	m.osText.SetWidth(35)
	m.osText.SetAlignment(types.TaRightJustify)
	m.osText.SetCaption(metadata.CEFFormOsLabelText)
	m.osText.SetParent(m)

	m.osBox = lcl.NewComboBox(m)
	m.osBox.SetName("ChromiumDirFormOSBox")
	m.osBox.SetLeft(m.osText.Left() + m.osText.Width() + 4)
	m.osBox.SetTop(rowTop)
	m.osBox.SetWidth(90)
	m.osBox.SetReadOnly(true)
	m.osBox.AnchorSideTop().SetControl(m.osText)
	m.osBox.AnchorSideTop().SetSide(types.AsrCenter)
	m.osBox.SetStyle(types.CsDropDownList)
	m.osBox.SetBorderStyle(types.BsSingle)
	for _, osName := range cmdcef.SupportedOSList {
		m.osBox.Items().Add(osName)
	}
	m.osBox.SetItemIndex(0) // 默认 windows
	m.osBox.SetOnChange(m.onOSChange)
	m.osBox.SetParent(m)

	// ARCH 标签 + 下拉框
	m.archText = lcl.NewLabel(m)
	m.archText.SetName("ChromiumDirFormARCHText")
	m.archText.SetLeft(m.osBox.Left() + m.osBox.Width() + 10)
	m.archText.SetTop(rowTop)
	m.archText.SetAutoSize(false)
	m.archText.SetWidth(35)
	m.archText.SetAlignment(types.TaRightJustify)
	m.archText.SetCaption(metadata.CEFFormArchLabelText)
	m.archText.SetParent(m)

	m.archBox = lcl.NewComboBox(m)
	m.archBox.SetName("ChromiumDirFormARCHBox")
	m.archBox.SetLeft(m.archText.Left() + m.archText.Width() + 4)
	m.archBox.SetTop(rowTop)
	m.archBox.SetWidth(80)
	m.archBox.SetReadOnly(true)
	m.archBox.AnchorSideTop().SetControl(m.archText)
	m.archBox.AnchorSideTop().SetSide(types.AsrCenter)
	m.archBox.SetStyle(types.CsDropDownList)
	m.archBox.SetBorderStyle(types.BsSingle)
	m.archBox.SetOnChange(m.onArchChange)
	m.archBox.SetParent(m)

	// 初始化: 设置默认 OS 并联动 ARCH
	m.initOSArchDefault()
	// 填充版本列表
	m.populateVersionList()
}

// setupProgressSection 创建下载进度区域 (进度条 + 状态标签)
func (m *TChromiumDirForm) setupProgressSection(nextTop func(int32) int32) {
	m.progressBar = lcl.NewProgressBar(m)
	m.progressBar.SetName("ChromiumDirFormProgressBar")
	m.progressBar.SetLeft(100)
	m.progressBar.SetTop(nextTop(35))
	m.progressBar.SetWidth(370)
	m.progressBar.SetHeight(20)
	m.progressBar.SetParent(m)
	m.progressBar.SetVisible(false)

	m.statusLabel = lcl.NewLabel(m)
	m.statusLabel.SetName("ChromiumDirFormStatusLabel")
	m.statusLabel.SetLeft(100)
	m.statusLabel.SetTop(nextTop(22))
	m.statusLabel.SetCaption(metadata.CEFFormStatusLabelCaption)
	m.statusLabel.Font().SetColor(colors.RGBToColor(128, 128, 128))
	m.statusLabel.Font().SetSize(8)
	m.statusLabel.SetParent(m)
}

// setupActionButtons 创建底部操作按钮 (默认目录、取消、确定)
func (m *TChromiumDirForm) setupActionButtons(nextTop func(int32) int32) {
	btnTop := nextTop(25)

	defaultBtnRect := types.TRect{Left: 100, Top: btnTop}
	defaultBtnRect.SetWidth(100)
	defaultBtnRect.SetHeight(25)
	m.defaultBtn = newGrayButton(m, defaultBtnRect, metadata.GI18n.Dict("ChromiumDirFormDefaultBtn.Text"), m.defaultBtnClick)
	m.defaultBtn.SetName("ChromiumDirFormDefaultBtn")

	cancelBtnRect := types.TRect{Left: 325, Top: btnTop}
	cancelBtnRect.SetWidth(60)
	cancelBtnRect.SetHeight(25)
	m.cancelBtn = newGrayButton(m, cancelBtnRect, metadata.GI18n.Dict("ChromiumDirFormCancelBtn.Text"), m.cancelBtnClick)
	m.cancelBtn.SetName("ChromiumDirFormCancelBtn")

	confirmBtnRect := types.TRect{Left: cancelBtnRect.Left + 60 + 20, Top: btnTop}
	confirmBtnRect.SetWidth(60)
	confirmBtnRect.SetHeight(25)
	m.confirmBtn = newBlueButton(m, confirmBtnRect, metadata.GI18n.Dict("ChromiumDirFormConfirmBtn.Text"), m.confirmBtnClick)
	m.confirmBtn.SetName("ChromiumDirFormConfirmBtn")
}

// ==================== OS/ARCH/版本 与 URL ====================

// selectedOS 返回当前选中的系统名
func (m *TChromiumDirForm) selectedOS() string {
	idx := m.osBox.ItemIndex()
	if idx < 0 || int(idx) >= len(cmdcef.SupportedOSList) {
		return runtime.GOOS
	}
	return cmdcef.SupportedOSList[idx]
}

// selectedArch 返回当前选中的架构名
func (m *TChromiumDirForm) selectedArch() string {
	idx := m.archBox.ItemIndex()
	osName := m.selectedOS()
	archs := cmdcef.OSArchMap[osName]
	if idx < 0 || int(idx) >= len(archs) {
		return runtime.GOARCH
	}
	return archs[idx]
}

// osArchVersion 返回 os_arch_version 格式的标识, 用于目录名和清单 key
func (m *TChromiumDirForm) osArchVersion(version string) string {
	return fmt.Sprintf("%s_%s_%s", m.selectedOS(), m.selectedArch(), version)
}

// initOSArchDefault 初始化 OS/ARCH 下拉框默认值, 优先从已配置的 CEF 版本中解析, 否则使用当前系统架构
func (m *TChromiumDirForm) initOSArchDefault() {
	// 从已配置的 CEF 版本中解析 OS 和 ARCH
	targetOS, targetArch := m.resolveConfiguredOSArch()

	// 设置默认 OS
	osIdx := 0
	for i, osName := range cmdcef.SupportedOSList {
		if osName == targetOS {
			osIdx = i
			break
		}
	}
	m.osBox.SetItemIndex(int32(osIdx))

	// 填充对应架构并设置默认
	m.populateArchListWithDefault(targetArch)
}

// resolveConfiguredOSArch 从已配置的 CEF 版本中解析 OS 和 ARCH, 否则返回当前系统架构
func (m *TChromiumDirForm) resolveConfiguredOSArch() (osName, arch string) {
	version := config.Config.Chromium.Version
	if bean.GProject != nil && bean.GProject.GUIRenderFramework == bean.GUIRenderFramework_CEF {
		version = bean.GProject.FrameworkVersion
	}
	if version != "" {
		parts := strings.SplitN(version, "_", 3)
		if len(parts) >= 2 {
			cfgOS := parts[0]
			cfgArch := parts[1]
			// 验证 OS 和 ARCH 是否有效
			osValid := false
			for _, os := range cmdcef.SupportedOSList {
				if os == cfgOS {
					osValid = true
					break
				}
			}
			archValid := false
			if archs, ok := cmdcef.OSArchMap[cfgOS]; ok {
				for _, a := range archs {
					if a == cfgArch {
						archValid = true
						break
					}
				}
			}
			if osValid && archValid {
				return cfgOS, cfgArch
			}
		}
	}
	return runtime.GOOS, runtime.GOARCH
}

// populateArchList 根据当前选中的 OS 填充架构下拉框, 默认选中当前系统架构
func (m *TChromiumDirForm) populateArchList() {
	m.populateArchListWithDefault(runtime.GOARCH)
}

// populateArchListWithDefault 根据当前选中的 OS 填充架构下拉框, 并选中指定的默认架构
func (m *TChromiumDirForm) populateArchListWithDefault(defaultArch string) {
	osName := m.selectedOS()
	archs := cmdcef.OSArchMap[osName]
	m.archBox.Items().Clear()
	for _, arch := range archs {
		m.archBox.Items().Add(arch)
	}
	// 默认选中指定架构
	archIdx := 0
	for i, arch := range archs {
		if arch == defaultArch {
			archIdx = i
			break
		}
	}
	if archIdx >= len(archs) {
		archIdx = 0
	}
	m.archBox.SetItemIndex(int32(archIdx))
}

// onOSChange 系统切换时联动更新架构下拉框和版本列表
func (m *TChromiumDirForm) onOSChange(sender lcl.IObject) {
	m.populateArchList()
	// 保留当前版本选择, 重建版本列表以更新 ✓ 标记
	curVer := m.selectedVersion()
	m.populateVersionList()
	// 恢复之前的版本选择
	for i, ver := range m.versionList {
		if ver == curVer {
			m.versionBox.SetItemIndex(int32(i))
			break
		}
	}
	m.updateConfirmBtnTextOnly()
}

// onArchChange 架构切换时重建版本列表并更新按钮文字
func (m *TChromiumDirForm) onArchChange(sender lcl.IObject) {
	curVer := m.selectedVersion()
	m.populateVersionList()
	for i, ver := range m.versionList {
		if ver == curVer {
			m.versionBox.SetItemIndex(int32(i))
			break
		}
	}
	m.updateConfirmBtnTextOnly()
}

// onVersionChange 版本切换时仅更新按钮文字, 不重建列表
func (m *TChromiumDirForm) onVersionChange(sender lcl.IObject) {
	m.updateConfirmBtnTextOnly()
}

// updateConfirmBtnTextOnly 仅更新按钮文字, 不重建版本列表
func (m *TChromiumDirForm) updateConfirmBtnTextOnly() {
	if m.isVersionInstalled() {
		m.confirmBtn.SetText(metadata.GI18n.Dict("ChromiumDirFormConfirmBtn.TextUse"))
	} else {
		m.confirmBtn.SetText(metadata.GI18n.Dict("ChromiumDirFormConfirmBtn.Text"))
	}
}

// populateVersionList 填充版本列表
func (m *TChromiumDirForm) populateVersionList() {
	m.installedSet = m.buildInstalledSet()
	m.versionList = m.sortedVersions()
	m.versionBox.Items().Clear()
	for _, ver := range m.versionList {
		label := ver
		oav := m.osArchVersion(ver)
		if m.installedSet[oav] {
			label += " ✓"
		}
		m.versionBox.Items().Add(label)
	}
	version := config.Config.Chromium.Version
	if bean.GProject != nil && bean.GProject.GUIRenderFramework == bean.GUIRenderFramework_CEF {
		version = bean.GProject.FrameworkVersion
	}
	if m.versionBox.Items().Count() > 0 {
		// 尝试恢复上次选中的版本
		savedVersion := version
		selectIdx := int32(0)
		if savedVersion != "" {
			for i, ver := range m.versionList {
				if m.osArchVersion(ver) == savedVersion {
					selectIdx = int32(i)
					break
				}
			}
		}
		m.versionBox.SetItemIndex(selectIdx)
	}
}

// sortedVersions 从 DesignerConfig.Chromium 读取版本号并排序
func (m *TChromiumDirForm) sortedVersions() []string {
	return cmdcef.Versions()
}

// buildInstalledSet 返回已安装版本的集合, key 为 os_arch_version
func (m *TChromiumDirForm) buildInstalledSet() map[string]bool {
	installed := make(map[string]bool)
	for _, oav := range cmdcef.InstalledVersions(m.selectedOS(), m.selectedArch()) {
		installed[oav] = true
	}
	return installed
}

// selectedVersion 通过索引返回真实版本号
func (m *TChromiumDirForm) selectedVersion() string {
	idx := m.versionBox.ItemIndex()
	if idx < 0 || int(idx) >= len(m.versionList) {
		return ""
	}
	return m.versionList[idx]
}

// isVersionInstalled 判断当前选中版本是否已安装 (基于 os_arch_version)
func (m *TChromiumDirForm) isVersionInstalled() bool {
	ver := m.selectedVersion()
	if ver == "" {
		return false
	}
	oav := m.osArchVersion(ver)
	return m.installedSet[oav]
}

// updateConfirmBtnText 根据安装状态更新按钮文字 (仅更新按钮, 不重建列表)
func (m *TChromiumDirForm) updateConfirmBtnText() {
	m.updateConfirmBtnTextOnly()
}
func (m *TChromiumDirForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	m.dlMu.Lock()
	defer m.dlMu.Unlock()
	if m.dlState == downloadRunning || m.dlState == downloadExtracting {
		if m.dlCancel != nil {
			m.dlCancel()
		}
	}
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

// ==================== 按钮状态机 ====================

// confirmBtnClick 确认/下载/停止/继续 按钮
func (m *TChromiumDirForm) confirmBtnClick(sender lcl.IObject) {
	m.dlMu.Lock()
	state := m.dlState
	m.dlMu.Unlock()

	// 空闲状态: 已安装直接确认, 未安装走下载
	if state == downloadIdle && m.isVersionInstalled() {
		oav := m.osArchVersion(m.selectedVersion())
		m.Version = oav
		// 记录当前 CEF 版本到全局配置
		if err := cmdcef.UseInstalled(oav, bean.GProject); err != nil {
			event.ConsoleWriteError("CEF configuration failed:", err.Error())
			return
		}
		m.Confirmed = true
		m.Close()
		return
	}

	switch state {
	case downloadIdle:
		m.startDownload()
	case downloadRunning:
		m.pauseDownload()
	case downloadPaused:
		m.resumeDownload()
	case downloadExtracting:
		m.stopExtract()
	}
}

// setDownloadState 设置下载状态 (需在主线程调用)
func (m *TChromiumDirForm) setDownloadState(state downloadState) {
	m.dlMu.Lock()
	m.dlState = state
	m.dlMu.Unlock()
	m.confirmBtn.SetEnabled(true)
	switch state {
	case downloadIdle:
		m.confirmBtn.SetText(metadata.GI18n.Dict("ChromiumDirFormConfirmBtn.Text"))
		m.confirmBtn.SetColor(blueBtnColor)
	case downloadRunning, downloadExtracting:
		m.confirmBtn.SetText(metadata.GI18n.Dict("ChromiumDirFormConfirmBtn.TextStop"))
		m.confirmBtn.SetColor(colors.RGBToColor(255, 127, 127))
	case downloadPaused:
		m.confirmBtn.SetText(metadata.GI18n.Dict("ChromiumDirFormConfirmBtn.TextRun"))
		m.confirmBtn.SetColor(blueBtnColor)
	case downloadCompleted:
		m.confirmBtn.SetText(metadata.GI18n.Dict("ChromiumDirFormConfirmBtn.TextSuccess"))
		m.confirmBtn.SetColor(colors.RGBToColor(46, 204, 113))
	}
}

// resetToIdle 恢复窗口到默认空闲状态
func (m *TChromiumDirForm) resetToIdle() {
	m.setDownloadState(downloadIdle)
	m.progressBar.SetVisible(false)
	m.statusLabel.SetVisible(true)
	m.statusLabel.SetCaption(metadata.GI18n.Dict("ChromiumDirFormStatusLabel.Caption"))
	m.defaultBtn.SetEnabled(true)
	m.dirBtn.SetEnabled(true)
	m.osBox.SetEnabled(true)
	m.archBox.SetEnabled(true)
	m.versionBox.SetEnabled(true)
	m.confirmBtn.SetOnClick(m.confirmBtnClick)
}

// ==================== 下载流程 ====================

// startDownload 开始下载 CEF
func (m *TChromiumDirForm) startDownload() {
	dir := m.dirEdit.Text()
	if dir == "" {
		dir = config.Config.Chromium.DefaultDir()
		m.dirEdit.SetText(dir)
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		event.ConsoleWriteError("Invalid directory path:", err.Error())
		m.statusLabel.SetCaption(metadata.CEFFormStatusLabelCaptionInvalid)
		return
	}

	// 获取选中版本
	version := m.selectedVersion()
	if version == "" {
		event.ConsoleWriteError("Please select a CEF version")
		m.statusLabel.SetCaption(metadata.CEFFormStatusLabelCaptionSelectCEF)
		return
	}

	osName := m.selectedOS()
	arch := m.selectedArch()
	downloadURL := cmdcef.BuildDownloadURL(version, osName, arch)
	if downloadURL == "" {
		event.ConsoleWriteError("Download URL not found for version:", version)
		m.statusLabel.SetCaption(metadata.CEFFormStatusLabelCaptionURLNotFound)
		return
	}

	// 创建目录
	if !tool.IsExist(absDir) {
		if err = os.MkdirAll(absDir, os.ModePerm); err != nil {
			event.ConsoleWriteError("Failed to create directory:", err.Error())
			m.statusLabel.SetCaption(metadata.CEFFormStatusLabelCaptionFailedCreateDirectory + ": " + err.Error())
			return
		}
	}

	// 更新配置
	m.dlVersion = version

	// 锁定 UI
	m.setDownloadState(downloadRunning)
	m.progressBar.SetVisible(true)
	m.progressBar.SetPosition(0)
	m.statusLabel.SetVisible(true)
	m.statusLabel.SetCaption(metadata.CEFFormStatusLabelCaptionCaptionFailedPreDown)
	m.defaultBtn.SetEnabled(false)
	m.dirBtn.SetEnabled(false)
	m.osBox.SetEnabled(false)
	m.archBox.SetEnabled(false)
	m.versionBox.SetEnabled(false)

	targetPath := filepath.Join(absDir, cmdcef.ArchiveFileName(version, osName, arch))
	event.ConsoleWriteInfo("Start downloading CEF:", downloadURL)
	event.ConsoleWriteInfo("Target:", targetPath)

	ctx, cancel := context.WithCancel(context.Background())
	m.dlMu.Lock()
	m.dlCancel = cancel
	m.dlMu.Unlock()
	go m.installCEF(ctx, cmdcef.InstallOptions{
		Dir:        absDir,
		Version:    version,
		OS:         osName,
		Arch:       arch,
		Project:    bean.GProject,
		OnProgress: m.handleInstallProgress,
	})
}

// handleInstallProgress 同步 cmd/cef 的安装进度到窗口状态。
func (m *TChromiumDirForm) handleInstallProgress(progress cmdcef.Progress) {
	if progress.Kind == cmdcef.ProgressInfo {
		if progress.Message != "" {
			event.ConsoleWriteInfo(progress.Message)
		}
		m.updateInstallInfo(progress.Message)
		return
	}
	if progress.Kind == cmdcef.ProgressDownload {
		m.updateInstallProgress(downloadRunning, progress, "Downloading...")
		return
	}
	if progress.Kind == cmdcef.ProgressExtract {
		m.updateInstallProgress(downloadExtracting, progress, "Extracting...")
	}
}

func (m *TChromiumDirForm) updateInstallInfo(message string) {
	lcl.RunOnMainThreadAsync(func(id uint32) {
		if m.closing {
			return
		}
		state := m.currentDownloadState()
		if state == downloadPaused || state == downloadCompleted {
			return
		}
		m.progressBar.SetVisible(true)
		if strings.HasPrefix(message, "Start downloading CEF:") {
			m.setDownloadState(downloadRunning)
			m.progressBar.SetPosition(0)
			m.statusLabel.SetCaption(metadata.CEFFormStatusLabelCaptionCaptionFailedPreDown)
			return
		}
		if strings.HasPrefix(message, "Extracting CEF:") {
			m.setDownloadState(downloadExtracting)
			m.progressBar.SetPosition(0)
			m.statusLabel.SetCaption(metadata.GI18n.Dict("ChromiumDirFormStatusLabel.CaptionSuccessDownloadUnZip"))
			return
		}
		if message != "" {
			m.statusLabel.SetCaption(message)
		}
	})
}

func (m *TChromiumDirForm) updateInstallProgress(state downloadState, progress cmdcef.Progress, fallback string) {
	current := progress.Current
	total := progress.Total
	message := progress.Message
	if message == "" {
		message = fallback
	}
	lcl.RunOnMainThreadAsync(func(id uint32) {
		if m.closing {
			return
		}
		currentState := m.currentDownloadState()
		if currentState == downloadPaused || currentState == downloadCompleted {
			return
		}
		if currentState == downloadExtracting && state == downloadRunning {
			return
		}
		m.setDownloadState(state)
		m.progressBar.SetVisible(true)
		m.statusLabel.SetVisible(true)
		if total > 0 {
			percent := int32(current * 100 / total)
			if percent < 0 {
				percent = 0
			}
			if percent > 100 {
				percent = 100
			}
			m.progressBar.SetPosition(percent)
			m.statusLabel.SetCaption(fmt.Sprintf("%s (%d%%)", message, percent))
			return
		}
		m.statusLabel.SetCaption(message)
	})
}

func (m *TChromiumDirForm) currentDownloadState() downloadState {
	m.dlMu.Lock()
	defer m.dlMu.Unlock()
	return m.dlState
}

// pauseDownload 暂停下载
func (m *TChromiumDirForm) pauseDownload() {
	m.dlMu.Lock()
	if m.dlCancel != nil {
		m.dlCancel()
		m.dlCancel = nil
	}
	m.dlMu.Unlock()

	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.setDownloadState(downloadPaused)
		m.statusLabel.SetCaption(metadata.CEFFormStatusLabelCaptionCaptionPaused)
		m.osBox.SetEnabled(true)
		m.archBox.SetEnabled(true)
		m.versionBox.SetEnabled(true)
	})
}

// stopExtract 停止解压, 恢复窗口到默认状态
func (m *TChromiumDirForm) stopExtract() {
	m.dlMu.Lock()
	if m.dlCancel != nil {
		m.dlCancel()
		m.dlCancel = nil
	}
	m.dlMu.Unlock()
	// 即时反馈: 按钮变为"停止中", 禁用防重复点击
	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.confirmBtn.SetText(metadata.GI18n.Dict("ChromiumDirFormConfirmBtn.TextStopping"))
		m.confirmBtn.SetEnabled(false)
	})
}

// resumeDownload 继续下载
func (m *TChromiumDirForm) resumeDownload() {
	m.startDownload()
}

// installCEF 执行 CEF 安装流程。
func (m *TChromiumDirForm) installCEF(ctx context.Context, options cmdcef.InstallOptions) {
	defer func() {
		m.dlMu.Lock()
		m.dlCancel = nil
		m.dlMu.Unlock()
	}()
	result, err := cmdcef.EnsureInstalled(ctx, options)
	if err != nil {
		if ctx.Err() != nil {
			m.onInstallCanceled()
			return
		}
		m.onDownloadError(err.Error())
		return
	}
	m.Version = result.OSArchVersion
	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.setDownloadState(downloadCompleted)
		m.progressBar.SetPosition(100)
		m.statusLabel.SetCaption(metadata.GI18n.Dict("ChromiumDirFormStatusLabel.CaptionSuccessInstall"))
		m.confirmBtn.SetEnabled(true)
		m.confirmBtn.SetText(metadata.GI18n.Dict("ChromiumDirFormConfirmBtn.TextSuccess"))
		m.confirmBtn.SetColor(colors.RGBToColor(46, 204, 113))
		m.confirmBtn.SetOnClick(m.onCompleteBtnClick)
	})
}

func (m *TChromiumDirForm) onInstallCanceled() {
	lcl.RunOnMainThreadAsync(func(id uint32) {
		if m.closing {
			return
		}
		state := m.currentDownloadState()
		if state == downloadExtracting {
			m.resetToIdle()
			m.statusLabel.SetCaption(metadata.GI18n.Dict("ChromiumDirFormStatusLabel.Caption"))
			return
		}
		if state == downloadRunning {
			m.setDownloadState(downloadPaused)
			m.statusLabel.SetCaption(metadata.CEFFormStatusLabelCaptionCaptionPaused)
			m.osBox.SetEnabled(true)
			m.archBox.SetEnabled(true)
			m.versionBox.SetEnabled(true)
		}
	})
}

func (m *TChromiumDirForm) onDownloadError(msg string) {
	event.ConsoleWriteError("CEF download error:", msg)
	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.resetToIdle()
		m.statusLabel.SetCaption(metadata.GI18n.Dict("ChromiumDirFormStatusLabel.CaptionFailDownload") +
			": " + msg)
	})
}

func (m *TChromiumDirForm) onCompleteBtnClick(sender lcl.IObject) {
	m.Confirmed = true
	m.Close()
}

// ==================== CEF 版本辅助函数 ====================

// GetInstalledCEFVersions 返回当前系统架构下已安装的 CEF 版本列表 (os_arch_version 格式)
func GetInstalledCEFVersions() []string {
	return cmdcef.InstalledVersions(runtime.GOOS, runtime.GOARCH)
}
