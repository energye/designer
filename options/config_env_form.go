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
	"errors"
	"github.com/energye/designer/designer"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"os"
	"path/filepath"
	"strings"
)

var (
	envFormWidth   = int32(555)
	envFormHeight  = int32(120)
	supportedLangs = []struct {
		code string
		name string
	}{
		{"zh-CN", "简体中文"},
		{"en-US", "English (US)"},
		//{"ja", "日本語"},
		//{"ko", "한국어"},
	}
)

// 运行项目(环境)配置窗口
func runEnvConfig() {
	// 显示运行项目(环境)配置窗口
	lcl.RunOnMainThreadAsync(func(id uint32) {
		form := NewEnvForm()
		form.ShowModal()
		//form.Show()
	})
}

func NewEnvForm() *TEnvForm {
	newEngForm := lcl.NewEngForm(nil)
	newForm := &TEnvForm{TEngForm: *newEngForm.(*lcl.TEngForm)}
	newForm.FormCreate(newEngForm)
	newForm.SetOnCloseQuery(newForm.OnCloseQuery)
	newForm.SetOnClose(newForm.OnClose)
	return newForm
}

type TEnvForm struct {
	lcl.TEngForm
	closing   bool
	selectDir lcl.ISelectDirectoryDialog

	goRootText lcl.ILabel
	goRootBox  lcl.IComboBox
	goRootBtn  *wg.TButton

	currentLang  string
	langComboBox lcl.IComboBox

	// 操作按钮
	cancelBtn *wg.TButton
	saveBtn   *wg.TButton
}

func (m *TEnvForm) FormCreate(sender lcl.IObject) {
	logs.Debug("TEnvForm FormCreate")

	m.SetCaption("环境配置")
	m.SetWidth(envFormWidth)
	m.SetHeight(envFormHeight)
	//constr := m.Constraints()
	//constr.SetMaxWidth(envFormWidth)
	//constr.SetMaxHeight(envFormHeight)
	//constr.SetMinWidth(envFormWidth)
	//constr.SetMinHeight(envFormHeight)
	m.SetVisible(false)
	m.SetDoubleBuffered(true)
	m.SetBorderIcons(types.NewSet(types.BiSystemMenu))
	SetWindowCenterByMainWindow(m)
	m.SetColor(colors.ClWhite)

	m.selectDir = lcl.NewSelectDirectoryDialog(m)

	gTop := int32(0)
	nextTop := func(top int32) int32 {
		gTop += top
		return gTop
	}
	{
		m.goRootText = lcl.NewLabel(m)
		m.goRootText.SetLeft(10)
		m.goRootText.SetTop(15)
		m.goRootText.SetCaption("Go Root")
		m.goRootText.SetParent(m)
		m.goRootBox = lcl.NewComboBox(m)
		m.goRootBox.SetLeft(10 + 55)
		m.goRootBox.SetTop(nextTop(10))
		m.goRootBox.SetWidth(440)
		m.goRootBox.SetDoubleBuffered(true)
		//m.goRootBox.SetStyle(types.CsDropDownList)
		m.goRootBox.SetBorderStyle(types.BsSingle)
		m.goRootBox.SetParent(m)
		// 优先从配置
		env := config.Config.Env[bean.GProject.Name]
		if env != nil && len(env.GoRoot) > 0 {
			goRoot := env.GoRoot
			for _, option := range goRoot {
				m.goRootBox.Items().Add(option)
			}
			selectIndex := env.GoRootSelectIndex
			if selectIndex >= 0 && selectIndex < int32(len(goRoot)) {
				m.goRootBox.SetText(goRoot[selectIndex])
				_ = SetGoRootPath(goRoot[selectIndex])
			}
		} else if cmdGoRoot != "" {
			// 从命令行中获取
			m.goRootBox.Items().Add(cmdGoRoot)
			m.goRootBox.SetText(cmdGoRoot)
		}

		m.goRootBtn = wg.NewButton(m)
		m.goRootBtn.SetIconFormBytes(resources.Images("menu/menu_project_open.png"))
		m.goRootBtn.SetRadius(3)
		cusRect := types.TRect{Left: m.goRootBox.Left() + m.goRootBox.Width() + 5, Top: 10}
		cusRect.SetWidth(30)
		if tool.IsLinux {
			cusRect.SetHeight(35)
		} else {
			cusRect.SetHeight(25)
		}
		m.goRootBtn.SetBoundsRect(cusRect)
		m.goRootBtn.SetColor(grayBtnColor)
		m.goRootBtn.SetBorderColor(wg.BbdNone, grayBtnColor)
		m.goRootBtn.SetCursor(types.CrHandPoint)
		m.goRootBtn.SetParent(m)
		m.goRootBtn.SetOnClick(m.goRootClick)
	}
	{

		m.langComboBox = lcl.NewComboBox(m)
		m.langComboBox.SetName("LangComboBox")
		m.langComboBox.SetLeft(m.goRootBox.Left())
		m.langComboBox.SetTop(nextTop(30))
		m.langComboBox.SetWidth(m.goRootBox.Width())
		m.langComboBox.SetStyle(types.CsDropDownList)
		selectIndex := 0
		for idx, lang := range supportedLangs {
			m.langComboBox.Items().Add(lang.name)
			if lang.code == config.Config.EnvLang {
				selectIndex = idx
			}
		}
		m.langComboBox.SetItemIndex(int32(selectIndex))
		m.langComboBox.SetParent(m)
		m.langComboBox.SetOnChange(m.onLangChange)
	}

	{
		cancelBtnRect := types.TRect{Left: 400, Top: nextTop(35)}
		cancelBtnRect.SetWidth(60)
		cancelBtnRect.SetHeight(25)
		m.cancelBtn = wg.NewButton(m)
		m.cancelBtn.SetText("关　闭")
		m.cancelBtn.Font().SetSize(8)
		m.cancelBtn.SetRadius(3)
		m.cancelBtn.SetBoundsRect(cancelBtnRect)
		m.cancelBtn.SetColor(grayBtnColor)
		m.cancelBtn.SetCursor(types.CrHandPoint)
		m.cancelBtn.SetParent(m)
		m.cancelBtn.SetOnClick(m.closeClick)

		saveBtnRect := types.TRect{Left: cancelBtnRect.Left + cancelBtnRect.Width() + 20, Top: cancelBtnRect.Top}
		saveBtnRect.SetWidth(60)
		saveBtnRect.SetHeight(25)
		m.saveBtn = wg.NewButton(m)
		m.saveBtn.SetText("保　存")
		m.saveBtn.Font().SetSize(8)
		m.saveBtn.Font().SetColor(colors.ClWhite)
		m.saveBtn.SetRadius(3)
		m.saveBtn.SetBoundsRect(saveBtnRect)
		m.saveBtn.SetColor(blueBtnColor)
		m.saveBtn.SetCursor(types.CrHandPoint)
		m.saveBtn.SetParent(m)
		m.saveBtn.SetOnClick(m.saveClick)
	}
	//(&hook.TWindowHook{Form: m}).Hook()
}

func (m *TEnvForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
	m.closing = true
}

func (m *TEnvForm) OnClose(sender lcl.IObject, closeAction *types.TCloseAction) {
	*closeAction = types.CaFree
}

func (m *TEnvForm) closeClick(sender lcl.IObject) {
	m.Close()
}

func (m *TEnvForm) saveClick(sender lcl.IObject) {
	goRoot := strings.TrimSpace(m.goRootBox.Text())
	event.ConsoleWriteInfo("Environment Configuration - Save", goRoot)
	if goRoot != "" {
		err := SetGoRootPath(goRoot)
		if err == nil {
			// 更新到 .energy 配置文件
			config.UpdateEnvGoRoot(bean.GProject.Name, goRoot)
		} else {
			event.ConsoleWriteError("Environment Configuration - Save failed:", err.Error(), goRoot)
			return
		}
	}
	config.UpdateConfig()
	event.ConsoleWriteInfo("Environment Configuration - Save-Completed")
}

func (m *TEnvForm) goRootClick(sender lcl.IObject) {
	m.selectDir.SetTitle("设置 Go SDK 目录")
	if m.selectDir.Execute() {
		goRoot := m.selectDir.FileName()
		bin := filepath.Join(goRoot, "bin")
		goCmd := "go"
		if tool.IsWindows {
			goCmd = "go.exe"
		}
		if !tool.IsExist(filepath.Join(bin, goCmd)) {
			api.ShowMessage("目录不是有效的 Go SDK 目录")
			return
		}
		isAdd := true
		for i := int32(0); i < m.goRootBox.Items().Count(); i++ {
			tmpItem := m.goRootBox.Items().Strings(i)
			if tmpItem == goRoot {
				isAdd = false
				break
			}
		}
		if isAdd {
			m.goRootBox.Items().Add(goRoot)
			m.goRootBox.SetText(goRoot)
		}
	}
}

var (
	currentEnvPath string
)

// SetGoRootPath 设置Go根路径到系统PATH环境变量中
//   - path: 要添加到PATH中的Go根路径
func SetGoRootPath(goSDKHome string) error {
	if currentEnvPath == "" {
		currentEnvPath = os.Getenv("PATH")
	}
	goSDKHome = strings.TrimSpace(goSDKHome)
	if goSDKHome == "" {
		return errors.New("Go SDK HOME 为空")
	}
	absPath, err := filepath.Abs(goSDKHome)
	if err != nil {
		return errors.New("路径格式非法：" + err.Error())
	}
	bin := filepath.Join(absPath, "bin")
	goCmd := "go"
	if tool.IsWindows {
		goCmd = "go.exe"
	}
	if !tool.IsExist(filepath.Join(bin, goCmd)) {
		return errors.New("目录不是有效的 Go SDK 目录")
	}
	pathList := strings.Split(currentEnvPath, string(os.PathListSeparator))
	newPathList := make([]string, 0, len(pathList))
	for _, p := range pathList {
		pTrim := strings.TrimSpace(p)
		if pTrim == "" {
			continue
		}
		normalizedP, err := filepath.Abs(pTrim)
		if err != nil {
			continue
		}
		if normalizedP != bin {
			newPathList = append(newPathList, p)
		}
	}
	pathList = newPathList
	newPath := bin + string(os.PathListSeparator) + strings.Join(pathList, string(os.PathListSeparator))
	if err = os.Setenv("GOROOT", absPath); err != nil {
		return errors.New("设置 GoRoot 环境变量失败：" + err.Error())
	}
	if err = os.Setenv("PATH", newPath); err != nil {
		return errors.New("设置 PATH 环境变量失败：" + err.Error())
	}
	return nil
}

func (m *TEnvForm) onLangChange(sender lcl.IObject) {
	idx := m.langComboBox.ItemIndex()
	if idx < 0 || int(idx) >= len(supportedLangs) {
		return
	}
	lang := supportedLangs[idx].code
	if lang == m.currentLang {
		return
	}
	m.currentLang = lang
	if m.currentLang != "" {
		err := designer.SwitchLocalesI18n(m.currentLang)
		if err != nil {
			event.ConsoleWriteError("Environment Configuration - i18n:", err.Error())
			return
		}
		config.Config.EnvLang = m.currentLang
		config.UpdateConfig()
		event.ConsoleWriteInfo("Environment Configuration - Switch i18n-Completed")
	}
}
