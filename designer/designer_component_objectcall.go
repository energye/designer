package designer

import (
	"errors"
	"fmt"
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/mapper"
	"github.com/energye/designer/pkg/message"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/types"
	"reflect"
	"strings"
)

// 组件对象函数调用

func methodNameToSet(name string) string {
	name = tool.FirstToUpper(name)
	return "Set" + name
}

// 更新组件属性到对象
func (m *DesigningComponent) UpdateComponentPropertyToObject(updateNodeData *vtedit.TEditNodeData) {
	m.drag.Hide()
	defer m.drag.Show()
	logs.Debug("更新组件:", m.ClassName(), "属性:", updateNodeData.EditNodeData.Name)
	// 检查当前组件属性是否允许更新
	if rs := m.CheckCanUpdateProp(updateNodeData); rs == err.RsSuccess {
		logs.Info("检查允许更新属性, 该属性", updateNodeData.EditNodeData.Name, "调用 API 更新, 同时更新节点数据")
		ref := &reflector{object: m.originObject, data: updateNodeData}
		result, err := ref.callMethod()
		_ = result
		if err != nil {
			logs.Error("调用 API 更新组件属性失败", err.Error())
		} else {
			logs.Info("调用 API 更新组件属性成功, 更新节点数据")
			m.UpdateTreeNode(updateNodeData)
		}
	} else if rs == err.RsIgnoreProp {
		logs.Info("检查允许更新属性, 该属性", updateNodeData.EditNodeData.Name, "忽略 API 更新, 只更新节点数据")
		m.UpdateTreeNode(updateNodeData)
	} else {
		// 更新失败
		switch rs {
		case err.RsDuplicateName: // 重复的组件名
			logs.Error("重复的组件名 检查允许更新属性失败, RS:", rs, "恢复节点内的组件名")
			// 恢复节点内的组件名
			updateNodeData.SetEditValue(m.Name())
			inspector.componentProperty.propertyTree.InvalidateNode(updateNodeData.AffiliatedNode)
		default:
			logs.Error("重复的组件名 检查允许更新属性失败, RS:", rs)
		}
	}
}

// 更新组件树节点信息
// 在设计组件属性修改后同步修改组件树节点可见值
func (m *DesigningComponent) UpdateTreeNode(updateNodeData *vtedit.TEditNodeData) {
	if !m.node.IsValid() {
		logs.Error("更新组件树失败, 当前设计组件节点无效")
		return
	}
	data := updateNodeData.EditNodeData
	propName := strings.ToLower(data.Name)
	logs.Debug("更新组件树, 尝试更新属性:", data.Name)
	switch propName {
	case "name":
		m.node.SetText(m.TreeName())
		if m.componentType == CtForm {
			m.ownerFormTab.sheet.SetCaption(m.Name())
		}
	}
}

// 检查是否允许更新属性
func (m *DesigningComponent) CheckCanUpdateProp(updateNodeData *vtedit.TEditNodeData) err.ResultStatus {
	if !m.node.IsValid() {
		// 无效节点对象
		return err.RsNotValid
	}
	data := updateNodeData.EditNodeData
	propName := strings.ToLower(data.Name)
	switch propName {
	case "name":
		// 在当前设计面板只有唯一一个组件的名
		if m.ownerFormTab.IsDuplicateName(m, data.EditValue()) {
			logs.Error("修改组件名失败, 该组件名已存在", data.EditValue())
			message.Info("修改组件名失败", "组件名 ["+data.EditValue()+"] 已存在", 200, 100)
			return err.RsDuplicateName
		}
	case "enabled", "visible":
		return err.RsIgnoreProp

	}
	return err.RsSuccess
}

// 反射调用函数
type reflector struct {
	object any
	data   *vtedit.TEditNodeData
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

func (m *reflector) convertArgsValue() (args []any) {
	switch m.data.EditNodeData.Type {
	case vtedit.PdtText:
		// string
		args = append(args, m.data.EditNodeData.StringValue)
	case vtedit.PdtInt:
		// int
		args = append(args, m.data.EditNodeData.IntValue)
	case vtedit.PdtInt64:
		// int64
		args = append(args, int64(m.data.EditNodeData.IntValue))
	case vtedit.PdtFloat:
		// float
		args = append(args, m.data.EditNodeData.FloatValue)
	case vtedit.PdtCheckBox:
		// bool
		data := m.data.AffiliatedNode.ToGo()
		if pData := vtedit.GetPropertyNodeData(data.Parent); pData != nil {
			dataList := pData.EditNodeData.CheckBoxValue
			var vals []int32
			for _, item := range dataList {
				if item.Checked {
					if v := mapper.GetLCL(item.Name); v == nil {
						logs.Error("[更新组件属性失败] TSet集合取types值不存在 常量名:", item.Name)
						return nil
					} else {
						vals = append(vals, v.(int32))
					}
				}
			}
			set := types.NewSet(vals...)
			args = append(args, set)
		} else {
			args = append(args, m.data.EditNodeData.Checked)
		}
	case vtedit.PdtCheckBoxList:
		// TSet 集合
		dataList := m.data.EditNodeData.CheckBoxValue
		set := types.NewSet()
		for _, item := range dataList {
			if item.Checked {
				if v := mapper.GetLCL(item.Name); v == nil {
					logs.Error("[更新组件属性失败] TSet集合取types值不存在 常量名:", item.Name)
					return nil
				} else {
					set = set.Include(v.(int32))
				}
			}
		}
		args = append(args, set)
	case vtedit.PdtComboBox:
		// const
		args = append(args, m.data.EditNodeData.StringValue)
	case vtedit.PdtColorSelect:
		// uint32
		args = append(args, uint32(m.data.EditNodeData.IntValue))
	default:
		logs.Error("[更新组件属性失败] 未实现的类型:", m.data.EditNodeData.Type)
		return nil
	}
	return
}

func (m *reflector) findMethodName() string {
	var methodName string
	switch m.data.Type() {
	case vtedit.PdtCheckBox:
		node := m.data.AffiliatedNode.ToGo()
		parentNode := node.Parent
		// 有父节点 PdtCheckBoxList
		if pData := vtedit.GetPropertyNodeData(parentNode); pData != nil {
			methodName = pData.Name()
		} else {
			methodName = m.data.Name()
		}
	default:
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
	case vtedit.PdtCheckBox:
		// checkbox 需要从父节点获得所属实际节点
		node := m.data.AffiliatedNode.ToGo()
		parentNode := node.Parent
		if pData := vtedit.GetPropertyNodeData(parentNode); pData != nil {
			data = pData // 使用父节点
		}
	}
	// 方法是用于遍历对象路径, 当当前节点具有父节点时且父节点为 class 时查找出所有对象目录
	// 找到所有对象目录后从顶层对象开始调用, 直到返回当前属性所在的对象
	// todo 1: 可能存在的问题, 某父对象不是class一定是错误的
	// todo 2: 当属性（对象方法）不正确时需要做特殊处理转换, 例如: Pen() >= PenToPen() 等等
	iterObjectName := func(data *vtedit.TEditNodeData) {
		var paths []string
		pData := data.Parent
		for pData != nil {
			if pData.Type() == vtedit.PdtClass { //todo 1
				paths = append(paths, pData.Name())
			} else {
				// 不正确, 直接退出
				panic("递归遍历属性节点对象路径错误, 对象非class类型, 节点必须为class类型: " + pData.Name())
			}
			pData = pData.Parent
		}
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
	}
	iterObjectName(data)
	return
}

// 调用方法
func (m *reflector) callMethod() ([]any, error) {
	object := m.findObject()
	methodName := m.findMethodName()

	method := m.findMethod(object, methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("方法 %v 未找到", methodName)
	}

	args := m.convertArgsValue()

	mType := method.Type()
	if mType.NumIn() != len(args) {
		return nil, fmt.Errorf("参数数量不匹配 需要: %v 实际: %v", mType.NumIn(), len(args))
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

	logs.Debug("调用方法开始:", methodName, "参数值:", args)
	// 调用方法
	results := method.Call(in)

	// 转换结果
	out := make([]any, len(results))
	for i, result := range results {
		out[i] = result.Interface()
	}
	logs.Debug("调用方法结束:", methodName, "返回值:", out)

	return out, nil
}

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
	return reflect.Value{}, errors.New("参数类型转换失败")
}
