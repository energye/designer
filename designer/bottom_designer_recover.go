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

package designer

import (
	"encoding/json"
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	projBean "github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	uiBean "github.com/energye/designer/uigen/bean"
	"github.com/energye/lcl/lcl"
	"os"
	"path/filepath"
	"sync"
)

// 恢复 FormTab
// 从 UI 布局文件恢复
// 触发功能:
//  1. 打开 xxx.ui 布局文件
//	1.1 恢复成功后, 提示更新到项目配置？？TODO 不同目录是否可以？还是复制到此项目目录？
//  2. 打开项目配置文件 xxx.egp, 根据 ui_forms 字段恢复所有窗体
//	2.1 恢复所有窗体对象到设计器

type TRecoverForm struct {
	components []uiBean.TUIComponent
	property   []uiBean.TProperty
}

// 恢复设计窗体组件树
//
//	恢复当前设计窗体的所有设计组件对象
//	创建设计组件对象并添加到项目管理树节点, 不恢复组件属性值
func (m *FormTab) RecoverComponentTree() {
	ProjectTreeBeginUpdate()
	defer ProjectTreeEndUpdate()

	var createTree func(childList []uiBean.TUIComponent, parent *TDesigningComponent)
	createTree = func(childList []uiBean.TUIComponent, parent *TDesigningComponent) {
		for _, child := range childList {
			if newDesComp := GetDesignerComponent(parent.FormTab, 0, 0, child.ClassName); newDesComp != nil {
				newDesComp.RecoverProperty = child.Properties
				// 设置组件节点关联
				newDesComp.SetParent(parent)
				// 添加到项目管理树
				parent.AddChild(newDesComp)
				// 恢复子组件
				createTree(child.Child, newDesComp)
			}
		}
	}
	// 项目管理-窗体根节点
	m.FormRoot.RecoverProperty = m.recover.property
	m.AddFormNode()
	// 创建子节点
	createTree(m.recover.components, m.FormRoot)
}

// 恢复设计窗体组件属性值
func (m *FormTab) RecoverComponentPropertyValue() {
	var iterateCreate func(childList []*TDesigningComponent)
	iterateCreate = func(childList []*TDesigningComponent) {
		for _, child := range childList {
			if child.RecoverProperty != nil {
				recoverDesignerComponentProperty(child)
				child.RecoverProperty = nil
			}
			iterateCreate(child.Child)
		}
	}
	if m.FormRoot.RecoverProperty != nil {
		recoverDesignerComponentProperty(m.FormRoot)
		m.FormRoot.RecoverProperty = nil
	}
	iterateCreate(m.FormRoot.Child)
}

// 恢复设计的组件属性
// 1. 调用 api 设置属性
// 2. 组件属性列表对应的属性Edit值
func recoverDesignerComponentProperty(component *TDesigningComponent) {
	if component == nil || component.RecoverProperty == nil {
		return
	}
	component.GetProps()
	propertyList := component.RecoverProperty
	for _, property := range propertyList {
		propNodeData := component.FindNodeDataByNamePaths(property)
		if propNodeData != nil {
			if propNodeData.Type() == consts.PdtCheckBoxList {
				// 转换 set 集合为 hashSet
				set := tool.SetToHashSet(property.Value)
				for _, checkBox := range propNodeData.EditNodeData.CheckBoxValue {
					checkBox.Checked = set.Contains(checkBox.Name)
				}
			}
			// 设置属性值
			propNodeData.SetEditValue(property.Value)
			// 更新 api
			component.recoverCallAPI(property.Name, propNodeData)
			if component.ComponentType == consts.CtNonVisual {
				component.objectNonWrap.TextFollowShow()
			}
		}
	}
}

// RecoverDesignerFormTab 恢复设计窗体, 非线程安全
// 只恢复当前项目下的窗体
// path: 当前项目路径
// project: 项目对象
// loadUIForm: 要加载的 UI 窗体对象, 如果 nil 表示加载所有窗体, 否则只加载当前这个窗体
func RecoverDesignerFormTab(path string, project *projBean.TProject, loadUIForm *projBean.TUIForm) {
	if designer == nil {
		return
	}
	if loadUIForm != nil {
		// 只加载这个窗体
	} else {
		// 加载所有
		wg := sync.WaitGroup{}
		wg.Add(len(project.UIForms))
		var activeForm *FormTab
		for _, uiForm := range project.UIForms {
			tempUIForm := uiForm
			// 判断窗体是整已存在
			if tab := designer.GetFormTab(tempUIForm.Id); tab != nil {
				wg.Done()
				if project.ActiveUIForm == tab.Id {
					activeForm = tab
				}
				continue
			}
			uiFilePath := filepath.Join(path, consts.LayoutsDir, tempUIForm.UIFile)
			data, err := os.ReadFile(uiFilePath)
			if err != nil {
				wg.Done()
				event.ConsoleWriteError("恢复设计窗体, 读取UI布局文件错误:", err.Error())
				continue
			}
			uiComponent := &uiBean.TUIComponent{}
			err = json.Unmarshal(data, uiComponent)
			if err != nil {
				wg.Done()
				event.ConsoleWriteError("恢复设计窗体, 解析窗体布局错误:", err.Error())
				continue
			}
			// 在UI线程操作
			lcl.RunOnMainThreadAsync(func(id uint32) {
				// 创建一个设计窗体
				formTab := designer.addDesignerFormTab(tempUIForm.Id)
				// 恢复模式，在 tab page 显示时恢复组件属性
				formTab.recover = &TRecoverForm{
					components: uiComponent.Child,
					property:   uiComponent.Properties,
				}
				// 设置属性
				formTab.sheet.Button().SetCaption(tempUIForm.Name)
				// 创建项目树节点->窗体节点
				formTab.RecoverComponentTree()
				// 激活最后设计的窗体
				if project.ActiveUIForm == formTab.Id {
					activeForm = formTab
				}
				wg.Done()
			})
		}
		// 等待所有设计窗体创建完
		wg.Wait()
		logs.Debug("RecoverDesignerFormTab end")
		lcl.RunOnMainThreadAsync(func(id uint32) {
			designer.tab.RecalculatePosition()
			if activeForm != nil {
				activeForm.FormRoot.node.SetExpanded(true)
				// 隐藏掉所有 form tab page
				designer.tab.HideAllActivated()
				// 激活显示当前默认的 form tab
				designer.ActiveFormTab(activeForm)
			}
		})
	}
}
