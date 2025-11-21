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

package helperform

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
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

var (
	bgColor     = colors.RGBToColor(56, 57, 60)
	bgTextColor = colors.ClGray
)

func NewCreateProjectForm() *TCreateProjectForm {
	designerForm := &TCreateProjectForm{}
	lcl.Application.NewForm(designerForm)
	return designerForm
}

type TCreateProjectForm struct {
	lcl.TEngForm
	oldWndPrc       uintptr
	goVersionOK     bool
	box             lcl.IPanel
	baseGroupBox    lcl.IGroupBox
	projNameText    lcl.ILabel
	projNameEdit    lcl.IEdit
	projPathText    lcl.ILabel
	projPathEdit    lcl.IEdit
	projPathBtn     *wg.TButton
	projPathDir     lcl.ISelectDirectoryDialog
	projTempText    lcl.ILabel
	projTempBox     lcl.IComboBox
	goVersionText   lcl.ILabel
	goVersionStatus *wg.TButton
	modGroupBox     lcl.IGroupBox
	modText         lcl.ILabel
	modBox          lcl.IComboBox
	cancelBtn       *wg.TButton
	createBtn       *wg.TButton
}

func (m *TCreateProjectForm) FormCreate(sender lcl.IObject) {
	logs.Info("TCreateProjectForm FormCreate")
	//defaultSize := int32(510)
	m.SetCaption("新建项目")
	m.SetWidth(555)
	m.SetHeight(555)
	constr := m.Constraints()
	constr.SetMaxWidth(555)
	constr.SetMaxHeight(555)
	constr.SetMinWidth(555)
	constr.SetMinHeight(555)
	//m.SetColor(bgColor)
	m.SetBorderIcons(types.NewSet())
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
	m.modGroupBox.SetCaption("模块依赖")
	m.modGroupBox.SetFont(fontLabel)
	m.modGroupBox.SetParent(m.box)

	m.baseGroupBox = lcl.NewGroupBox(m)
	m.baseGroupBox.SetAlign(types.AlTop)
	m.baseGroupBox.SetHeight(245)
	m.baseGroupBox.BorderSpacing().SetAround(6)
	m.baseGroupBox.SetCaption("基础信息")
	m.baseGroupBox.SetFont(fontLabel)
	m.baseGroupBox.SetParent(m.box)

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
		m.projPathEdit.SetWidth(290)
		m.projPathEdit.SetFont(fontText)
		m.projPathEdit.SetReadOnly(true)
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

		m.projPathDir = lcl.NewSelectDirectoryDialog(m)
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
		m.modText.SetCaption("依赖来源")
		m.modText.SetFont(fontLabel)
		m.modText.SetParent(m.modGroupBox)

		m.modBox = lcl.NewComboBox(m)
		m.modBox.SetBounds(120, 15, textWidth, 36)
		m.modBox.SetFont(fontText)
		m.modBox.SetReadOnly(true)
		m.modBox.SetStyle(types.CsDropDownList)
		m.modBox.SetBorderStyle(types.BsSingle)
		m.modBox.Items().Add("远程仓库 (从远程拉取)")
		m.modBox.Items().Add("本地路径 (离线/手动指定)")
		m.modBox.SetItemIndex(0)
		m.modBox.SetParent(m.modGroupBox)
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
			m.goVersionOK = compareVersions(version, "1.20") == 1
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.goVersionStatus.SetText(buf.String())
				if m.goVersionOK {
					m.goVersionStatus.SetColor(colors.ClGreen)
					m.goVersionStatus.SetIconFavoriteFormBytes(resources.Images("button/laugh.png"))
				} else {
					m.goVersionStatus.SetColor(colors.ClRed)
					m.goVersionStatus.SetIconFavoriteFormBytes(resources.Images("button/weep.png"))
				}
				m.goVersionStatus.ForcePaint(func() {
					m.goVersionStatus.Invalidate()
				})
			})
		}
		result = true
	}
	cmd.Command("go", "version")
}

func (m *TCreateProjectForm) projPathClick(sender lcl.IObject) {
	m.projPathDir.SetTitle("选择目录")
	if m.projPathDir.Execute() {
		dir := m.projPathDir.FileName()
		m.projPathEdit.SetText(dir)
	}
}

func (m *TCreateProjectForm) closeClick(sender lcl.IObject) {
	m.Close()
}

func (m *TCreateProjectForm) createClick(sender lcl.IObject) {
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
