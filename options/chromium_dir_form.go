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
	"archive/tar"
	"compress/bzip2"
	"context"
	"fmt"
	"github.com/energye/designer/resources/metadata"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
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
	errExtractStopped     = fmt.Errorf("extract stopped by user")
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
	closing    bool
	selectDir  lcl.ISelectDirectoryDialog
	dlState    downloadState
	dlCancel   context.CancelFunc
	dlMu       sync.Mutex
	dlProgress int64
	dlTotal    int64
	dlVersion  string // 当前正在下载的版本, 用于暂停后检测版本变更
	dlStop     bool   // 解压停止信号
	Confirmed  bool   // 用户是否点击了"完成"确认按钮
	Version    string // 已安装完成的 CEF 版本

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
	for _, osName := range supportedOSList {
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

// supportedOSList 支持的操作系统列表
var supportedOSList = []string{"windows", "linux", "darwin"}

// osArchMap 每个系统支持的架构列表
var osArchMap = map[string][]string{
	"windows": {"amd64", "386"},
	"linux":   {"amd64", "386", "arm64", "arm"},
	"darwin":  {"amd64", "arm64"},
}

// cefOSArchMap Go runtime 架构名到 CEF 下载链接架构名的映射
var cefOSArchMap = map[string]map[string]string{
	"windows": {"amd64": "windows64", "386": "windows32"},
	"linux":   {"amd64": "linux64", "386": "linux32", "arm64": "linuxarm64", "arm": "linuxarm"},
	"darwin":  {"amd64": "macosx64", "arm64": "macosarm64"},
}

// selectedOS 返回当前选中的系统名
func (m *TChromiumDirForm) selectedOS() string {
	idx := m.osBox.ItemIndex()
	if idx < 0 || int(idx) >= len(supportedOSList) {
		return runtime.GOOS
	}
	return supportedOSList[idx]
}

// selectedArch 返回当前选中的架构名
func (m *TChromiumDirForm) selectedArch() string {
	idx := m.archBox.ItemIndex()
	osName := m.selectedOS()
	archs := osArchMap[osName]
	if idx < 0 || int(idx) >= len(archs) {
		return runtime.GOARCH
	}
	return archs[idx]
}

// osArchVersion 返回 os_arch_version 格式的标识, 用于目录名和清单 key
func (m *TChromiumDirForm) osArchVersion(version string) string {
	return fmt.Sprintf("%s_%s_%s", m.selectedOS(), m.selectedArch(), version)
}

// initOSArchDefault 初始化 OS/ARCH 下拉框默认值为当前系统和架构
func (m *TChromiumDirForm) initOSArchDefault() {
	// 设置默认 OS
	currentOS := runtime.GOOS
	osIdx := 0
	for i, osName := range supportedOSList {
		if osName == currentOS {
			osIdx = i
			break
		}
	}
	m.osBox.SetItemIndex(int32(osIdx))

	// 填充对应架构并设置默认
	m.populateArchList()
}

// populateArchList 根据当前选中的 OS 填充架构下拉框
func (m *TChromiumDirForm) populateArchList() {
	osName := m.selectedOS()
	archs := osArchMap[osName]
	m.archBox.Items().Clear()
	for _, arch := range archs {
		m.archBox.Items().Add(arch)
	}
	// 默认选中当前架构
	currentArch := runtime.GOARCH
	archIdx := 0
	for i, arch := range archs {
		if arch == currentArch {
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
	if m.versionBox.Items().Count() > 0 {
		// 尝试恢复上次选中的版本
		savedVersion := config.Config.Chromium.Version
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
	chromiumMap := config.DesignerConfig.Chromium
	versions := make([]string, 0, len(chromiumMap))
	for ver := range chromiumMap {
		versions = append(versions, ver)
	}
	sort.Slice(versions, func(i, j int) bool {
		return compareVersion(versions[i], versions[j]) < 0
	})
	return versions
}

// buildInstalledSet 返回已安装版本的集合, key 为 os_arch_version
func (m *TChromiumDirForm) buildInstalledSet() map[string]bool {
	manifest := config.Config.Chromium.LoadCEFManifest()
	installed := make(map[string]bool)
	for oav := range manifest {
		if config.Config.Chromium.IsCEFInstalled(oav) {
			installed[oav] = true
		}
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

// compareVersion 比较两个版本号, 如 "109.1.18" vs "127.3.5"
func compareVersion(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}
	for i := 0; i < maxLen; i++ {
		var aNum, bNum int
		if i < len(aParts) {
			aNum, _ = strconv.Atoi(aParts[i])
		}
		if i < len(bParts) {
			bNum, _ = strconv.Atoi(bParts[i])
		}
		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
	}
	return 0
}

// cefOSArch 返回 CEF 下载链接中的 {osarch} 部分
func cefOSArch(osName, arch string) string {
	if osMap, ok := cefOSArchMap[osName]; ok {
		if cefArch, ok := osMap[arch]; ok {
			return cefArch
		}
	}
	return "windows64"
}

// buildDownloadURL 构建 CEF 下载链接
func buildDownloadURL(version, osName, arch string) string {
	urlTemplate := config.DesignerConfig.Chromium.Get(version)
	if urlTemplate == "" {
		return ""
	}
	url := strings.ReplaceAll(urlTemplate, "{version}", version)
	url = strings.ReplaceAll(url, "{osarch}", cefOSArch(osName, arch))
	return url
}

// cefArchiveFileName 从下载 URL 中提取压缩包文件名 (URL 解码)
func cefArchiveFileName(version, osName, arch string) string {
	dlURL := buildDownloadURL(version, osName, arch)
	if dlURL == "" {
		return fmt.Sprintf("cef_binary_%s_%s_minimal.tar.bz2", version, cefOSArch(osName, arch))
	}
	// 取 URL 最后一段作为文件名, 并解码 %2B 等编码字符
	idx := strings.LastIndex(dlURL, "/")
	if idx >= 0 {
		name := dlURL[idx+1:]
		if decoded, err := url.PathUnescape(name); err == nil {
			return decoded
		}
		return name
	}
	return fmt.Sprintf("cef_binary_%s_%s_minimal.tar.bz2", version, cefOSArch(osName, arch))
}

// ==================== 窗口事件 ====================

func (m *TChromiumDirForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	m.dlMu.Lock()
	defer m.dlMu.Unlock()
	if m.dlState == downloadRunning {
		if m.dlCancel != nil {
			m.dlCancel()
		}
	}
	if m.dlState == downloadExtracting {
		m.dlStop = true
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
		config.Config.Chromium.Version = oav
		// 仅 CEF 项目更新项目配置
		if bean.GProject != nil && bean.GProject.GUIRenderFramework == bean.GUIRenderFramework_CEF {
			bean.GProject.FrameworkVersion = oav
		}
		config.UpdateConfig()
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
	downloadURL := buildDownloadURL(version, osName, arch)
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
	config.Config.Chromium.Dir = absDir
	m.dlVersion = version

	// 锁定 UI
	m.setDownloadState(downloadRunning)
	m.progressBar.SetVisible(true)
	m.statusLabel.SetVisible(true)
	m.statusLabel.SetCaption(metadata.CEFFormStatusLabelCaptionCaptionFailedPreDown)
	m.defaultBtn.SetEnabled(false)
	m.dirBtn.SetEnabled(false)
	m.osBox.SetEnabled(false)
	m.archBox.SetEnabled(false)
	m.versionBox.SetEnabled(false)

	targetPath := filepath.Join(absDir, cefArchiveFileName(version, osName, arch))
	event.ConsoleWriteInfo("Start downloading CEF:", downloadURL)
	event.ConsoleWriteInfo("Target:", targetPath)

	go m.doDownload(downloadURL, targetPath, version)
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
	m.dlStop = true
	m.dlMu.Unlock()
	// 即时反馈: 按钮变为"停止中", 禁用防重复点击
	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.confirmBtn.SetText(metadata.GI18n.Dict("ChromiumDirFormConfirmBtn.TextStopping"))
		m.confirmBtn.SetEnabled(false)
	})
}

// resumeDownload 继续下载
func (m *TChromiumDirForm) resumeDownload() {
	version := m.selectedVersion()
	osName := m.selectedOS()
	arch := m.selectedArch()

	// 检测版本是否变更
	versionChanged := version != m.dlVersion

	absDir, _ := filepath.Abs(m.dirEdit.Text())
	targetPath := filepath.Join(absDir, cefArchiveFileName(version, osName, arch))

	if versionChanged {
		// 版本变更, 删除旧的不完整文件
		oldPath := filepath.Join(absDir, cefArchiveFileName(m.dlVersion, osName, arch))
		if tool.IsExist(oldPath) {
			os.Remove(oldPath)
			event.ConsoleWriteInfo("Removed old partial file:", oldPath)
		}
		m.dlVersion = version
		m.dlProgress = 0
		m.dlTotal = 0
	}

	downloadURL := buildDownloadURL(version, osName, arch)

	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.setDownloadState(downloadRunning)
		m.osBox.SetEnabled(false)
		m.archBox.SetEnabled(false)
		m.versionBox.SetEnabled(false)
		if versionChanged {
			m.statusLabel.SetCaption(metadata.GI18n.Dict("ChromiumDirFormStatusLabel.CaptionReDownload"))
		} else {
			m.statusLabel.SetCaption(metadata.GI18n.Dict("ChromiumDirFormStatusLabel.CaptionFailedResumeDown"))
		}
	})

	go m.doDownload(downloadURL, targetPath, version)
}

// doDownload 执行下载, 支持断点续传
func (m *TChromiumDirForm) doDownload(url, targetPath, version string) {
	ctx, cancel := context.WithCancel(context.Background())
	m.dlMu.Lock()
	m.dlCancel = cancel
	m.dlMu.Unlock()

	defer func() {
		m.dlMu.Lock()
		m.dlCancel = nil
		m.dlMu.Unlock()
	}()

	// 检查已下载的文件大小, 用于断点续传
	var existingSize int64
	if info, err := os.Stat(targetPath); err == nil {
		existingSize = info.Size()
	}

	// 先用 HEAD 检查远程文件大小, 判断是否已完整下载
	remoteSize := m.checkRemoteFileSize(ctx, url)
	if ctx.Err() != nil {
		return
	}

	if existingSize > 0 && remoteSize > 0 && existingSize >= remoteSize {
		// 文件已完整下载, 跳过下载直接完成
		event.ConsoleWriteInfo("CEF archive already downloaded:", targetPath)
		m.dlProgress = existingSize
		m.dlTotal = remoteSize
		m.onDownloadCompleted(version, targetPath)
		return
	}

	// 构建请求, 支持断点续传
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		m.onDownloadError(err.Error())
		return
	}

	if existingSize > 0 {
		// 设置 Range 头实现断点续传
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", existingSize))
		event.ConsoleWriteInfo("Resuming download from:", formatSize(existingSize))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		m.onDownloadError(err.Error())
		return
	}
	defer resp.Body.Close()

	// 判断服务器是否支持断点续传
	var startSize int64
	if resp.StatusCode == http.StatusPartialContent {
		// 206: 服务器支持 Range, 从断点继续
		startSize = existingSize
		m.dlTotal = remoteSize
	} else if resp.StatusCode == http.StatusOK {
		// 200: 服务器不支持 Range, 从头开始
		startSize = 0
		m.dlTotal = resp.ContentLength
		existingSize = 0
	} else {
		m.onDownloadError(fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status))
		return
	}

	m.dlProgress = startSize

	// 打开文件: 续传用追加, 新下载用创建
	var out *os.File
	if startSize > 0 {
		out, err = os.OpenFile(targetPath, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		out, err = os.Create(targetPath)
	}
	if err != nil {
		m.onDownloadError(err.Error())
		return
	}
	defer out.Close()

	// 带进度的读取
	buf := make([]byte, 32*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				m.onDownloadError(writeErr.Error())
				return
			}
			m.dlProgress += int64(n)
			m.updateProgress()
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			if ctx.Err() != nil {
				// 被取消, 进度已保存, 下次可续传
				return
			}
			m.onDownloadError(readErr.Error())
			return
		}
	}

	m.onDownloadCompleted(version, targetPath)
}

// checkRemoteFileSize 通过 HEAD 请求获取远程文件大小
func (m *TChromiumDirForm) checkRemoteFileSize(ctx context.Context, url string) int64 {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return -1
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return resp.ContentLength
	}
	return -1
}

// ==================== UI 更新 ====================

// updateProgress 更新下载进度条 (在下载协程中调用)
func (m *TChromiumDirForm) updateProgress() {
	m.updateProgressStatus("Downloading... %s / %s", formatSize(m.dlProgress), formatSize(m.dlTotal))
}

// onDownloadError 下载出错 (在下载协程中调用)
func (m *TChromiumDirForm) onDownloadError(msg string) {
	event.ConsoleWriteError("CEF download error:", msg)
	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.resetToIdle()
		m.statusLabel.SetCaption(metadata.GI18n.Dict("ChromiumDirFormStatusLabel.CaptionFailDownload") +
			": " + msg)
	})
}

// onDownloadCompleted 下载完成, 开始解压
func (m *TChromiumDirForm) onDownloadCompleted(version, targetPath string) {
	event.ConsoleWriteInfo("CEF download completed:", targetPath)
	m.dlStop = false

	oav := m.osArchVersion(version)

	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.setDownloadState(downloadExtracting)
		m.statusLabel.SetCaption(metadata.GI18n.Dict("ChromiumDirFormStatusLabel.CaptionSuccessDownloadUnZip"))
		m.progressBar.SetPosition(0)
	})

	destDir := filepath.Join(config.Config.Chromium.Dir, oav)

	files, err := m.extractTarBz2(targetPath, destDir)
	if err != nil {
		if err == errExtractStopped {
			event.ConsoleWriteInfo("CEF extract stopped by user")
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.resetToIdle()
			})
		} else {
			event.ConsoleWriteError("CEF extract error:", err.Error())
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.resetToIdle()
				m.statusLabel.SetCaption(metadata.GI18n.Dict("ChromiumDirFormStatusLabel.CaptionFailUnZip") +
					": " + err.Error())
			})
		}
		return
	}

	// 保存安装清单
	if err = config.Config.Chromium.SaveCEFManifest(oav, files); err != nil {
		event.ConsoleWriteError("Failed to save CEF manifest:", err.Error())
	}

	// 通过清单校验安装完整性
	if !config.Config.Chromium.IsCEFInstalled(oav) {
		event.ConsoleWriteError("CEF installation verification failed")
		lcl.RunOnMainThreadAsync(func(id uint32) {
			m.resetToIdle()
			m.statusLabel.SetCaption(metadata.GI18n.Dict("ChromiumDirFormStatusLabel.CaptionFailVerify"))
		})
		return
	}

	event.ConsoleWriteInfo("CEF installed to:", destDir, "files:", fmt.Sprintf("%d", len(files)))
	m.Version = oav
	// 记录当前 CEF 版本到全局配置
	config.Config.Chromium.Version = oav
	// 仅 CEF 项目更新项目配置
	if bean.GProject != nil && bean.GProject.GUIRenderFramework == bean.GUIRenderFramework_CEF {
		bean.GProject.FrameworkVersion = oav
	}
	config.UpdateConfig()

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

// onCompleteBtnClick 安装完成后的确认按钮点击
func (m *TChromiumDirForm) onCompleteBtnClick(sender lcl.IObject) {
	m.Confirmed = true
	m.Close()
}

// ==================== 解压 ====================

// stopReader 包装 io.Reader, 在每次读取时检查停止信号
type stopReader struct {
	r    io.Reader
	stop *bool
}

func (s *stopReader) Read(p []byte) (int, error) {
	if *s.stop {
		return 0, errExtractStopped
	}
	if len(p) > 32*1024 {
		p = p[:32*1024]
	}
	n, err := s.r.Read(p)
	if *s.stop {
		return n, errExtractStopped
	}
	return n, err
}

// extractTarBz2 解压 .tar.bz2 到指定目录
// 提取顶层 Release/ 和 Resources/ 目录内的文件到 destDir, 去掉这两个前缀
// macOS Resources/ 在 Release/ 内部, 自然不会被顶层 Resources/ 匹配
// 返回已记录的文件信息列表(排除 cefExcludeFiles)
func (m *TChromiumDirForm) extractTarBz2(archivePath, destDir string) ([]config.CEFFileInfo, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sf := &stopReader{r: f, stop: &m.dlStop}
	bz2Reader := bzip2.NewReader(sf)
	tarReader := tar.NewReader(bz2Reader)

	// 第一遍扫描: 确定根目录前缀 + 统计需要解压的文件数
	var rootDir string // 如 "cef_binary_xxx/"
	var totalFiles int64
	for {
		if m.dlStop {
			return nil, errExtractStopped
		}
		header, err := tarReader.Next()
		if err != nil {
			if m.dlStop {
				return nil, errExtractStopped
			}
			if err == io.EOF {
				break
			}
			return nil, err
		}
		name := header.Name
		// 从第一个包含 /Release/ 的路径提取根目录
		if rootDir == "" {
			if idx := strings.Index(name, "/Release/"); idx >= 0 {
				rootDir = name[:idx+1] // 包含尾部 /
			}
		}
		if rootDir != "" {
			rel := extractRelPath(name, rootDir)
			if rel != "" && header.Typeflag == tar.TypeReg {
				if !config.IsCEFExcludeFile(filepath.Base(rel)) {
					totalFiles++
				}
			}
		}
	}

	if rootDir == "" {
		return nil, fmt.Errorf("release directory not found in archive")
	}

	event.ConsoleWriteInfo("CEF archive root:", rootDir, "files:", fmt.Sprintf("%d", totalFiles))

	// 重新打开进行实际解压
	f.Close()
	f, err = os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sf = &stopReader{r: f, stop: &m.dlStop}
	bz2Reader = bzip2.NewReader(sf)
	tarReader = tar.NewReader(bz2Reader)

	if err = os.MkdirAll(destDir, os.ModePerm); err != nil {
		return nil, err
	}

	m.dlProgress = 0
	m.dlTotal = totalFiles

	var files []config.CEFFileInfo
	copyBuf := make([]byte, 32*1024)

	for {
		if m.dlStop {
			return nil, errExtractStopped
		}

		header, err := tarReader.Next()
		if err != nil {
			if m.dlStop {
				return nil, errExtractStopped
			}
			if err == io.EOF {
				break
			}
			return nil, err
		}

		name := header.Name
		rel := extractRelPath(name, rootDir)
		if rel == "" {
			continue
		}

		target := filepath.Join(destDir, rel)

		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(target, os.ModePerm); err != nil {
				return nil, err
			}
		case tar.TypeReg:
			if err = os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
				return nil, err
			}
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return nil, err
			}
			writeErr := m.copyWithStop(outFile, tarReader, copyBuf)
			outFile.Close()
			if writeErr != nil {
				return nil, writeErr
			}
			// 记录文件信息(排除 cefExcludeFiles)
			fileName := filepath.Base(rel)
			if !config.IsCEFExcludeFile(fileName) {
				files = append(files, config.CEFFileInfo{Name: rel, Size: header.Size})
			}

			m.dlProgress++
			m.updateExtractProgress()
		case tar.TypeSymlink:
			if err = os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
				return nil, err
			}
			os.Remove(target)
			if err = os.Symlink(header.Linkname, target); err != nil {
				event.ConsoleWriteInfo("Symlink skipped:", target, err.Error())
			}
		}
	}
	return files, nil
}

// extractRelPath 从 tar 路径中提取 Release/ 或 Resources/ 内的相对路径
// rootDir 为根目录前缀, 如 "cef_binary_xxx/"
// 匹配: rootDir + "Release/" + rel 或 rootDir + "Resources/" + rel
func extractRelPath(name, rootDir string) string {
	if strings.HasPrefix(name, rootDir+"Release/") {
		return strings.TrimPrefix(name, rootDir+"Release/")
	}
	if strings.HasPrefix(name, rootDir+"Resources/") {
		return strings.TrimPrefix(name, rootDir+"Resources/")
	}
	return ""
}

func (m *TChromiumDirForm) copyWithStop(dst io.Writer, src io.Reader, buf []byte) error {
	for {
		if m.dlStop {
			return errExtractStopped
		}
		nr, readErr := src.Read(buf)
		if nr > 0 {
			nw, writeErr := dst.Write(buf[:nr])
			if writeErr != nil {
				return writeErr
			}
			if nw != nr {
				return io.ErrShortWrite
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				return nil
			}
			if m.dlStop {
				return errExtractStopped
			}
			return readErr
		}
	}
}

// updateExtractProgress 更新解压进度
func (m *TChromiumDirForm) updateExtractProgress() {
	m.updateProgressStatus("Extracting... %d / %d", m.dlProgress, m.dlTotal)
}

// updateProgressStatus 通用进度更新 (在下载/解压协程中调用)
func (m *TChromiumDirForm) updateProgressStatus(format string, args ...any) {
	progress := m.dlProgress
	total := m.dlTotal
	statusText := fmt.Sprintf(format, args...)
	lcl.RunOnMainThreadAsync(func(id uint32) {
		if total > 0 {
			percent := int32(progress * 100 / total)
			m.progressBar.SetPosition(percent)
			m.statusLabel.SetCaption(fmt.Sprintf("%s (%d%%)", statusText, percent))
		} else {
			m.statusLabel.SetCaption(statusText)
		}
	})
}

// formatSize 格式化文件大小
func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// ==================== CEF 版本辅助函数 ====================

// GetInstalledCEFVersions 返回当前系统架构下已安装的 CEF 版本列表 (os_arch_version 格式)
func GetInstalledCEFVersions() []string {
	manifest := config.Config.Chromium.LoadCEFManifest()
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH
	var versions []string
	for oav := range manifest {
		if config.Config.Chromium.IsCEFInstalled(oav) {
			// 检查是否匹配当前系统架构
			prefix := fmt.Sprintf("%s_%s_", currentOS, currentArch)
			if strings.HasPrefix(oav, prefix) {
				versions = append(versions, oav)
			}
		}
	}
	sort.Slice(versions, func(i, j int) bool {
		return compareVersion(extractVersionFromOAV(versions[i]), extractVersionFromOAV(versions[j])) > 0
	})
	return versions
}

// extractVersionFromOAV 从 os_arch_version 格式中提取版本号
func extractVersionFromOAV(oav string) string {
	parts := strings.SplitN(oav, "_", 3)
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}
