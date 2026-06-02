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
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
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

	// 版本选择
	versionText lcl.ILabel
	versionBox  lcl.IComboBox

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
	m.SetCaption("CEF 框架设置")
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
		m.dirEdit.SetTop(nextTop(25))
		m.dirEdit.SetWidth(350)
		m.dirEdit.SetDoubleBuffered(true)
		m.dirEdit.SetTextHint(config.Config.Chromium.DefaultDir())
		m.dirEdit.SetLabelPosition(types.LpLeft)
		// 如果已有配置的目录, 显示出来
		if config.Config.Chromium.Dir != "" {
			m.dirEdit.SetText(config.Config.Chromium.Dir)
		}
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

	// 版本选择 — 所有可用版本, 已安装标记
	{
		m.versionText = lcl.NewLabel(m)
		m.versionText.SetLeft(15)
		m.versionText.SetTop(nextTop(35))
		m.versionText.SetCaption("CEF 版本:")
		m.versionText.SetParent(m)

		m.versionBox = lcl.NewComboBox(m)
		m.versionBox.SetName("ChromiumDirFormVersionBox")
		m.versionBox.SetLeft(80)
		m.versionBox.SetTop(m.versionText.Top() - 3)
		m.versionBox.SetWidth(200)
		m.versionBox.SetReadOnly(true)
		m.versionBox.SetStyle(types.CsDropDownList)
		m.versionBox.SetBorderStyle(types.BsSingle)
		// 填充版本, 已安装标记
		installed := m.installedVersionSet()
		versions := m.sortedVersions()
		for _, ver := range versions {
			if installed[ver] {
				m.versionBox.Items().Add(ver + " (已安装)")
			} else {
				m.versionBox.Items().Add(ver)
			}
		}
		if m.versionBox.Items().Count() > 0 {
			m.versionBox.SetItemIndex(0)
		}
		m.versionBox.SetOnChange(m.onVersionChange)
		m.versionBox.SetParent(m)
	}

	// 下载进度 (初始隐藏)
	{
		m.progressBar = lcl.NewProgressBar(m)
		m.progressBar.SetName("ChromiumDirFormProgressBar")
		m.progressBar.SetLeft(80)
		m.progressBar.SetTop(nextTop(35))
		m.progressBar.SetWidth(350)
		m.progressBar.SetHeight(20)
		m.progressBar.SetParent(m)
		m.progressBar.SetVisible(false)

		m.statusLabel = lcl.NewLabel(m)
		m.statusLabel.SetLeft(80)
		m.statusLabel.SetTop(nextTop(22))
		m.statusLabel.SetCaption("选择目录和版本后, 点击确定开始下载")
		m.statusLabel.Font().SetColor(colors.RGBToColor(128, 128, 128))
		m.statusLabel.Font().SetSize(8)
		m.statusLabel.SetParent(m)
	}

	// 操作按钮
	{
		defaultBtnRect := types.TRect{Left: 15, Top: nextTop(30)}
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

	// 初始按钮文字
	if m.isVersionInstalled() {
		m.confirmBtn.SetText("使 用")
	}
}

// ==================== 版本与URL ====================

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

// installedVersionSet 返回已安装版本的集合
func (m *TChromiumDirForm) installedVersionSet() map[string]bool {
	manifest := config.Config.Chromium.LoadCEFManifest()
	installed := make(map[string]bool)
	for ver := range manifest {
		if config.Config.Chromium.IsCEFInstalled(ver) {
			installed[ver] = true
		}
	}
	return installed
}

// selectedVersion 返回下拉框选中的版本号(去掉标记后缀)
func (m *TChromiumDirForm) selectedVersion() string {
	idx := m.versionBox.ItemIndex()
	if idx < 0 {
		return ""
	}
	text := m.versionBox.Items().Strings(idx)
	return strings.TrimSuffix(text, " (已安装)")
}

// isVersionInstalled 判断下拉框选中的版本是否已安装
func (m *TChromiumDirForm) isVersionInstalled() bool {
	idx := m.versionBox.ItemIndex()
	if idx < 0 {
		return false
	}
	text := m.versionBox.Items().Strings(idx)
	return strings.HasSuffix(text, " (已安装)")
}

// onVersionChange 版本切换时更新按钮文字
func (m *TChromiumDirForm) onVersionChange(sender lcl.IObject) {
	if m.isVersionInstalled() {
		m.confirmBtn.SetText("使 用")
	} else {
		m.confirmBtn.SetText("确 定")
	}
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
			fmt.Sscanf(aParts[i], "%d", &aNum)
		}
		if i < len(bParts) {
			fmt.Sscanf(bParts[i], "%d", &bNum)
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
func cefOSArch() string {
	switch runtime.GOOS {
	case "windows":
		if runtime.GOARCH == "386" {
			return "windows32"
		}
		return "windows64"
	case "linux":
		if runtime.GOARCH == "386" {
			return "linux32"
		}
		return "linux64"
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return "macosarm64"
		}
		return "macos64"
	default:
		return "windows64"
	}
}

// buildDownloadURL 构建 CEF 下载链接
func buildDownloadURL(version string) string {
	urlTemplate := config.DesignerConfig.Chromium.Get(version)
	if urlTemplate == "" {
		return ""
	}
	url := strings.ReplaceAll(urlTemplate, "{version}", version)
	url = strings.ReplaceAll(url, "{osarch}", cefOSArch())
	return url
}

// cefArchiveFileName 返回 CEF 压缩包文件名
func cefArchiveFileName(version string) string {
	return fmt.Sprintf("cef_binary_%s_%s_client.tar.bz2", version, cefOSArch())
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
		m.Version = m.selectedVersion()
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
		m.confirmBtn.SetText("确 定")
		m.confirmBtn.SetColor(blueBtnColor)
	case downloadRunning, downloadExtracting:
		m.confirmBtn.SetText("停 止")
		m.confirmBtn.SetColor(colors.RGBToColor(255, 127, 127))
	case downloadPaused:
		m.confirmBtn.SetText("继 续")
		m.confirmBtn.SetColor(blueBtnColor)
	}
}

// resetToIdle 恢复窗口到默认空闲状态
func (m *TChromiumDirForm) resetToIdle() {
	m.setDownloadState(downloadIdle)
	m.progressBar.SetVisible(false)
	m.statusLabel.SetVisible(true)
	m.statusLabel.SetCaption("选择目录和版本后, 点击确定开始下载")
	m.defaultBtn.SetEnabled(true)
	m.dirBtn.SetEnabled(true)
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
		m.statusLabel.SetCaption("目录路径无效")
		return
	}

	// 获取选中版本
	version := m.selectedVersion()
	if version == "" {
		event.ConsoleWriteError("Please select a CEF version")
		m.statusLabel.SetCaption("请选择 CEF 版本")
		return
	}
	downloadURL := buildDownloadURL(version)
	if downloadURL == "" {
		event.ConsoleWriteError("Download URL not found for version:", version)
		m.statusLabel.SetCaption("未找到该版本的下载地址")
		return
	}

	// 创建目录
	if !tool.IsExist(absDir) {
		if err = os.MkdirAll(absDir, os.ModePerm); err != nil {
			event.ConsoleWriteError("Failed to create directory:", err.Error())
			m.statusLabel.SetCaption("创建目录失败: " + err.Error())
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
	m.statusLabel.SetCaption("准备下载...")
	m.defaultBtn.SetEnabled(false)
	m.dirBtn.SetEnabled(false)
	m.versionBox.SetEnabled(false)

	targetPath := filepath.Join(absDir, cefArchiveFileName(version))
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
		m.statusLabel.SetCaption("下载已暂停, 可更换版本后继续")
		m.versionBox.SetEnabled(true)
	})
}

// stopExtract 停止解压, 恢复窗口到默认状态
func (m *TChromiumDirForm) stopExtract() {
	m.dlMu.Lock()
	m.dlStop = true
	m.dlMu.Unlock()
}

// resumeDownload 继续下载
func (m *TChromiumDirForm) resumeDownload() {
	version := m.selectedVersion()

	// 检测版本是否变更
	versionChanged := version != m.dlVersion

	absDir, _ := filepath.Abs(m.dirEdit.Text())
	targetPath := filepath.Join(absDir, cefArchiveFileName(version))

	if versionChanged {
		// 版本变更, 删除旧的不完整文件
		oldPath := filepath.Join(absDir, cefArchiveFileName(m.dlVersion))
		if tool.IsExist(oldPath) {
			os.Remove(oldPath)
			event.ConsoleWriteInfo("Removed old partial file:", oldPath)
		}
		m.dlVersion = version
		m.dlProgress = 0
		m.dlTotal = 0
	}

	downloadURL := buildDownloadURL(version)

	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.setDownloadState(downloadRunning)
		m.versionBox.SetEnabled(false)
		if versionChanged {
			m.statusLabel.SetCaption("版本已更换, 重新下载...")
		} else {
			m.statusLabel.SetCaption("继续下载...")
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

// updateProgress 更新进度条 (在下载协程中调用)
func (m *TChromiumDirForm) updateProgress() {
	progress := m.dlProgress
	total := m.dlTotal
	lcl.RunOnMainThreadAsync(func(id uint32) {
		if total > 0 {
			percent := int32(progress * 100 / total)
			m.progressBar.SetPosition(percent)
			m.statusLabel.SetCaption(fmt.Sprintf("下载中... %s / %s (%d%%)", formatSize(progress), formatSize(total), percent))
		} else {
			m.statusLabel.SetCaption(fmt.Sprintf("下载中... %s", formatSize(progress)))
		}
	})
}

// onDownloadError 下载出错 (在下载协程中调用)
func (m *TChromiumDirForm) onDownloadError(msg string) {
	event.ConsoleWriteError("CEF download error:", msg)
	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.resetToIdle()
		m.statusLabel.SetCaption("下载失败: " + msg)
	})
}

// onDownloadCompleted 下载完成, 开始解压
func (m *TChromiumDirForm) onDownloadCompleted(version, targetPath string) {
	event.ConsoleWriteInfo("CEF download completed:", targetPath)
	m.dlStop = false
	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.setDownloadState(downloadExtracting)
		m.statusLabel.SetCaption("下载完成, 正在解压...")
		m.progressBar.SetPosition(0)
	})

	destDir := filepath.Join(config.Config.Chromium.Dir, version)

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
				m.statusLabel.SetCaption("解压失败: " + err.Error())
			})
		}
		return
	}

	// 保存安装清单
	if err = config.Config.Chromium.SaveCEFManifest(version, files); err != nil {
		event.ConsoleWriteError("Failed to save CEF manifest:", err.Error())
	}

	// 通过清单校验安装完整性
	if !config.Config.Chromium.IsCEFInstalled(version) {
		event.ConsoleWriteError("CEF installation verification failed")
		lcl.RunOnMainThreadAsync(func(id uint32) {
			m.resetToIdle()
			m.statusLabel.SetCaption("安装校验失败")
		})
		return
	}

	event.ConsoleWriteInfo("CEF installed to:", destDir, "files:", fmt.Sprintf("%d", len(files)))
	m.Version = version
	config.UpdateConfig()

	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.setDownloadState(downloadCompleted)
		m.progressBar.SetPosition(100)
		m.statusLabel.SetCaption("安装完成!")
		m.confirmBtn.SetEnabled(true)
		m.confirmBtn.SetText("完 成")
		m.confirmBtn.SetColor(colors.RGBToColor(46, 204, 113))
		m.confirmBtn.SetOnClick(func(sender lcl.IObject) {
			m.Confirmed = true
			m.Close()
		})
	})
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
// 只解压 Release/ 目录下的文件, 去掉 cef_binary_xxx/Release/ 前缀
// 返回已记录的文件信息列表(排除 cefExcludeFiles)
func (m *TChromiumDirForm) extractTarBz2(archivePath, destDir string) ([]config.CEFFileInfo, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bz2Reader := bzip2.NewReader(f)
	tarReader := tar.NewReader(bz2Reader)

	// 第一遍扫描: 找到 Release/ 前缀 + 统计需要解压的文件数
	var releasePrefix string
	var totalFiles int64
	for {
		if m.dlStop {
			return nil, errExtractStopped
		}
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		name := header.Name
		if releasePrefix == "" {
			if idx := strings.Index(name, "/Release/"); idx >= 0 {
				releasePrefix = name[:idx+len("/Release/")]
			}
		}
		if releasePrefix != "" && strings.HasPrefix(name, releasePrefix) {
			relPath := strings.TrimPrefix(name, releasePrefix)
			if relPath != "" && header.Typeflag == tar.TypeReg {
				if !config.IsCEFExcludeFile(filepath.Base(relPath)) {
					totalFiles++
				}
			}
		}
	}

	if releasePrefix == "" {
		return nil, fmt.Errorf("Release directory not found in archive")
	}

	event.ConsoleWriteInfo("CEF archive release prefix:", releasePrefix, "files:", fmt.Sprintf("%d", totalFiles))

	// 重新打开进行实际解压
	f.Close()
	f, err = os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bz2Reader = bzip2.NewReader(f)
	tarReader = tar.NewReader(bz2Reader)

	if err = os.MkdirAll(destDir, os.ModePerm); err != nil {
		return nil, err
	}

	m.dlProgress = 0
	m.dlTotal = totalFiles

	var files []config.CEFFileInfo

	for {
		if m.dlStop {
			return nil, errExtractStopped
		}

		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		name := header.Name
		if !strings.HasPrefix(name, releasePrefix) {
			continue
		}
		relPath := strings.TrimPrefix(name, releasePrefix)
		if relPath == "" {
			continue
		}

		target := filepath.Join(destDir, relPath)

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
			sr := &stopReader{r: tarReader, stop: &m.dlStop}
			if _, err = io.Copy(outFile, sr); err != nil {
				outFile.Close()
				return nil, err
			}
			outFile.Close()

			// 记录文件信息(排除 cefExcludeFiles)
			fileName := filepath.Base(relPath)
			if !config.IsCEFExcludeFile(fileName) {
				files = append(files, config.CEFFileInfo{Name: relPath, Size: header.Size})
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

// updateExtractProgress 更新解压进度
func (m *TChromiumDirForm) updateExtractProgress() {
	current := m.dlProgress
	total := m.dlTotal
	lcl.RunOnMainThreadAsync(func(id uint32) {
		if total > 0 {
			percent := int32(current * 100 / total)
			m.progressBar.SetPosition(percent)
			m.statusLabel.SetCaption(fmt.Sprintf("解压中... %d / %d (%d%%)", current, total, percent))
		} else {
			m.statusLabel.SetCaption(fmt.Sprintf("解压中... %d 个文件", current))
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
