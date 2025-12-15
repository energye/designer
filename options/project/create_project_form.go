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
	"fmt"
	"github.com/energye/designer/designer"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/tool/command"
	"github.com/energye/lcl/tool/exec"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/font"
	"github.com/energye/widget/wg"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 在设计器中创建项目
// 功能
// 主要: 在指定目录创建一个 energy 应用, 应用有默认模板
// 1. 应用名
// 2. 应用目录
// 3. 所需依赖模块(lcl, cef, webview), 从网络下载, 或设计器内绑定(✔️)
// 4. 模块模式： go.mod (✔️), go.work

var (
	createProjectFormWidth  = int32(525)
	createProjectFormHeight = int32(355)
	minGoVersion            = "1.20"
)

var (
	bgColor          = colors.RGBToColor(56, 57, 60)
	bgTextColor      = colors.ClGray
	modSelectOptions = []string{"本地路径 (内置框架)", "远程仓库 (远程拉取)"}
)

// NewCreateProjectForm 创建一个新的项目创建表单实例
// 该函数初始化一个 TCreateProjectForm 结构体，并通过 lcl.Application.NewForm 方法将其注册为应用程序窗体
func NewCreateProjectForm() *TCreateProjectForm {
	newEngForm := lcl.NewEngForm(nil)
	newForm := &TCreateProjectForm{TEngForm: newEngForm.(*lcl.TEngForm)}
	//lcl.Application.NewForm(newForm) // 不使用原因：go debug 模式有问题
	newForm.FormCreate(newForm)
	newForm.SetOnCloseQuery(newForm.OnCloseQuery)
	newForm.SetOnClose(newForm.OnClose)
	return newForm
}

type TCreateProjectForm struct {
	*lcl.TEngForm
	closing     bool
	goVersionOK bool
	one         sync.Once
	box         lcl.IPanel
	selectDir   lcl.ISelectDirectoryDialog

	// 基础信息部分
	baseGroupBox    lcl.ILabel
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
	modGroupBox   lcl.ILabel
	modErrorLabel lcl.ILabel
	modText       lcl.ILabel
	modBox        lcl.IComboBox
	//modLocalBox   lcl.IPanel
	//modRemoteBox  lcl.IPanel

	//modLocalDirText    lcl.ILabel
	//modLocalDirEdit    lcl.IEdit
	//modLocalDirBtn     *wg.TButton
	//modLCLCheckBox     lcl.ICheckBox
	//modCEFCheckBox     lcl.ICheckBox
	//modWebviewCheckBox lcl.ICheckBox

	// 操作按钮
	cancelBtn *wg.TButton
	createBtn *wg.TButton
}

func (m *TCreateProjectForm) FormCreate(sender lcl.IObject) {
	logs.Debug("TCreateProjectForm FormCreate")
	m.SetCaption("新建项目")
	m.SetWidth(createProjectFormWidth)
	m.SetHeight(createProjectFormHeight)
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

func (m *TCreateProjectForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	m.closing = true
}

func (m *TCreateProjectForm) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	*closeAction = types.CaFree
}

func (m *TCreateProjectForm) initComponents() {
	fontSize := int32(12)
	if tool.IsLinux {
		fontSize = 10
	}
	fontLabel := lcl.NewFont()
	fontLabel.SetName("微软雅黑")
	//fontLabel.SetStyle(types.NewSet(types.FsBold))
	fontLabel.SetSize(12)
	fontLabel.SetCharSet(font.CHINESEBIG5_CHARSET)
	//fontLabel.SetColor(colors.ClGreen)
	fontText := lcl.NewFont()
	fontText.SetName("微软雅黑 Light")
	fontText.SetSize(fontSize)

	left := int32(35)
	textWidth := int32(355)

	m.baseGroupBox = lcl.NewLabel(m)
	m.baseGroupBox.SetLeft(10)
	m.baseGroupBox.SetTop(10)
	m.baseGroupBox.SetCaption("新建项目-基础信息")
	m.baseGroupBox.SetFont(fontLabel)
	m.baseGroupBox.Font().SetSize(10)
	m.baseGroupBox.SetParent(m.box)

	m.modGroupBox = lcl.NewLabel(m)
	m.modGroupBox.SetTop(230)
	m.modGroupBox.SetLeft(10)
	m.modGroupBox.SetCaption("ENERGY框架-模块依赖")
	m.modGroupBox.SetFont(fontLabel)
	m.modGroupBox.Font().SetSize(10)
	m.modGroupBox.SetParent(m.box)

	m.selectDir = lcl.NewSelectDirectoryDialog(m)
	{
		errorFont := lcl.NewFont()
		errorFont.SetColor(colors.RGBToColor(255, 127, 127))
		errorFont.SetName("微软雅黑 Light")
		errorFont.SetCharSet(font.CHINESEBIG5_CHARSET)
		errorFont.SetSize(8)
		m.baseErrorLabel = lcl.NewLabel(m)
		m.baseErrorLabel.SetFont(errorFont)
		m.baseErrorLabel.SetVisible(false)
		m.baseErrorLabel.SetParent(m.box)
		m.modErrorLabel = lcl.NewLabel(m)
		m.modErrorLabel.SetFont(errorFont)
		m.modErrorLabel.SetVisible(false)
		m.modErrorLabel.SetParent(m.box)
	}
	baseTop := int32(40)
	{
		m.projNameText = lcl.NewLabel(m)
		m.projNameText.SetLeft(left)
		m.projNameText.SetTop(baseTop)
		m.projNameText.SetCaption("项目名称")
		m.projNameText.SetFont(fontLabel)
		m.projNameText.SetParent(m.box)

		m.projNameEdit = lcl.NewEdit(m)
		m.projNameEdit.SetLeft(120)
		m.projNameEdit.SetTop(baseTop - 5)
		m.projNameEdit.SetWidth(textWidth)
		m.projNameEdit.SetFont(fontText)
		m.projNameEdit.SetDoubleBuffered(true)
		m.projNameEdit.SetParentColor(false)
		m.projNameEdit.SetParent(m.box)
		m.projNameEdit.SetTextHint("YourProjectName")
	}
	{
		m.projPathText = lcl.NewLabel(m)
		m.projPathText.SetLeft(left)
		m.projPathText.SetTop(baseTop + 50)
		m.projPathText.SetCaption("项目路径")
		m.projPathText.SetFont(fontLabel)
		m.projPathText.SetParent(m.box)

		m.projPathEdit = lcl.NewEdit(m)
		m.projPathEdit.SetLeft(120)
		m.projPathEdit.SetTop(baseTop + 45)
		m.projPathEdit.SetWidth(textWidth - 65) // 是目录选择按钮的 宽度+Left(5)
		m.projPathEdit.SetFont(fontText)
		m.projPathEdit.SetDoubleBuffered(true)
		//m.projPathEdit.SetReadOnly(true)
		m.projPathEdit.SetTextHint("/your/app/path/name")
		m.projPathEdit.SetParent(m.box)

		m.projPathBtn = wg.NewButton(m)
		m.projPathBtn.SetIconFormBytes(resources.Images("actions/add.png"))
		m.projPathBtn.SetRadius(3)
		cusRect := types.TRect{Left: m.projPathEdit.Left() + m.projPathEdit.Width() + 5, Top: baseTop + 45}
		cusRect.SetWidth(60)
		if tool.IsLinux {
			cusRect.SetHeight(35)
		} else {
			cusRect.SetHeight(30)
		}
		m.projPathBtn.SetBoundsRect(cusRect)
		m.projPathBtn.SetParent(m.box)
		m.projPathBtn.SetOnClick(m.projPathClick)
	}
	{
		m.projTempText = lcl.NewLabel(m)
		m.projTempText.SetLeft(left)
		m.projTempText.SetTop(baseTop + 100)
		m.projTempText.SetCaption("项目模板")
		m.projTempText.SetFont(fontLabel)
		m.projTempText.SetParent(m.box)

		m.projTempBox = lcl.NewComboBox(m)
		m.projTempBox.SetBounds(120, baseTop+95, textWidth, 35)
		m.projTempBox.SetFont(fontText)
		m.projTempBox.SetReadOnly(true)
		m.projTempBox.SetStyle(types.CsDropDownList)
		m.projTempBox.SetBorderStyle(types.BsSingle)
		m.projTempBox.SetDoubleBuffered(true)
		m.projTempBox.Items().Add("默认预设模板")
		m.projTempBox.SetItemIndex(0)
		m.projTempBox.SetParent(m.box)
	}
	{
		m.goVersionText = lcl.NewLabel(m)
		m.goVersionText.SetLeft(left)
		m.goVersionText.SetTop(baseTop + 150)
		m.goVersionText.SetCaption(" Go 版本")
		m.goVersionText.SetFont(fontLabel)
		m.goVersionText.SetParent(m.box)

		m.goVersionStatus = wg.NewButton(m)
		m.goVersionStatus.SetText("检测本地")
		m.goVersionStatus.SetFont(fontText)
		m.goVersionStatus.Font().SetColor(colors.ClWhite)
		m.goVersionStatus.Font().SetStyle(types.NewSet(types.FsBold))
		m.goVersionStatus.SetRadius(3)
		goVersionRect := types.TRect{Left: 120, Top: baseTop + 145}
		goVersionRect.SetWidth(textWidth)
		goVersionRect.SetHeight(30)
		m.goVersionStatus.SetBoundsRect(goVersionRect)
		m.goVersionStatus.SetColor(colors.ClGray)
		m.goVersionStatus.SetShowHint(true)
		m.goVersionStatus.SetHint("在本机检测已安装可用的 Go 版本\n绿色: 检测成功, 支持\n红色: 检测失败, 不支持")
		m.goVersionStatus.SetParent(m.box)
	}

	baseTop = 210
	{
		m.modText = lcl.NewLabel(m)
		m.modText.SetLeft(left)
		m.modText.SetTop(baseTop + 50)
		m.modText.SetCaption("模块来源")
		m.modText.SetFont(fontLabel)
		m.modText.SetParent(m.box)

		m.modBox = lcl.NewComboBox(m)
		m.modBox.SetBounds(120, baseTop+45, textWidth, 36)
		m.modBox.SetFont(fontText)
		m.modBox.SetReadOnly(true)
		m.modBox.SetStyle(types.CsDropDownList)
		m.modBox.SetBorderStyle(types.BsSingle)
		for _, option := range modSelectOptions {
			m.modBox.Items().Add(option)
		}
		m.modBox.SetItemIndex(0)
		m.modBox.SetParent(m.box)
		m.modBox.SetOnChange(m.modBoxChange)

		//m.modLocalBox = lcl.NewPanel(m)
		//m.modLocalBox.SetBevelOuter(types.BvNone)
		//m.modLocalBox.SetTop(baseTop + 95)
		//m.modLocalBox.SetHeight(100)
		//m.modLocalBox.SetWidth(createProjectFormWidth)
		//m.modLocalBox.SetVisible(true)
		//m.modLocalBox.SetParent(m.box)
		//
		//m.modRemoteBox = lcl.NewPanel(m)
		//m.modRemoteBox.SetBevelOuter(types.BvNone)
		//m.modRemoteBox.SetTop(baseTop + 95)
		//m.modRemoteBox.SetHeight(100)
		//m.modRemoteBox.SetWidth(createProjectFormWidth)
		//m.modRemoteBox.SetVisible(false)
		//m.modRemoteBox.SetParent(m.box)

		//m.modLocalDirText = lcl.NewLabel(m)
		//m.modLocalDirText.SetLeft(left)
		//m.modLocalDirText.SetTop(5)
		//m.modLocalDirText.SetCaption("框架目录")
		//m.modLocalDirText.SetFont(fontLabel)
		//m.modLocalDirText.SetParent(m.modLocalBox)
		//
		//m.modLocalDirEdit = lcl.NewEdit(m)
		//m.modLocalDirEdit.SetLeft(120)
		//m.modLocalDirEdit.SetTop(0)
		//m.modLocalDirEdit.SetWidth(textWidth - 65) // 是目录选择按钮的 宽度+Left(5)
		//m.modLocalDirEdit.SetFont(fontText)
		//m.modLocalDirEdit.SetReadOnly(true)
		//m.modLocalDirEdit.SetDoubleBuffered(true)
		//m.modLocalDirEdit.SetText(frameworks.Path)
		//m.modLocalDirEdit.SetParent(m.modLocalBox)
		//
		//m.modLocalDirBtn = wg.NewButton(m)
		//m.modLocalDirBtn.SetIconFormBytes(resources.Images("actions/add.png"))
		//m.modLocalDirBtn.SetRadius(3)
		//modLocalDirRect := types.TRect{Left: m.modLocalDirEdit.Left() + m.modLocalDirEdit.Width() + 5, Top: 0}
		//modLocalDirRect.SetWidth(60)
		//if tool.IsLinux {
		//	modLocalDirRect.SetHeight(35)
		//} else {
		//	modLocalDirRect.SetHeight(30)
		//}
		//m.modLocalDirBtn.SetHint("默认 ENERGY Designer 安装目录")
		//m.modLocalDirBtn.SetShowHint(true)
		//m.modLocalDirBtn.SetBoundsRect(modLocalDirRect)
		//m.modLocalDirBtn.SetParent(m.modLocalBox)
		//m.modLocalDirBtn.SetOnClick(m.modLocalDirBtnClick)
		//
		//m.modLCLCheckBox = lcl.NewCheckBox(m)
		//m.modLCLCheckBox.SetFont(fontText)
		//m.modLCLCheckBox.SetLeft(left)
		//m.modLCLCheckBox.SetTop(modLocalDirRect.Top + modLocalDirRect.Height() + 20)
		//m.modLCLCheckBox.SetCaption("LCL (Native UI)")
		//m.modLCLCheckBox.SetHint("Lazarus Component Library")
		//m.modLCLCheckBox.SetShowHint(true)
		//m.modLCLCheckBox.SetChecked(true)
		//m.modLCLCheckBox.SetEnabled(frameworks.EnableLCL)
		//m.modLCLCheckBox.SetParent(m.modLocalBox)
		//
		//m.modCEFCheckBox = lcl.NewCheckBox(m)
		//m.modCEFCheckBox.SetFont(fontText)
		//m.modCEFCheckBox.SetLeft(left + 150)
		//m.modCEFCheckBox.SetTop(modLocalDirRect.Top + modLocalDirRect.Height() + 20)
		//m.modCEFCheckBox.SetCaption("CEF (Web UI)")
		//m.modCEFCheckBox.SetHint("Chromium Embedded Framework")
		//m.modCEFCheckBox.SetShowHint(true)
		//m.modCEFCheckBox.SetChecked(true)
		//m.modCEFCheckBox.SetEnabled(frameworks.EnableCEF)
		//m.modCEFCheckBox.SetParent(m.modLocalBox)
		//
		//m.modWebviewCheckBox = lcl.NewCheckBox(m)
		//m.modWebviewCheckBox.SetFont(fontText)
		//m.modWebviewCheckBox.SetLeft(left + 285)
		//m.modWebviewCheckBox.SetTop(modLocalDirRect.Top + modLocalDirRect.Height() + 20)
		//m.modWebviewCheckBox.SetCaption("WebView (Web UI)")
		//m.modWebviewCheckBox.SetHint("System Runtime Framework")
		//m.modWebviewCheckBox.SetShowHint(true)
		//m.modWebviewCheckBox.SetChecked(true)
		//m.modWebviewCheckBox.SetEnabled(frameworks.EnableWV)
		//m.modWebviewCheckBox.SetParent(m.modLocalBox)
	}
	{
		m.cancelBtn = wg.NewButton(m)
		m.cancelBtn.SetText("关 闭")
		m.cancelBtn.SetFont(fontText)
		m.cancelBtn.Font().SetColor(colors.ClWhite)
		m.cancelBtn.Font().SetStyle(types.NewSet(types.FsBold))
		m.cancelBtn.SetRadius(3)
		cancelBtnRect := types.TRect{Left: 250, Top: baseTop + 100}
		cancelBtnRect.SetWidth(100)
		cancelBtnRect.SetHeight(40)
		m.cancelBtn.SetBoundsRect(cancelBtnRect)
		m.cancelBtn.SetColor(colors.RGBToColor(255, 127, 127))
		m.cancelBtn.SetParent(m.box)
		m.cancelBtn.SetOnClick(m.closeClick)

		m.createBtn = wg.NewButton(m)
		m.createBtn.SetText("创 建")
		m.createBtn.SetFont(fontText)
		m.createBtn.Font().SetStyle(types.NewSet(types.FsBold))
		m.createBtn.Font().SetColor(colors.ClWhite)
		m.createBtn.SetRadius(3)
		createBtnRect := types.TRect{Left: cancelBtnRect.Left + cancelBtnRect.Width() + 30, Top: cancelBtnRect.Top}
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
	logs.Debug("TCreateProjectForm Show")
	m.one.Do(func() {
		addSize := int32(20)
		br := m.BoundsRect()
		br.SetWidth(createProjectFormWidth)
		br.SetHeight(createProjectFormHeight + addSize)
		m.SetBoundsRect(br) // trigger WM_NCCALCSIZE hook msg
		constr := m.Constraints()
		constr.SetMaxWidth(createProjectFormWidth)
		constr.SetMaxHeight(createProjectFormHeight + addSize)
		constr.SetMinWidth(createProjectFormWidth)
		constr.SetMinHeight(createProjectFormHeight + addSize)
		m.WorkAreaCenter()
		go m.checkGoVersion()
	})
}

func (m *TCreateProjectForm) checkGoVersion() {
	if m.goVersionStatus == nil {
		return
	}
	time.Sleep(time.Second / 2)
	result := false
	cmd := command.NewCMD()
	cmd.IsPrint = false
	cmd.HideWindow = true
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
					m.goVersionStatus.SetColor(colors.ClGreen)
					m.goVersionStatus.SetIconFavoriteFormBytes(resources.Images("button/laugh.png"))
				} else {
					m.goVersionStatus.SetColor(colors.ClRed)
					m.goVersionStatus.SetIconFavoriteFormBytes(resources.Images("button/weep.png"))
				}
				m.goVersionStatus.SetText(buf.String())
			})
		}
		result = true
	}
	cmd.Command("go", "version")
}

// 框架存放目录选择
//func (m *TCreateProjectForm) modLocalDirBtnClick(sender lcl.IObject) {
//	m.selectDir.SetTitle("框架安装目录")
//	m.selectDir.SetInitialDir(m.modLocalDirEdit.Text())
//	if m.selectDir.Execute() {
//		dir := m.selectDir.FileName()
//		m.modLocalDirEdit.SetText(dir)
//	}
//}

// 项目存放目录选择
func (m *TCreateProjectForm) projPathClick(sender lcl.IObject) {
	m.selectDir.SetTitle("新建项目")
	if m.projPathEdit.Text() == "" {
		initDir := gPath
		if initDir == "" {
			initDir = exec.Dir
		} else {
			initDir = filepath.Join(initDir, "../")
		}
		m.selectDir.SetInitialDir(initDir)
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
	// 禁用按钮
	m.cancelBtn.SetDisable(true)
	m.createBtn.SetDisable(true)
	recoverBtn := func() {
		if !m.closing {
			// 恢复按钮
			m.cancelBtn.SetDisable(false)
			m.createBtn.SetDisable(false)
		}
	}
	// 框架安装目录
	//frameworkDir := m.modLocalDirEdit.Text()
	//enableLCL := m.modLCLCheckBox.Checked()
	//enableCEF := m.modCEFCheckBox.Checked()
	//enableWV := m.modWebviewCheckBox.Checked()
	projectName := m.projNameEdit.Text()
	projectDir := m.projPathEdit.Text()
	// 检查创建项目
	if checkCreate(projectDir) {
		// 触发文件修改监听事件
		event.Emit(event.TTrigger{Name: event.ListenFileChange})
		// 允许创建
		// 重置设计器
		designer.ResetDesigner()
		go func() {
			// 更新设计器配置框架目录
			//if config.UpdateFrameworkDir(frameworkDir) {
			//	config.UpdateConfig()
			//	// 更新框架目录
			//	frameworks.Path = config.Config.FrameworkDir
			//}
			//// 释放 LCL 库
			//frameworks.ExtractLCL(enableLCL)
			//// 释放 CEF 库
			//frameworks.ExtractCEF(enableCEF)
			//// 释放 WebView 库
			//frameworks.ExtractWV(enableWV)
			// 运行创建项目
			if doRunCreate(projectName, projectDir) {
				// go.mod
				event.ConsoleWriteInfo("go mod tidy")
				cmd := command.NewCMD()
				cmd.IsPrint = false
				cmd.HideWindow = true
				cmd.Dir = projectDir
				cmd.Console = func(data string, level command.Level) {
					event.ConsoleWriteInfo(data)
				}
				cmd.Command("go", "mod", "tidy")
				// 恢复按钮
				recoverBtn()
				designer.UpdateDesignerTitle(fmt.Sprintf("%v (%v)", gProject.Name, gPath))
			}

		}()
	} else {
		recoverBtn()
	}
}

func (m *TCreateProjectForm) projIconPreviewPaintBackground(sender lcl.IObject, canvas lcl.ICanvas, rect types.TRect) {
	cell := int32(8)
	bmp := lcl.NewBitmap()
	bmp.SetPixelFormat(types.Pf24bit)
	bmp.SetSize(rect.Width(), rect.Height())
	bmpCanvas := bmp.Canvas()
	bmpCanvas.BrushToBrush().SetColor(colors.ClWhite)
	bmpCanvas.FillRectWithIntX4(0, 0, bmp.Width(), bmp.Height())
	bmpCanvas.BrushToBrush().SetColor(colors.ClLtGray)
	for i := 0; i < int(bmp.Width()/cell); i++ {
		for j := 0; j < int(bmp.Height()/cell); j++ {
			if (i%2 != 0) == (j%2 != 0) {
				bmpCanvas.FillRectWithIntX4(int32(i)*cell, int32(j)*cell, int32(i+1)*cell, int32(j+1)*cell)
			}
		}
	}
	sourceRect := types.TRect{Left: 0, Top: 0}
	sourceRect.SetWidth(bmp.Width())
	sourceRect.SetHeight(bmp.Height())
	canvas.CopyRectWithRectX2Canvas(rect, bmpCanvas, sourceRect)
}

// 模块选择
func (m *TCreateProjectForm) modBoxChange(sender lcl.IObject) {
	if m.modBox.ItemIndex() == 1 {
		m.showError(m.modErrorLabel, m.modText.BoundsRect(), "＊暂不支持从远程下载")
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

// 验证输入
func (m *TCreateProjectForm) validateInputs() bool {
	if strings.TrimSpace(m.projNameEdit.Text()) == "" {
		m.projNameEdit.SetFocus()
		m.showError(m.baseErrorLabel, m.projNameText.BoundsRect(), "＊项目名为空")
		return false
	}
	selectProjectPath := strings.TrimSpace(m.projPathEdit.Text())
	if selectProjectPath == "" {
		m.projPathEdit.SetFocus()
		m.showError(m.baseErrorLabel, m.projPathText.BoundsRect(), "＊项目目录为空")
		return false
	}
	//if !tool.IsExist(selectProjectPath) {
	//	//m.showError(m.baseErrorLabel, m.projPathText.BoundsRect(), "＊目录不存在")
	//	//isCreate := api.MessageDlg("目录不存在是否创建？", types.MtCustom,
	//	//	types.NewSet(types.MbYes, types.MbNo), types.MbNo) == types.IdYes
	//	//if isCreate {
	//	err := os.MkdirAll(selectProjectPath, os.ModePerm)
	//	if err != nil {
	//		m.showError(m.baseErrorLabel, m.projPathText.BoundsRect(), "＊创建目录失败")
	//		return false
	//	}
	//	//}
	//}
	if m.modBox.ItemIndex() == 1 {
		m.showError(m.modErrorLabel, m.modText.BoundsRect(), "＊暂不支持从远程下载")
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
