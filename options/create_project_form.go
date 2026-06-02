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
	"fmt"
	"github.com/energye/designer/designer"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/designer/resources/metadata"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/tool/command"
	"github.com/energye/lcl/tool/exec"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	createProjectFormWidth  = int32(500)
	createProjectFormHeight = int32(210)
	minGoVersion            = "1.20"
)

func init() {
	//if tool.IsLinux {
	//	createProjectFormWidth = 550
	//	createProjectFormHeight = 300
	//}
}

// NewCreateProjectForm 创建一个新的项目创建表单实例
// 该函数初始化一个 TCreateProjectForm 结构体，并通过 lcl.Application.NewForm 方法将其注册为应用程序窗体
func NewCreateProjectForm() *TCreateProjectForm {
	newEngForm := lcl.NewEngForm(nil)
	newForm := &TCreateProjectForm{TEngForm: *newEngForm.(*lcl.TEngForm)}
	//lcl.Application.NewForm(newForm) // 不使用原因：go debug 模式有问题
	newForm.FormCreate(newForm)
	newForm.SetOnCloseQuery(newForm.OnCloseQuery)
	newForm.SetOnClose(newForm.OnClose)
	return newForm
}

type TCreateProjectForm struct {
	lcl.TEngForm
	closing     bool
	goVersionOK bool
	one         sync.Once
	box         lcl.IPanel
	selectDir   lcl.ISelectDirectoryDialog

	// 基础信息部分
	projNameEdit    lcl.ILabeledEdit
	projPathEdit    lcl.ILabeledEdit
	projPathBtn     *wg.TButton
	goVersionStatus *wg.TButton

	// GUI 渲染框架
	guiRenderFrameworkText lcl.ILabel
	guiRenderFrameworkBox  lcl.IComboBox

	// 操作按钮
	cancelBtn *wg.TButton
	createBtn *wg.TButton
}

func (m *TCreateProjectForm) FormCreate(sender lcl.IObject) {
	logs.Debug("TCreateProjectForm FormCreate")
	m.SetName("CreateProjectForm")
	m.SetCaption(metadata.GI18n.Dict("CreateProjectForm.Caption"))
	m.SetWidth(createProjectFormWidth)
	m.SetHeight(createProjectFormHeight)
	m.SetVisible(false)
	m.SetDoubleBuffered(true)
	m.SetBorderIcons(types.NewSet(types.BiSystemMenu))
	m.box = lcl.NewPanel(m)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetAlign(types.AlClient)
	m.box.SetColor(colors.ClWhite)
	m.box.SetParent(m)
	m.SetOnShow(m.onShow)
	m.initComponents()
	SetWindowCenterByMainWindow(m)

	//(&hook.TWindowHook{Form: m}).Hook()
}

func (m *TCreateProjectForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	m.closing = true
}

func (m *TCreateProjectForm) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	*closeAction = types.CaFree
}

func (m *TCreateProjectForm) initComponents() {
	//left := int32(45)
	textWidth := int32(355)
	gTop := int32(0)
	nextTop := func(top int32) int32 {
		gTop += top
		return gTop
	}

	m.selectDir = lcl.NewSelectDirectoryDialog(m)
	m.selectDir.SetName("CreateProjectFormSelectDir")
	{
		m.projNameEdit = lcl.NewLabeledEdit(m)
		m.projNameEdit.SetName("CreateProjectFormProjectName")
		m.projNameEdit.SetLeft(100)
		if tool.IsLinux {
			m.projNameEdit.SetTop(nextTop(10))
		} else {
			m.projNameEdit.SetTop(nextTop(20))
		}
		m.projNameEdit.SetWidth(textWidth)
		m.projNameEdit.SetDoubleBuffered(true)
		m.projNameEdit.SetParentColor(false)
		m.projNameEdit.SetParent(m.box)
		m.projNameEdit.SetTextHint("your project name, e.g: myapp")
		m.projNameEdit.SetLabelPosition(types.LpLeft)
		m.projNameEdit.SetText("")
		projNameText := m.projNameEdit.EditLabel()
		projNameText.SetCaption(metadata.GI18n.Dict("CreateProjectFormProjectName.EditLabel.Caption"))
	}
	{

		m.projPathEdit = lcl.NewLabeledEdit(m)
		m.projPathEdit.SetName("CreateProjectFormProjectPath")
		m.projPathEdit.SetLeft(100)
		if tool.IsLinux {
			m.projPathEdit.SetTop(nextTop(40))
		} else {
			m.projPathEdit.SetTop(nextTop(30))
		}
		m.projPathEdit.SetWidth(textWidth - 40)
		m.projPathEdit.SetDoubleBuffered(true)
		m.projPathEdit.SetTextHint("/your/app/path/name")
		m.projPathEdit.SetParent(m.box)
		m.projPathEdit.SetLabelPosition(types.LpLeft)
		m.projPathEdit.SetText("")
		projPathText := m.projPathEdit.EditLabel()
		projPathText.SetCaption(metadata.GI18n.Dict("CreateProjectFormProjectPath.EditLabel.Caption"))

		m.projPathBtn = wg.NewButton(m)
		m.projPathBtn.SetIconFormBytes(resources.Images("menu/menu_project_open.png"))
		m.projPathBtn.SetRadius(3)
		cusRect := types.TRect{Left: m.projPathEdit.Left() + m.projPathEdit.Width() + 5, Top: m.projPathEdit.Top()}
		cusRect.SetWidth(35)
		if tool.IsLinux {
			cusRect.SetHeight(35)
		} else {
			cusRect.SetHeight(25)
		}
		m.projPathBtn.SetBoundsRect(cusRect)
		m.projPathBtn.SetColor(grayBtnColor)
		m.projPathBtn.SetBorderColor(wg.BbdNone, grayBtnColor)
		m.projPathBtn.SetCursor(types.CrHandPoint)
		m.projPathBtn.SetParent(m.box)
		m.projPathBtn.SetOnClick(m.projPathClick)
	}
	{
		goVersionText := lcl.NewLabel(m)
		goVersionText.SetName("CreateProjectFormGoVersionLabel")
		goVersionText.SetAutoSize(false)
		goVersionText.SetLeft(0)
		goVersionText.SetWidth(96)
		goVersionText.SetAlignment(types.TaRightJustify)
		goVersionText.SetTop(nextTop(40))
		goVersionText.SetCaption(metadata.GI18n.Dict("CreateProjectFormGoVersionLabel.Caption"))
		goVersionText.SetParent(m.box)

		m.goVersionStatus = wg.NewButton(m)
		m.goVersionStatus.SetRadius(0)
		goVersionRectTop := goVersionText.Top()
		if !tool.IsLinux {
			goVersionRectTop = goVersionRectTop - 5
		}
		goVersionRect := types.TRect{Left: 100, Top: goVersionRectTop}
		goVersionRect.SetWidth(textWidth)
		goVersionRect.SetHeight(25)
		color := wg.LightenColor(colors.RGBToColor(214, 234, 242), 0.5)
		m.goVersionStatus.SetBoundsRect(goVersionRect)
		m.goVersionStatus.SetColor(color)
		m.goVersionStatus.SetBorderColor(wg.BbdNone, wg.DarkenColor(color, 0.2))
		m.goVersionStatus.SetEnterColor(color, color)
		m.goVersionStatus.SetDownColor(color, color)
		m.goVersionStatus.AnchorSideTop().SetControl(goVersionText)
		m.goVersionStatus.AnchorSideTop().SetSide(types.AsrCenter)
		m.goVersionStatus.SetParent(m.box)
	}

	{
		m.guiRenderFrameworkText = lcl.NewLabel(m)
		m.guiRenderFrameworkText.SetName("CreateProjectFormGuiRenderFrameworkLabel")
		m.guiRenderFrameworkText.SetAutoSize(false)
		m.guiRenderFrameworkText.SetWidth(96)
		m.guiRenderFrameworkText.SetAlignment(types.TaRightJustify)
		m.guiRenderFrameworkText.SetLeft(0)
		m.guiRenderFrameworkText.SetTop(nextTop(35))
		m.guiRenderFrameworkText.SetCaption(metadata.GI18n.Dict("CreateProjectFormGuiRenderFrameworkLabel.Caption"))
		m.guiRenderFrameworkText.SetParent(m.box)

		m.guiRenderFrameworkBox = lcl.NewComboBox(m)
		m.guiRenderFrameworkBox.SetBounds(100, m.guiRenderFrameworkText.Top(), textWidth, 36)
		m.guiRenderFrameworkBox.AnchorSideTop().SetControl(m.guiRenderFrameworkText)
		m.guiRenderFrameworkBox.AnchorSideTop().SetSide(types.AsrCenter)

		m.guiRenderFrameworkBox.SetReadOnly(true)
		m.guiRenderFrameworkBox.SetStyle(types.CsDropDownList)
		m.guiRenderFrameworkBox.SetBorderStyle(types.BsSingle)
		bean.GUIRenderFrameworks.Iterate(func(gui string, guiDesc string) bool {
			m.guiRenderFrameworkBox.Items().Add(guiDesc)
			return false
		})
		m.guiRenderFrameworkBox.SetItemIndex(0)
		m.guiRenderFrameworkBox.SetParent(m.box)
	}
	{
		cancelBtnRect := types.TRect{Left: 325, Top: nextTop(45)}
		cancelBtnRect.SetWidth(50)
		cancelBtnRect.SetHeight(25)
		m.cancelBtn = wg.NewButton(m)
		m.cancelBtn.SetName("CreateProjectFormCloseBtn")
		m.cancelBtn.SetText(metadata.GI18n.Dict("CreateProjectFormCloseBtn.Caption"))
		m.cancelBtn.Font().SetSize(8)
		m.cancelBtn.SetRadius(3)
		m.cancelBtn.SetBoundsRect(cancelBtnRect)
		m.cancelBtn.SetColor(grayBtnColor)
		m.cancelBtn.SetCursor(types.CrHandPoint)
		m.cancelBtn.SetParent(m.box)
		m.cancelBtn.SetOnClick(m.closeClick)

		createBtnRect := types.TRect{Left: cancelBtnRect.Left + cancelBtnRect.Width() + 20, Top: cancelBtnRect.Top}
		createBtnRect.SetWidth(60)
		createBtnRect.SetHeight(25)
		m.createBtn = wg.NewButton(m)
		m.createBtn.SetName("CreateProjectFormCreateBtn")
		m.createBtn.SetText(metadata.GI18n.Dict("CreateProjectFormCreateBtn.Caption"))
		m.createBtn.Font().SetColor(colors.ClWhite)
		m.createBtn.Font().SetSize(8)
		m.createBtn.SetRadius(3)
		m.createBtn.SetBoundsRect(createBtnRect)
		m.createBtn.SetColor(blueBtnColor)
		m.createBtn.SetCursor(types.CrHandPoint)
		m.createBtn.SetParent(m.box)
		m.createBtn.SetOnClick(m.createClick)
	}
}

// 窗口显示事件
func (m *TCreateProjectForm) onShow(sender lcl.IObject) {
	logs.Debug("TCreateProjectForm Show")
	m.one.Do(func() {
		br := m.BoundsRect()
		br.SetWidth(createProjectFormWidth)
		br.SetHeight(createProjectFormHeight)
		m.SetBoundsRect(br) // trigger WM_NCCALCSIZE hook msg
		//constr := m.Constraints()
		//constr.SetMaxWidth(createProjectFormWidth)
		//constr.SetMaxHeight(createProjectFormHeight)
		//constr.SetMinWidth(createProjectFormWidth)
		//constr.SetMinHeight(createProjectFormHeight)
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
					m.goVersionStatus.SetIconFavoriteFormBytes(resources.Images("button/ok_btn_16x16.png"))
				} else {
					m.goVersionStatus.SetIconFavoriteFormBytes(resources.Images("button/err_btn_16x16.png"))
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
	m.selectDir.SetTitle(metadata.GI18n.Dict("CreateProjectFormSelectDir.Caption"))
	if m.projPathEdit.Text() == "" {
		initDir := bean.GPath
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

	projectName := m.projNameEdit.Text()
	projectDir := m.projPathEdit.Text()
	guiRenderFramework := m.guiRenderFrameworkBox.Text()
	guiRenderFrameworkGUI := bean.GUIRenderFramework("")
	bean.GUIRenderFrameworks.Iterate(func(gui, guiDesc string) bool {
		if guiRenderFramework == guiDesc {
			guiRenderFrameworkGUI = gui
			return true
		}
		return false
	})

	create := &CreateProject{Name: projectName, Dir: projectDir, GuiRenderFramework: guiRenderFrameworkGUI}
	if guiRenderFrameworkGUI == bean.GUIRenderFramework_CEF {
		// 弹出 CEF 版本选择/安装窗口
		chromiumForm := NewChromiumDirForm()
		chromiumForm.ShowModal()
		// 用户未点确认或安装未完成, 中止创建
		if !chromiumForm.Confirmed || chromiumForm.Version == "" ||
			!config.Config.Chromium.IsCEFInstalled(chromiumForm.Version) {
			return
		}
		create.FrameworkVersion = chromiumForm.Version
	}

	// 检查创建项目
	if checkCreate(projectDir) {
		m.Close()
		// 触发文件修改监听事件
		event.Emit(event.TTrigger{Name: event.ListenFileChange})
		// 重置设计器
		designer.ResetDesigner()
		go func() {
			// 运行创建项目
			if doRunCreate(create) {
				lcl.RunOnMainThreadAsync(func(id uint32) {
					designer.ProjectTreeClear()
					event.Emit(event.TTrigger{Name: event.ListenProjectSrcFileChange, Payload: event.TPayload{Type: event.ProjectSrcScan}})
				})
				designer.CMDGoModDepsUpdate()
				designer.UpdateDesignerTitle(fmt.Sprintf("%v (%v)", bean.GProject.Name, bean.GPath))
				designer.SetEnableFuncComponent(true)
				//recoverBtn()
			}
		}()
	} else {
		//recoverBtn()
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
		return false
	}
	selectProjectPath := strings.TrimSpace(m.projPathEdit.Text())
	if selectProjectPath == "" {
		m.projPathEdit.SetFocus()
		return false
	}
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
