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
	"strings"
	"sync"
	"time"
)

// PERF 诊断开关 (临时): DESIGNER_PERF=1 时输出恢复各阶段耗时
var perfDiag = os.Getenv("DESIGNER_PERF") != ""

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
	perfTreeStart := time.Now()
	var perfCreateCost time.Duration
	var perfCreateCount int
	ProjectTreeBeginUpdate()
	defer ProjectTreeEndUpdate()

	var createTree func(childList []uiBean.TUIComponent, parent *TDesigningComponent)
	createTree = func(childList []uiBean.TUIComponent, parent *TDesigningComponent) {
		for _, child := range childList {
			perfStart := time.Now()
			newDesComp := GetDesignerComponent(parent.FormTab, 0, 0, child.ClassName)
			perfCreateCost += time.Since(perfStart)
			perfCreateCount++
			if newDesComp != nil {
				newDesComp.RecoverProperty = child.Properties
				// 设置组件节点关联
				newDesComp.SetParent(parent)
				// 添加到项目管理树节点
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
	if perfDiag {
		logs.Info("PERF RecoverComponentTree components:", perfCreateCount,
			"createCost:", perfCreateCost, "total:", time.Since(perfTreeStart))
	}
}

// 恢复设计窗体组件属性值
func (m *FormTab) RecoverComponentPropertyValue() {
	perfStart := time.Now()
	var perfCount int
	var iterateCreate func(childList []*TDesigningComponent)
	iterateCreate = func(childList []*TDesigningComponent) {
		for _, child := range childList {
			if child.RecoverProperty != nil {
				recoverDesignerComponentProperty(child)
				child.RecoverProperty = nil
				perfCount++
			}
			iterateCreate(child.Child)
		}
	}
	if m.FormRoot.RecoverProperty != nil {
		recoverDesignerComponentProperty(m.FormRoot)
		m.FormRoot.RecoverProperty = nil
	}
	iterateCreate(m.FormRoot.Child)
	if perfDiag {
		logs.Info("PERF RecoverComponentPropertyValue components:", perfCount, "cost:", time.Since(perfStart))
	}
}

// 恢复设计的组件属性
// 1. 调用 api 设置属性
// 2. 组件属性列表对应的属性Edit值
// Left/Top/Width/Height 几何属性合并为一次 SetBounds 直调:
// 设计模式下逐个属性 setter 走组件消息通知路径, 每次调用约 20-90ms;
// 而组件对象的 SetBounds 为快速路径, 100 组件布局恢复可从 10s+ 降到 1s 内。
func recoverDesignerComponentProperty(component *TDesigningComponent) {
	if component == nil || component.RecoverProperty == nil {
		return
	}
	perfStart := time.Now()
	component.GetProps()
	perfGetProps := time.Since(perfStart)
	perfCallStart := time.Now()
	var (
		geoLeft, geoTop, geoWidth, geoHeight int32
		hasLeft, hasTop, hasWidth, hasHeight bool
		applyGeo                             bool
	)
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
			// 几何属性收集后统一走 SetBounds 快速路径, 不逐个调用属性 setter
			propName := strings.ToLower(property.Name)
			geoHandled := false
			if component.ComponentType == consts.CtVisual {
				switch propName {
				case "left":
					geoLeft, hasLeft = anyToInt32(property.Value), true
					geoHandled = true
				case "top":
					geoTop, hasTop = anyToInt32(property.Value), true
					geoHandled = true
				case "width":
					geoWidth, hasWidth = anyToInt32(property.Value), true
					geoHandled = true
				case "height":
					geoHeight, hasHeight = anyToInt32(property.Value), true
					geoHandled = true
				}
			} else if component.ComponentType == consts.CtForm && (propName == "width" || propName == "height") {
				// 窗体位置由设计器管理, 只合并宽高
				if propName == "width" {
					geoWidth, hasWidth = anyToInt32(property.Value), true
				} else {
					geoHeight, hasHeight = anyToInt32(property.Value), true
				}
				geoHandled = true
			}
			if geoHandled {
				applyGeo = true
			} else {
				// 更新 api
				perfPropStart := time.Now()
				component.recoverCallAPI(property.Name, propNodeData)
				if perfDiag && time.Since(perfPropStart) > time.Millisecond {
					logs.Info("PERF prop slow", component.Name(), property.Name, time.Since(perfPropStart))
				}
				if component.ComponentType == consts.CtNonVisual {
					component.objectNonWrap.TextFollowShow()
				}
			}
		}
	}
	if applyGeo {
		perfPropStart := time.Now()
		br := component.Control().BoundsRect()
		l, t, w, h := br.Left, br.Top, br.Width(), br.Height()
		if hasLeft {
			l = geoLeft
		}
		if hasTop {
			t = geoTop
		}
		if hasWidth {
			w = geoWidth
		}
		if hasHeight {
			h = geoHeight
		}
		component.SetBounds(l, t, w, h)
		if perfDiag && time.Since(perfPropStart) > time.Millisecond {
			logs.Info("PERF geo slow", component.Name(), time.Since(perfPropStart))
		}
	}
	if perfDiag {
		logs.Info("PERF recoverProp", component.Name(), "getProps:", perfGetProps,
			"props:", len(propertyList), "callAPI:", time.Since(perfCallStart))
	}
}

// anyToInt32 将布局文件中的数值属性值转为 int32
func anyToInt32(v any) int32 {
	switch x := v.(type) {
	case float64:
		return int32(x)
	case float32:
		return int32(x)
	case int:
		return int32(x)
	case int64:
		return int32(x)
	case int32:
		return x
	case uint:
		return int32(x)
	case string:
		if i, err := tool.StrToInt32(x); err == nil {
			return i
		}
	}
	return 0
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
		// 加载所有窗体
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
				event.ConsoleWriteError("Failed to restore design window, error reading UI layout file:", err.Error())
				continue
			}
			uiComponent := &uiBean.TUIComponent{}
			err = json.Unmarshal(data, uiComponent)
			if err != nil {
				wg.Done()
				event.ConsoleWriteError("Failed to restore design window, error parsing form layout:", err.Error())
				continue
			}
			// 在UI线程操作
			lcl.RunOnMainThreadAsync(func(id uint32) {
				// 创建一个设计窗体
				formTab := designer.addDesignFormTab(&tempUIForm)
				if formTab == nil {
					event.ConsoleWriteError("recover design form error")
					return
				}
				// 恢复模式，在 tab page 显示时恢复组件属性
				formTab.recover = &TRecoverForm{
					components: uiComponent.Child,
					property:   uiComponent.Properties,
				}
				// 设置属性
				formTab.mainPage.Button().SetCaption(tempUIForm.Name)
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

// PERF 诊断 (临时): 恢复完成后对首个可视组件做基准测试,
// 对比恢复期间与恢复后的单次调用成本。
func perfBenchAfterRecover(m *FormTab) {
	var first *TDesigningComponent
	var walk func(c *TDesigningComponent)
	walk = func(c *TDesigningComponent) {
		if first != nil {
			return
		}
		for _, ch := range c.Child {
			if ch.ComponentType == consts.CtVisual {
				first = ch
				return
			}
			walk(ch)
		}
	}
	if m.FormRoot != nil {
		walk(m.FormRoot)
	}
	if first == nil || first.Control() == nil {
		return
	}
	bench := func(name string, fn func()) {
		fn()
		var total time.Duration
		for i := 0; i < 10; i++ {
			t := time.Now()
			fn()
			total += time.Since(t)
		}
		logs.Info("PERF bench", name, "avg:", total/10)
	}
	ctrl := first.Control()
	// 基准操作必须自恢复, 否则会污染组件的真实位置/尺寸/标题
	bench("SetCaption", func() {
		orig := ctrl.Caption()
		ctrl.SetCaption("B")
		ctrl.SetCaption(orig)
	})
	bench("SetHeight", func() {
		h := ctrl.Height()
		ctrl.SetHeight(h + 1)
		ctrl.SetHeight(h)
	})
	bench("SetLeft", func() {
		l := ctrl.Left()
		ctrl.SetLeft(l + 1)
		ctrl.SetLeft(l)
	})
	bench("SetBounds", func() { ctrl.SetBounds(ctrl.Left(), ctrl.Top(), ctrl.Width(), ctrl.Height()) })
	bench("BoundsRect读", func() { _ = ctrl.BoundsRect() })
}
