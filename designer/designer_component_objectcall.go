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
	"errors"
	"fmt"
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/mapper"
	"github.com/energye/designer/pkg/message"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"reflect"
	"strings"
	"sync/atomic"
	"time"
)

// 组件对象函数调用

var (
	updateComponentMutex atomic.Bool
	updateComponentDelay = 50 * time.Millisecond
)

func methodNameToSet(name string) string {
	name = tool.FirstToUpper(name)
	return "Set" + name
}

// 恢复调用 api
// property 当前恢复的属性
func (m *TDesigningComponent) recoverCallAPI(propertyName string, property *vtedit.TEditNodeData) {
	if property.Type() == consts.PdtMethod {
		// 绑定事件方法不更新 API
		return
	}
	if property == nil {
		logs.Error("恢复属性-属性名节点数据不存在:", propertyName)
	} else {
		switch rs := m.CheckCanUpdateProp(property); rs {
		case err.RsIgnoreProp:
			return
		}
		ref := &reflector{object: m.originObject, data: property, objectNonWrap: m.objectNonWrap}
		_, err := ref.callMethod()
		if err != nil {
			logs.Error("恢复属性-调用 API 更新组件属性失败", err.Error())
		} else {
			switch strings.ToLower(propertyName) {
			case "name":
				m.SetName(property.EditStringValue())
			}
		}
	}
}

// 执行更新组件绑定事件到代码, 该操作不更新 api 绑定, 而是直接更新代码
func (m *TDesigningComponent) doUpdateComponentBindEventToCode(updateNodeData *vtedit.TEditNodeData) {
	logs.Debug("更新组件:", m.ClassName(), "事件:", updateNodeData.Name(), "IsModify:", updateNodeData.IsModify())
	triggerUIGeneration(m, updateNodeData, event.CodeGenEvent)
}

// 执行更新组件属性到对象 api
func (m *TDesigningComponent) doUpdateComponentPropertyToObject(updateNodeData *vtedit.TEditNodeData) {
	if updateNodeData.Type() == consts.PdtMethod {
		// 绑定事件方法不更新 API
		return
	}
	logs.Debug("更新组件:", m.ClassName(), "属性:", updateNodeData.Name(), "IsModify:", updateNodeData.IsModify())
	if !m.node.IsValid() {
		// 无效节点对象
		return
	}
	// 检查当前组件属性是否允许更新
	if rs := m.CheckCanUpdateProp(updateNodeData); rs == err.RsSuccess {
		logs.Info("检查更新属性-成功, 该属性", updateNodeData.Name(), "调用 API 更新, 同时更新节点数据")
		ref := &reflector{object: m.originObject, data: updateNodeData, objectNonWrap: m.objectNonWrap}
		result, err := ref.callMethod()
		_ = result
		if err != nil {
			logs.Error("调用 API 更新组件属性失败", err.Error())
		} else {
			// 成功
			logs.Info("调用 API 更新组件属性成功, 更新节点数据")
			m.UpdateTreeNode(updateNodeData)
			// 属性修改-UI布局
			triggerUIGeneration(m, updateNodeData, event.CodeGenUI)
		}
	} else if rs == err.RsIgnoreProp { // 忽略的属性, 成功的一种
		logs.Info("检查更新属性-忽略, 该属性", updateNodeData.Name(), "忽略 API 更新, 只更新节点数据")
		m.UpdateTreeNode(updateNodeData)
		// 属性修改-UI布局
		triggerUIGeneration(m, updateNodeData, event.CodeGenUI)
	} else {
		// 更新失败
		switch rs {
		case err.RsDuplicateName: // 重复的组件名
			logs.Error("检查更新属性-组件名重复 检查更新属性失败, RS:", rs, "恢复节点内的组件名")
			// 恢复节点内的组件名
			updateNodeData.SetEditValue(m.Name())
			m.propertyTree.InvalidateNode(updateNodeData.AffiliatedNode)
		default:
			logs.Error("检查更新属性-其它错误 检查更新属性失败, RS:", rs)
		}
	}
}

// 更新组件属性到对象
func (m *TDesigningComponent) UpdateComponentPropertyToObject(updateNodeData *vtedit.TEditNodeData) {
	if updateComponentMutex.Load() {
		return
	}
	updateComponentMutex.Store(true)
	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.doUpdateComponentPropertyToObject(updateNodeData)
		updateComponentMutex.Store(false)
	})
}

// 更新组件绑定事件到代码
func (m *TDesigningComponent) UpdateComponentBindEventToCode(updateNodeData *vtedit.TEditNodeData) {
	if updateComponentMutex.Load() {
		return
	}
	updateComponentMutex.Store(true)
	m.doUpdateComponentBindEventToCode(updateNodeData)
	updateComponentMutex.Store(false)
}

// GetRecvMethods 获取接收方法列表
// 返回与当前设计组件关联的接收方法信息数组
// 这些方法通常是在代码生成过程中需要绑定到组件上的事件处理方法
func (m *TDesigningComponent) GetRecvMethods() []*dast.TFuncInfo {
	return m.FormTab.recvMethods
}

func (m *TDesigningComponent) GetClassName() string {
	className := m.ClassName()
	return className
}

func (m *TDesigningComponent) GetName() string {
	name := m.Name()
	return name
}

func (m *TDesigningComponent) GetMod() string {
	if m.mod == "" {
		// 默认 LCL
		m.mod = consts.ModLCL
	}
	return m.mod
}

// 更新组件树节点信息
// 在设计组件属性修改后同步修改组件树节点可见值
func (m *TDesigningComponent) UpdateTreeNode(updateNodeData *vtedit.TEditNodeData) error {
	if !m.node.IsValid() {
		logs.Error("更新组件树失败, 当前设计组件节点无效")
		return errors.New("更新组件树失败, 当前设计组件节点无效")
	}
	data := updateNodeData.EditNodeData
	propName := strings.ToLower(data.Name)
	logs.Debug("更新组件树, 尝试更新属性:", data.Name)
	switch propName {
	case "name":
		// 同步更新组件名字段值
		m.SetName(data.EditStringValue())
		// 更新组件树名
		m.node.SetText(m.TreeName())
		// 窗体组件
		if m.ComponentType == consts.CtForm {
			// 更新设计窗体标签名
			m.FormTab.mainPage.SetCaption(m.Name())
		}
	}
	return nil
}

// 检查是否允许更新属性
func (m *TDesigningComponent) CheckCanUpdateProp(updateNodeData *vtedit.TEditNodeData) err.ResultStatus {
	propName := strings.ToLower(updateNodeData.Name())
	switch propName {
	case "name":
		if !tool.IsValidVariableName(updateNodeData.EditStringValue()) {
			logs.Error("修改组件名失败, 组件名非法, 只能使用英文字母 + 数字 + 下划线:", updateNodeData.EditStringValue())
			message.Info("修改组件名失败", "组件名 ["+updateNodeData.EditStringValue()+"] \n只能使用英文字母 + 数字 + 下划线", 200, 100)
			return err.RsFail
		}
		// 在当前设计面板只有唯一一个组件的名
		if m.FormTab.IsDuplicateName(m, updateNodeData.EditStringValue()) {
			logs.Error("修改组件名失败, 该组件名已存在", updateNodeData.EditStringValue())
			message.Info("修改组件名失败", "组件名 ["+updateNodeData.EditStringValue()+"] 已存在", 200, 100)
			return err.RsDuplicateName
		}
		// 当前应用唯一窗体名
		if m.ComponentType == consts.CtForm {
			if designer.IsDuplicateName(m, updateNodeData.EditStringValue()) {
				logs.Error("修改组件名失败, 该组件名已存在", updateNodeData.EditStringValue())
				message.Info("修改组件名失败", "组件名 ["+updateNodeData.EditStringValue()+"] 已存在", 200, 100)
				return err.RsDuplicateName
			}
		}
	case "enabled", "visible":
		// 忽略调用API的属性
		return err.RsIgnoreProp
	case "autosize", "borderstyle", "borderstyletoformborderstyle", "left", "top", "align", "anchors":
		// 忽略调用API的属性
		// Form 组件
		if m.ComponentType == consts.CtForm {
			return err.RsIgnoreProp
		}
	}
	return err.RsSuccess
}

// 反射调用函数
type reflector struct {
	object        any                      // 真实对象
	data          *vtedit.TEditNodeData    // 对象绑定数据
	objectNonWrap *TNonVisualComponentWrap // 非可视化对象, 只当前对象为非可视化对象时使用
}

// 查找方法（包含匿名嵌套字段的方法）
func (m *reflector) findMethod(val reflect.Value, methodName string) reflect.Value {
	if !val.IsValid() {
		return reflect.Value{}
	}
	// 如果是指针，先解引用
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 先尝试在当前类型中查找方法
	method := val.MethodByName(methodName)
	if method.IsValid() {
		return method
	}

	// 如果当前类型没有，尝试指针接收者
	if val.CanAddr() {
		method = val.Addr().MethodByName(methodName)
		if method.IsValid() {
			return method
		}
	}

	// 在匿名嵌套字段中查找方法
	return m.findMethodInEmbeddedFields(val, methodName)
}

// 在匿名嵌套字段中递归查找方法
func (m *reflector) findMethodInEmbeddedFields(val reflect.Value, methodName string) reflect.Value {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// 检查是否是匿名嵌套字段（嵌入字段）
		if field.Anonymous {
			embeddedField := val.Field(i)
			// 递归在嵌套字段中查找
			method := m.findMethod(embeddedField, methodName)
			if method.IsValid() {
				return method
			}
		}
	}
	return reflect.Value{}
}

// 根据当前数据节点的类型，将编辑数据转换为对应类型的参数值并返回。
func (m *reflector) convertArgsValue() (args []any) {
	switch m.data.Type() {
	case consts.PdtText, consts.PdtUint16:
		// string
		args = append(args, m.data.EditNodeData.StringValue)
	case consts.PdtInt:
		// int
		args = append(args, m.data.EditNodeData.IntValue)
	case consts.PdtInt64:
		// int64
		args = append(args, int64(m.data.EditNodeData.IntValue))
	case consts.PdtFloat:
		// float
		args = append(args, m.data.EditNodeData.FloatValue)
	case consts.PdtCheckBox:
		// bool
		data := m.data.AffiliatedNode.ToGo()
		if pData := vtedit.GetPropertyNodeData(data.Parent); pData != nil && pData.Type() != consts.PdtClass {
			// 当节点不是 class 时才处理
			dataList := pData.EditNodeData.CheckBoxValue
			var vals []int32
			for _, item := range dataList {
				if item.Checked {
					if v := mapper.GetLCL(item.Name); v == nil {
						logs.Error("更新组件属性失败 TSet集合取types值不存在 常量名:", item.Name)
						return nil
					} else {
						if val, err := tool.StrToInt32(tool.IntToString(v)); err != nil {
							logs.Error("更新组件属性失败 TSet集合取types值转换错误 常量名:", item.Name, "err:", err.Error())
						} else {
							vals = append(vals, val)
						}
					}
				}
			}
			set := types.NewSet(vals...)
			args = append(args, set)
		} else {
			args = append(args, m.data.EditNodeData.Checked)
		}
	case consts.PdtCheckBoxList:
		// TSet 集合
		dataList := m.data.EditNodeData.CheckBoxValue
		set := types.NewSet()
		for _, item := range dataList {
			if item.Checked {
				if v := mapper.GetLCL(item.Name); v == nil {
					logs.Error("更新组件属性失败 TSet集合取types值不存在 常量名:", item.Name)
					return nil
				} else {
					if val, err := tool.StrToInt32(tool.IntToString(v)); err != nil {
						logs.Error("更新组件属性失败 TSet集合取types值转换错误 常量名:", item.Name, "err:", err.Error())
					} else {
						set = set.Include(val)
					}
				}
			}
		}
		args = append(args, set)
	case consts.PdtComboBox:
		// const
		args = append(args, m.data.EditNodeData.StringValue)
	case consts.PdtColorSelect:
		// uint32
		args = append(args, uint32(m.data.EditNodeData.IntValue))
	case consts.PdtMethod:
		// 绑定事件方法不更新 API
	default:
		logs.Error("更新组件属性失败 未实现的类型:", m.data.Type(), "属性名名称:", m.data.Name(), "类:", lcl.AsObject(m.object).ToString())
		return nil
	}
	return
}

func (m *reflector) findMethodName() string {
	var methodName string
	switch m.data.Type() {
	case consts.PdtCheckBox:
		node := m.data.AffiliatedNode.ToGo()
		if node == nil {
			logs.Error("查找方法错误, 节点对象转换Go对象失败, 属性:", m.data.Name(), "类:", lcl.AsObject(m.object).ToString())
			return ""
		}
		parentNode := node.Parent
		// 有父节点 PdtCheckBoxList
		if pData := vtedit.GetPropertyNodeData(parentNode); pData != nil {
			if pData.Type() == consts.PdtClass {
				//父节点是 class 时使用当前属性名
				methodName = m.data.Name()
			} else {
				//父节点不是 class 时使用父节点属性名, 此时应为 PdtCheckBoxList
				methodName = pData.Name()
			}
		} else {
			// 没有父节点, 使用当前属性名
			methodName = m.data.Name()
		}
	default:
		// 其它默认当前属性名
		methodName = m.data.Name()
	}
	// Setter
	methodName = methodNameToSet(methodName)
	return methodName
}

func (m *reflector) findObject() (object reflect.Value) {
	object = reflect.ValueOf(m.object)
	data := m.data

	switch data.Type() {
	case consts.PdtCheckBox:
		// checkbox 需要从父节点获得所属实际节点
		node := m.data.AffiliatedNode.ToGo()
		if node == nil {
			logs.Error("查找对象错误, 节点对象转换Go对象失败, 属性:", m.data.Name(), "类:", lcl.AsObject(m.object).ToString())
			return
		}
		parentNode := node.Parent
		if pData := vtedit.GetPropertyNodeData(parentNode); pData != nil {
			if pData.Type() == consts.PdtClass {
				// 父节点是 class 时, object 使用父节点对象
				// 在下面 paths 时获取
			} else {
				//父节点不是 class 时使用父节点, 此时应为 PdtCheckBoxList
				data = pData // 使用父节点
			}
		}
	}
	// 方法是用于遍历对象路径, 当当前节点具有父节点时且父节点为 class 时查找出对象路径(paths)
	// 找到所有对象路径(paths)后从顶层对象开始调用, 直到返回当前属性所在的对象
	// todo 1: 可能存在的问题, 某父对象不是class一定是错误的
	// todo 2: 当属性（对象方法）不正确时需要做特殊处理转换, 例如: Pen() >= PenToPen() 等等
	paths := data.Paths()
	if len(paths) > 0 {
		for i := len(paths) - 1; i >= 0; i-- {
			name := paths[i] // todo 2
			in := make([]reflect.Value, 0)
			method := m.findMethod(object, name)
			results := method.Call(in)
			// 当前属性的所属对象
			object = results[0]
		}
	}
	return
}

// 调用方法
func (m *reflector) callMethod() ([]any, error) {
	var (
		object     reflect.Value
		methodName string
	)
	if m.objectNonWrap != nil && tool.Equal(m.data.Name(), "Left", "Top") {
		// 非可视化组件, 在修改位置时 Left, Top 使用包裹对象
		object = reflect.ValueOf(m.objectNonWrap.icon)
	} else {
		object = m.findObject()
	}
	methodName = m.findMethodName()
	if methodName == "" {
		return nil, fmt.Errorf("方法名称未找到 属性: %v 类: %v", m.data.Name(), lcl.AsObject(m.object).ToString())
	}
	method := m.findMethod(object, methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("方法 %v 未找到 属性: %v 类: %v", methodName, m.data.Name(), lcl.AsObject(m.object).ToString())
	}

	args := m.convertArgsValue()

	mType := method.Type()
	if mType.NumIn() != len(args) {
		return nil, fmt.Errorf("参数数量不匹配 需要: %v 实际: %v 属性: %v 类: %v", mType.NumIn(), len(args), m.data.Name(), lcl.AsObject(m.object).ToString())
	}

	// 准备参数
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValue := reflect.ValueOf(arg)
		targetType := mType.In(i)
		// 类型不同尝试转换
		if !argValue.Type().AssignableTo(targetType) {
			// 转换参数类型
			if convertValue, err := m.convertArgsType(arg, targetType); err != nil {
				return nil, fmt.Errorf("转换参数失败, index: %v 值: %v 需要类型: %v", i, arg, targetType.Name())
			} else {
				in[i] = convertValue
			}
		} else {
			in[i] = argValue
		}
		//logs.Debug("reflector callMethod targetType:", targetType, targetType.String(), targetType.Name())
	}

	logs.Debug("调用方法开始 对象:", object.Type().Name(), "方法:", methodName, "参数值:", args)
	// 调用方法
	results := method.Call(in)

	// 转换结果
	out := make([]any, len(results))
	for i, result := range results {
		out[i] = result.Interface()
	}
	logs.Debug("调用方法结束 对象:", object.Type().Name(), "方法:", methodName, "返回值:", out)
	return out, nil
}

// 将给定的值转换为目标类型
//
//	value: 需要转换的源值
//	targetType: 目标类型
func (m *reflector) convertArgsType(value any, targetType reflect.Type) (reflect.Value, error) {
	sourceValue := reflect.ValueOf(value)
	sourceType := sourceValue.Type()
	if sourceType.AssignableTo(targetType) {
		return sourceValue, nil
	}
	if sourceType.ConvertibleTo(targetType) {
		return sourceValue.Convert(targetType), nil
	}
	switch value.(type) {
	case string:
		val := mapper.GetLCL(value.(string))
		if val != nil {
			return reflect.ValueOf(val).Convert(targetType), nil
		}
	}
	if v, err := tool.ValueToTargetType(value, targetType); err == nil {
		return reflect.ValueOf(v), nil
	} else {
		return reflect.Value{}, err
	}
}
