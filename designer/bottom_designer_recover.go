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
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
	projBean "github.com/energye/designer/project/bean"
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

// 恢复窗体组件信息
// 只恢复一次
func (m *FormTab) Recover() {
	if m.recover == nil {
		return
	}
	tempRecover := m.recover
	// 置空
	m.recover = nil
	// 加载属性到设计器
	// 此步骤会初始化并填充设计组件实例
	m.FormRoot.LoadPropertyToInspector()
	// 添加到组件树
	m.AddFormNode()
	// 恢复属性
	recoverDesignerComponentProperty(tempRecover.property, m.FormRoot)
	// 恢复子组件
	recoverDesignerChildComponent(tempRecover.components, m.FormRoot)
	// 恢复的默认切换至当前Form编辑状态
	m.FormRoot.node.SetSelected(true)
	//m.switchComponentEditing(m.FormRoot)
	// 释放掉
	tempRecover.components = nil
	tempRecover.property = nil
}

// 恢复设计的子组件
func recoverDesignerChildComponent(childList []uiBean.TUIComponent, parent *TDesigningComponent) {
	for _, child := range childList {
		if create := GetRegisterComponent(child.ClassName); create != nil {
			newComp := create(parent.formTab, 0, 0)
			newComp.SetParent(parent)
			// 2. 添加到组件树
			parent.AddChild(newComp)
			// 加载属性到设计器
			// 此步骤会初始化并填充设计组件实例
			newComp.LoadPropertyToInspector()
			// 恢复组件属性
			recoverDesignerComponentProperty(child.Properties, newComp)
			// 恢复子组件
			recoverDesignerChildComponent(child.Child, newComp)
		}
	}
}

// 恢复设计的组件属性
// 1. 调用 api 设置属性
// 2. 组件属性列表对应的属性Edit值
func recoverDesignerComponentProperty(propertyList []uiBean.TProperty, component *TDesigningComponent) {
	for _, property := range propertyList {
		for _, prop := range component.PropertyList {
			if prop.Name() == property.Name {
				// 设置属性值
				prop.SetEditValue(property.Value)
				// 重新渲染该属性单元格, 这里应该不需要
				//component.propertyTree.InvalidateNode(prop.AffiliatedNode)
				// 更新 api
				component.doUpdateComponentPropertyToObject(prop)
				break
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
			uiFilePath := filepath.Join(path, project.Package, tempUIForm.UIFile)
			data, err := os.ReadFile(uiFilePath)
			if err != nil {
				wg.Done()
				logs.Error("恢复设计窗体, 读取UI布局文件错误:", err.Error())
				event.ConsoleWriteError("恢复设计窗体, 读取UI布局文件错误:", err.Error())
				continue
			}
			uiComponent := &uiBean.TUIComponent{}
			err = json.Unmarshal(data, uiComponent)
			if err != nil {
				wg.Done()
				logs.Error("恢复设计窗体, 解析窗体布局错误:", err.Error())
				event.ConsoleWriteError("恢复设计窗体, 解析窗体布局错误:", err.Error())
				continue
			}
			// 在UI线程操作
			lcl.RunOnMainThreadAsync(func(id uint32) {
				// 创建一个设计窗体
				formTab := designer.addDesignerFormTab(tempUIForm.Id)
				formTab.recover = &TRecoverForm{
					components: uiComponent.Child,
					property:   uiComponent.Properties,
				}
				// 设置属性
				formTab.sheet.Button().SetCaption(tempUIForm.Name)

				// 默认激活的窗体
				if project.ActiveUIForm == formTab.Id {
					activeForm = formTab
				}
				wg.Done()
			})
		}
		// 等待所有设计窗体创建完
		wg.Wait()
		lcl.RunOnMainThreadAsync(func(id uint32) {
			designer.tab.RecalculatePosition()
			if activeForm != nil {
				// 隐藏掉所有组件树
				designer.hideAllComponentTrees()
				// 隐藏掉所有 form tab
				designer.tab.HideAllActivated()
				// 激活显示当前默认的 form tab
				designer.ActiveFormTab(activeForm)
			}
		})
	}
}
