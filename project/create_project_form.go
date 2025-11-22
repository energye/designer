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
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/designer/resources/frameworks"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/tool/command"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/font"
	"github.com/energye/widget/wg"
	"strconv"
	"strings"
	"time"
)

// 在设计器中创建项目
// 功能
// 主要: 在指定目录创建一个 energy 应用, 应用有默认模板
// 1. 应用名
// 2. 应用目录
// 3. 所需依赖模块(lcl, cef, webview), 从网络下载, 或设计器内绑定(✔️)
// 4. 模块模式： go.mod (✔️), go.work

/*
# 设计器安装目录（如 Windows：C:\EnergyDesigner，Linux：/opt/EnergyDesigner）
EnergyDesigner/
├── designer.exe  # 设计器主程序
└── frameworks/  # 内置框架根目录
	└── energy/  # 框架模块（必须是标准 Go 模块）
		├── go.mod  # 框架自身的 go.mod（如 module github.com/your-org/energy）
		├── go.sum
		├── core/  # 框架核心代码
		└── v1.2.3/  # （可选）多版本支持，每个版本一个独立模块
			├── go.mod
			└── core/

module my-app  // 项目自身的模块路径

go 1.20  // 项目使用的 Go 版本

// 声明框架依赖（版本需与内置框架的 go.mod 一致）
require github.com/your-org/energy v1.2.3

// 关键：将框架的网络模块路径，替换为设计器内置的本地路径
replace github.com/your-org/energy v1.2.3 => C:/EnergyDesigner/internal/frameworks/energy/v1.2.3

// 若框架无多版本，直接指向根目录：
// replace github.com/your-org/energy => C:/EnergyDesigner/internal/frameworks/energy

// 设计器新建 mod 模式项目时，配置内置框架
func createModProject(projectPath, frameworkBuiltInPath string) error {
    // 1. 创建项目目录
    if err := os.MkdirAll(projectPath, 0755); err != nil {
        return err
    }
    // 2. 初始化 go.mod
    cmd := exec.Command("go", "mod", "init", "my-app")
    cmd.Dir = projectPath
    if err := cmd.Run(); err != nil {
        return err
    }
    // 3. 写入 go.mod（添加 require 和 replace）
    modContent := fmt.Sprintf(`module my-app

go 1.21

require github.com/your-org/energy v1.2.3

replace github.com/your-org/energy v1.2.3 => %s`, frameworkBuiltInPath)
    return os.WriteFile(filepath.Join(projectPath, "go.mod"), []byte(modContent), 0644)
}
*/

const (
	formWidth    = int32(555)
	formHeight   = int32(555)
	minGoVersion = "1.20"
)

var (
	bgColor          = colors.RGBToColor(56, 57, 60)
	bgTextColor      = colors.ClGray
	modSelectOptions = []string{"本地路径 (内置框架)", "远程仓库 (远程拉取)"}
)

func NewCreateProjectForm() *TCreateProjectForm {
	designerForm := &TCreateProjectForm{}
	lcl.Application.NewForm(designerForm)
	return designerForm
}

type TCreateProjectForm struct {
	lcl.TEngForm
	oldWndPrc   uintptr
	goVersionOK bool
	box         lcl.IPanel
	selectDir   lcl.ISelectDirectoryDialog
	// 基础信息部分
	baseGroupBox    lcl.IGroupBox
	baseErrorLabel  lcl.ILabel
	projNameText    lcl.ILabel
	projNameEdit    lcl.IEdit
	projPathText    lcl.ILabel
	projPathEdit    lcl.IEdit
	projPathBtn     *wg.TButton
	projTempText    lcl.ILabel
	projTempBox     lcl.IComboBox
	goVersionText   lcl.ILabel
	goVersionStatus *wg.TButton
	// 模块部分
	modGroupBox   lcl.IGroupBox
	modErrorLabel lcl.ILabel
	modText       lcl.ILabel
	modBox        lcl.IComboBox
	modLocalBox   lcl.IPanel
	modRemoteBox  lcl.IPanel

	modLocalDirText    lcl.ILabel
	modLocalDirEdit    lcl.IEdit
	modLocalDirBtn     *wg.TButton
	modLCLCheckBox     lcl.ICheckBox
	modCEFCheckBox     lcl.ICheckBox
	modWebviewCheckBox lcl.ICheckBox

	// 操作按钮
	cancelBtn *wg.TButton
	createBtn *wg.TButton
}

func (m *TCreateProjectForm) FormCreate(sender lcl.IObject) {
	logs.Info("TCreateProjectForm FormCreate")
	m.SetCaption("新建项目")
	m.SetWidth(formWidth)
	m.SetHeight(formHeight)
	constr := m.Constraints()
	constr.SetMaxWidth(formWidth)
	constr.SetMaxHeight(formHeight)
	constr.SetMinWidth(formWidth)
	constr.SetMinHeight(formHeight)
	m.SetFormStyle(types.FsStayOnTop)
	m.SetShowInTaskBar(types.StNever)
	//m.SetColor(bgColor)
	m.SetBorderIcons(types.NewSet(types.BiSystemMenu))
	m.WorkAreaCenter()
	m.box = lcl.NewPanel(m)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetAlign(types.AlClient)
	m.box.SetParent(m)
	m.initComponents()
	m.SetOnShow(m.onShow)
	//m._HookWndProcMessage()
}

func (m *TCreateProjectForm) initComponents() {
	fontLabel := lcl.NewFont()
	fontLabel.SetName("微软雅黑 Light")
	fontLabel.SetStyle(types.NewSet(types.FsBold))
	fontLabel.SetSize(12)
	fontLabel.SetCharSet(font.CHINESEBIG5_CHARSET)
	//fontLabel.SetColor(colors.ClGreen)
	fontText := lcl.NewFont()
	fontText.SetName("微软雅黑 Light")
	fontText.SetSize(12)

	left := int32(35)
	textWidth := int32(355)

	m.modGroupBox = lcl.NewGroupBox(m)
	m.modGroupBox.SetAlign(types.AlTop)
	//m.modGroupBox.SetTop(m.baseGroupBox.Height() + 25)
	m.modGroupBox.SetHeight(240)
	m.modGroupBox.BorderSpacing().SetAround(6)
	m.modGroupBox.SetCaption("ENERGY框架-模块依赖")
	m.modGroupBox.SetFont(fontLabel)
	m.modGroupBox.SetParent(m.box)

	m.baseGroupBox = lcl.NewGroupBox(m)
	m.baseGroupBox.SetAlign(types.AlTop)
	m.baseGroupBox.SetHeight(245)
	m.baseGroupBox.BorderSpacing().SetAround(6)
	m.baseGroupBox.SetCaption("新建项目-基础信息")
	m.baseGroupBox.SetFont(fontLabel)
	m.baseGroupBox.SetParent(m.box)

	m.selectDir = lcl.NewSelectDirectoryDialog(m)

	{
		errorFont := lcl.NewFont()
		errorFont.SetColor(colors.ClRed)
		errorFont.SetName("微软雅黑 Light")
		errorFont.SetCharSet(font.CHINESEBIG5_CHARSET)
		errorFont.SetSize(8)
		m.baseErrorLabel = lcl.NewLabel(m)
		m.baseErrorLabel.SetFont(errorFont)
		m.baseErrorLabel.SetVisible(false)
		m.baseErrorLabel.SetParent(m.baseGroupBox)
		m.modErrorLabel = lcl.NewLabel(m)
		m.modErrorLabel.SetFont(errorFont)
		m.modErrorLabel.SetVisible(false)
		m.modErrorLabel.SetParent(m.modGroupBox)
	}
	{
		m.projNameText = lcl.NewLabel(m)
		m.projNameText.SetLeft(left)
		m.projNameText.SetTop(20)
		m.projNameText.SetWidth(80)
		m.projNameText.SetCaption("项目名称")
		m.projNameText.SetFont(fontLabel)
		m.projNameText.SetParent(m.baseGroupBox)

		m.projNameEdit = lcl.NewEdit(m)
		m.projNameEdit.SetLeft(120)
		m.projNameEdit.SetTop(15)
		m.projNameEdit.SetWidth(textWidth)
		m.projNameEdit.SetFont(fontText)
		m.projNameEdit.SetParent(m.baseGroupBox)
		m.projNameEdit.SetTextHint("新建的项目名称")
	}
	{
		m.projPathText = lcl.NewLabel(m)
		m.projPathText.SetLeft(left)
		m.projPathText.SetTop(70)
		m.projPathText.SetWidth(80)
		m.projPathText.SetCaption("项目路径")
		m.projPathText.SetFont(fontLabel)
		m.projPathText.SetParent(m.baseGroupBox)

		m.projPathEdit = lcl.NewEdit(m)
		m.projPathEdit.SetLeft(120)
		m.projPathEdit.SetTop(65)
		m.projPathEdit.SetWidth(textWidth - 65) // 是目录选择按钮的 宽度+Left(5)
		m.projPathEdit.SetFont(fontText)
		//m.projPathEdit.SetReadOnly(true)
		m.projPathEdit.SetTextHint("项目的存放目录")
		m.projPathEdit.SetParent(m.baseGroupBox)

		m.projPathBtn = wg.NewButton(m)
		m.projPathBtn.SetIconFormBytes(resources.Images("actions/add.png"))
		m.projPathBtn.SetRadius(3)
		cusRect := types.TRect{Left: m.projPathEdit.Left() + m.projPathEdit.Width() + 5, Top: 65}
		cusRect.SetWidth(60)
		cusRect.SetHeight(30)
		m.projPathBtn.SetBoundsRect(cusRect)
		m.projPathBtn.SetParent(m.baseGroupBox)
		m.projPathBtn.SetOnClick(m.projPathClick)
	}
	{
		m.projTempText = lcl.NewLabel(m)
		m.projTempText.SetLeft(left)
		m.projTempText.SetTop(120)
		m.projTempText.SetWidth(80)
		m.projTempText.SetCaption("项目模板")
		m.projTempText.SetFont(fontLabel)
		m.projTempText.SetParent(m.baseGroupBox)

		m.projTempBox = lcl.NewComboBox(m)
		m.projTempBox.SetBounds(120, 115, textWidth, 36)
		m.projTempBox.SetFont(fontText)
		m.projTempBox.SetReadOnly(true)
		m.projTempBox.SetStyle(types.CsDropDownList)
		m.projTempBox.SetBorderStyle(types.BsSingle)
		m.projTempBox.Items().Add("默认预设模板")
		m.projTempBox.SetItemIndex(0)
		m.projTempBox.SetParent(m.baseGroupBox)
	}
	{
		m.goVersionText = lcl.NewLabel(m)
		m.goVersionText.SetLeft(left)
		m.goVersionText.SetTop(170)
		m.goVersionText.SetWidth(80)
		m.goVersionText.SetCaption(" Go 版本")
		m.goVersionText.SetFont(fontLabel)
		m.goVersionText.SetParent(m.baseGroupBox)

		m.goVersionStatus = wg.NewButton(m)
		m.goVersionStatus.SetText("检测本地")
		m.goVersionStatus.SetFont(fontText)
		m.goVersionStatus.Font().SetColor(colors.ClWhite)
		m.goVersionStatus.Font().SetStyle(types.NewSet(types.FsBold))
		m.goVersionStatus.SetRadius(3)
		goVersionRect := types.TRect{Left: m.goVersionText.Left() + m.goVersionText.Width() + 5, Top: 165}
		goVersionRect.SetWidth(textWidth)
		goVersionRect.SetHeight(30)
		m.goVersionStatus.SetBoundsRect(goVersionRect)
		m.goVersionStatus.SetColor(colors.ClGray)
		m.goVersionStatus.SetParent(m.baseGroupBox)
	}
	{
		m.modText = lcl.NewLabel(m)
		m.modText.SetLeft(left)
		m.modText.SetTop(20)
		m.modText.SetWidth(80)
		m.modText.SetCaption("模块来源")
		m.modText.SetFont(fontLabel)
		m.modText.SetParent(m.modGroupBox)

		m.modBox = lcl.NewComboBox(m)
		m.modBox.SetBounds(120, 15, textWidth, 36)
		m.modBox.SetFont(fontText)
		m.modBox.SetReadOnly(true)
		m.modBox.SetStyle(types.CsDropDownList)
		m.modBox.SetBorderStyle(types.BsSingle)
		for _, option := range modSelectOptions {
			m.modBox.Items().Add(option)
		}
		m.modBox.SetItemIndex(0)
		m.modBox.SetParent(m.modGroupBox)
		m.modBox.SetOnChange(m.modBoxChange)

		m.modLocalBox = lcl.NewPanel(m)
		m.modLocalBox.SetBevelOuter(types.BvNone)
		m.modLocalBox.SetAlign(types.AlBottom)
		m.modLocalBox.SetHeight(150)
		m.modLocalBox.SetVisible(true)
		m.modLocalBox.SetParent(m.modGroupBox)

		m.modRemoteBox = lcl.NewPanel(m)
		m.modRemoteBox.SetBevelOuter(types.BvNone)
		m.modRemoteBox.SetAlign(types.AlBottom)
		m.modRemoteBox.SetHeight(150)
		m.modRemoteBox.SetColor(colors.ClGray)
		m.modRemoteBox.SetVisible(false)
		m.modRemoteBox.SetParent(m.modGroupBox)

		m.modLocalDirText = lcl.NewLabel(m)
		m.modLocalDirText.SetLeft(left)
		m.modLocalDirText.SetTop(5)
		m.modLocalDirText.SetWidth(80)
		m.modLocalDirText.SetCaption("框架目录")
		m.modLocalDirText.SetFont(fontLabel)
		m.modLocalDirText.SetParent(m.modLocalBox)

		m.modLocalDirEdit = lcl.NewEdit(m)
		m.modLocalDirEdit.SetLeft(120)
		m.modLocalDirEdit.SetTop(0)
		m.modLocalDirEdit.SetWidth(textWidth - 65) // 是目录选择按钮的 宽度+Left(5)
		m.modLocalDirEdit.SetFont(fontText)
		m.modLocalDirEdit.SetReadOnly(true)
		m.modLocalDirEdit.SetText(frameworks.Path)
		m.modLocalDirEdit.SetParent(m.modLocalBox)

		m.modLocalDirBtn = wg.NewButton(m)
		m.modLocalDirBtn.SetIconFormBytes(resources.Images("actions/add.png"))
		m.modLocalDirBtn.SetRadius(3)
		modLocalDirRect := types.TRect{Left: m.modLocalDirEdit.Left() + m.modLocalDirEdit.Width() + 5, Top: 0}
		modLocalDirRect.SetWidth(60)
		modLocalDirRect.SetHeight(30)
		m.modLocalDirBtn.SetBoundsRect(modLocalDirRect)
		m.modLocalDirBtn.SetParent(m.modLocalBox)
		m.modLocalDirBtn.SetOnClick(m.modLocalDirBtnClick)

		m.modLCLCheckBox = lcl.NewCheckBox(m)
		m.modLCLCheckBox.SetFont(fontText)
		m.modLCLCheckBox.SetLeft(left)
		m.modLCLCheckBox.SetTop(modLocalDirRect.Top + modLocalDirRect.Height() + 20)
		m.modLCLCheckBox.SetCaption("LCL (Native UI)")
		m.modLCLCheckBox.SetHint("Lazarus Component Library")
		m.modLCLCheckBox.SetShowHint(true)
		m.modLCLCheckBox.SetChecked(true)
		m.modLCLCheckBox.SetEnabled(frameworks.EnableLCL)
		m.modLCLCheckBox.SetParent(m.modLocalBox)

		m.modCEFCheckBox = lcl.NewCheckBox(m)
		m.modCEFCheckBox.SetFont(fontText)
		m.modCEFCheckBox.SetLeft(left + 150)
		m.modCEFCheckBox.SetTop(modLocalDirRect.Top + modLocalDirRect.Height() + 20)
		m.modCEFCheckBox.SetCaption("CEF (Web UI)")
		m.modCEFCheckBox.SetHint("Chromium Embedded Framework")
		m.modCEFCheckBox.SetShowHint(true)
		//m.modCEFCheckBox.SetChecked(true)
		m.modCEFCheckBox.SetEnabled(frameworks.EnableCEF)
		m.modCEFCheckBox.SetParent(m.modLocalBox)

		m.modWebviewCheckBox = lcl.NewCheckBox(m)
		m.modWebviewCheckBox.SetFont(fontText)
		m.modWebviewCheckBox.SetLeft(left + 285)
		m.modWebviewCheckBox.SetTop(modLocalDirRect.Top + modLocalDirRect.Height() + 20)
		m.modWebviewCheckBox.SetCaption("WebView (Web UI)")
		m.modWebviewCheckBox.SetHint("System Runtime Framework")
		m.modWebviewCheckBox.SetShowHint(true)
		//m.modWebviewCheckBox.SetChecked(true)
		m.modWebviewCheckBox.SetEnabled(frameworks.EnableWV)
		m.modWebviewCheckBox.SetParent(m.modLocalBox)
	}
	{
		m.cancelBtn = wg.NewButton(m)
		m.cancelBtn.SetText("关 闭")
		m.cancelBtn.SetFont(fontText)
		m.cancelBtn.Font().SetColor(colors.ClWhite)
		m.cancelBtn.SetRadius(3)
		cancelBtnRect := types.TRect{Left: 310, Top: 505}
		cancelBtnRect.SetWidth(100)
		cancelBtnRect.SetHeight(40)
		m.cancelBtn.SetBoundsRect(cancelBtnRect)
		m.cancelBtn.SetColor(colors.RGBToColor(255, 127, 127))
		m.cancelBtn.SetParent(m.box)
		m.cancelBtn.SetOnClick(m.closeClick)

		m.createBtn = wg.NewButton(m)
		m.createBtn.SetText("创 建")
		m.createBtn.SetFont(fontText)
		m.createBtn.Font().SetColor(colors.ClWhite)
		m.createBtn.SetRadius(3)
		createBtnRect := types.TRect{Left: 430, Top: 505}
		createBtnRect.SetWidth(100)
		createBtnRect.SetHeight(40)
		m.createBtn.SetBoundsRect(createBtnRect)
		m.createBtn.SetColor(colors.RGBToColor(46, 204, 113))
		m.createBtn.SetParent(m.box)
		m.createBtn.SetOnClick(m.createClick)
	}
}

// 窗口显示事件
func (m *TCreateProjectForm) onShow(sender lcl.IObject) {
	//width := int32(555)
	//height := int32(515)
	//m.SetWidth(width)
	//m.SetHeight(height)
	//m.SetBoundsRect(m.BoundsRect()) // trigger WM_NCCALCSIZE hook msg
	//m.WorkAreaCenter()

	go m.checkGoVersion()
}

func (m *TCreateProjectForm) checkGoVersion() {
	time.Sleep(time.Second / 2)
	result := false
	cmd := command.NewCMD()
	cmd.IsPrint = false
	cmd.Console = func(data string, level command.Level) {
		if !result {
			logs.Debug(level, data)
			parts := strings.Fields(data)
			buf := tool.Buffer{}
			version := ""
			for i, part := range parts {
				if tool.Equal(part, "go", "version") {
					continue
				}
				if i == 2 {
					version = part[2:]
				}
				buf.WriteString(part, " ")
			}
			// 支持的最低Go版本
			m.goVersionOK = compareVersions(version, minGoVersion) == 1
			lcl.RunOnMainThreadAsync(func(id uint32) {
				if m.goVersionOK {
					buf.WriteString(" 支持")
					m.goVersionStatus.SetColor(colors.ClGreen)
					m.goVersionStatus.SetIconFavoriteFormBytes(resources.Images("button/laugh.png"))
				} else {
					buf.WriteString(" 不支持")
					m.goVersionStatus.SetColor(colors.ClRed)
					m.goVersionStatus.SetIconFavoriteFormBytes(resources.Images("button/weep.png"))
				}
				m.goVersionStatus.SetText(buf.String())
				m.goVersionStatus.ForcePaint(func() {
					m.goVersionStatus.Invalidate()
				})
			})
		}
		result = true
	}
	cmd.Command("go", "version")
}

// 项目存放目录选择
func (m *TCreateProjectForm) modLocalDirBtnClick(sender lcl.IObject) {
	m.selectDir.SetTitle("框架安装目录")
	m.selectDir.SetInitialDir(m.modLocalDirEdit.Text())
	if m.selectDir.Execute() {
		dir := m.selectDir.FileName()
		m.modLocalDirEdit.SetText(dir)
	}
}

// 项目存放目录选择
func (m *TCreateProjectForm) projPathClick(sender lcl.IObject) {
	m.selectDir.SetTitle("新建项目")
	if m.projPathEdit.Text() == "" {
		//m.projPathDir.SetFileName(exec.Dir)
	} else {
		m.selectDir.SetInitialDir(m.projPathEdit.Text())
	}
	if m.selectDir.Execute() {
		dir := m.selectDir.FileName()
		m.projPathEdit.SetText(dir)
	}
}

// 关闭
func (m *TCreateProjectForm) closeClick(sender lcl.IObject) {
	m.Close()
}

// 创建
func (m *TCreateProjectForm) createClick(sender lcl.IObject) {
	if !m.validateInputs() {
		return
	}
	//frameworkDir := m.modLocalDirEdit.Text()
	//isLCL := m.modLCLCheckBox.Checked()
	//isCEF := m.modCEFCheckBox.Checked()
	//isWV := m.modWebviewCheckBox.Checked()
	//
	//projectName := m.projNameEdit.Text()
	//projectDir := m.projPathEdit.Text()
	//doRunCreate()
}

// 模块选择
func (m *TCreateProjectForm) modBoxChange(sender lcl.IObject) {
	if m.modBox.ItemIndex() == 1 {
		m.showError(m.modErrorLabel, m.modBox.BoundsRect(), "＊暂不支持从远程下载")
		return
	}
	m.baseErrorLabel.Hide()
	m.modErrorLabel.Hide()
}

// 统一错误显示方法
func (m *TCreateProjectForm) showError(label lcl.ILabel, br types.TRect, message string) {
	label.SetLeft(br.Left)
	label.SetTop(br.Top + br.Height() + 5)
	label.SetCaption(message)
	label.Show()
}

func (m *TCreateProjectForm) validateInputs() bool {
	if strings.TrimSpace(m.projNameEdit.Text()) == "" {
		m.showError(m.baseErrorLabel, m.projNameEdit.BoundsRect(), "＊项目名为空")
		return false
	}
	selectProjectPath := strings.TrimSpace(m.projPathEdit.Text())
	if selectProjectPath == "" {
		m.showError(m.baseErrorLabel, m.projPathEdit.BoundsRect(), "＊项目目录为空")
		return false
	}
	//if !tool.IsExist(selectProjectPath) {
	//	m.showError(m.baseErrorLabel, m.projPathEdit.BoundsRect(), "＊目录不存在")
	//	return false
	//}
	if m.modBox.ItemIndex() == 1 {
		m.showError(m.modErrorLabel, m.modBox.BoundsRect(), "＊暂不支持从远程下载")
		return false
	}

	m.baseErrorLabel.Hide()
	m.modErrorLabel.Hide()
	return true
}

// compareVersions 比较两个版本号字符串的大小
//
// v1 - 第一个版本号字符串，格式如 "1.2.3"
// v2 - 第二个版本号字符串，格式如 "1.2.3"
// int - 比较结果：1表示v1大于v2，-1表示v1小于v2，0表示两者相等
func compareVersions(v1, v2 string) int {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")
	maxLen := len(v1Parts)
	if len(v2Parts) > maxLen {
		maxLen = len(v2Parts)
	}
	for i := 0; i < maxLen; i++ {
		part1 := getVersionPart(v1Parts, i)
		part2 := getVersionPart(v2Parts, i)

		if part1 > part2 {
			return 1
		} else if part1 < part2 {
			return -1
		}
	}
	return 0
}

func getVersionPart(parts []string, index int) int {
	if index >= len(parts) {
		return 0
	}
	num, err := strconv.Atoi(parts[index])
	if err != nil {
		return 0
	}
	return num
}

//func (m *TCreateProjectForm) wndProc(hwnd types.HWND, message uint32, wParam, lParam uintptr) uintptr {
//	switch message {
//	case messages.WM_DPICHANGED:
//		if !lcl.Application.Scaled() {
//			newWindowSize := (*types.TRect)(unsafe.Pointer(lParam))
//			win.SetWindowPos(m.Handle(), uintptr(0),
//				newWindowSize.Left, newWindowSize.Top, newWindowSize.Right-newWindowSize.Left, newWindowSize.Bottom-newWindowSize.Top,
//				win.SWP_NOZORDER|win.SWP_NOACTIVATE)
//		}
//	}
//	switch message {
//	case messages.WM_ACTIVATE:
//		win.ExtendFrameIntoClientArea(m.Handle(), win.Margins{CxLeftWidth: 1, CxRightWidth: 1, CyTopHeight: 1, CyBottomHeight: 1})
//	case messages.WM_NCCALCSIZE:
//		if wParam != 0 {
//			isMaximize := uint32(win.GetWindowLong(m.Handle(), win.GWL_STYLE))&win.WS_MAXIMIZE != 0
//			if isMaximize {
//				rect := (*types.TRect)(unsafe.Pointer(lParam))
//				monitor := win.MonitorFromRect(rect, win.MONITOR_DEFAULTTONULL)
//				if monitor != 0 {
//					var monitorInfo types.TMonitorInfo
//					monitorInfo.CbSize = types.DWORD(unsafe.Sizeof(monitorInfo))
//					if win.GetMonitorInfo(monitor, &monitorInfo) {
//						*rect = monitorInfo.RcWork
//					}
//				}
//			}
//			return 0
//		}
//	}
//
//	return win.CallWindowProc(m.oldWndPrc, uintptr(hwnd), message, wParam, lParam)
//}
//
//func (m *TCreateProjectForm) _HookWndProcMessage() {
//	wndProcCallback := syscall.NewCallback(m.wndProc)
//	m.oldWndPrc = win.SetWindowLongPtr(m.Handle(), win.GWL_WNDPROC, wndProcCallback)
//}

//func (m *TCreateProjectForm) _RestoreWndProc() {
//	if m.oldWndPrc != 0 {
//		win.SetWindowLongPtr(m.Handle(), win.GWL_WNDPROC, m.oldWndPrc)
//		m.oldWndPrc = 0
//	}
//}
