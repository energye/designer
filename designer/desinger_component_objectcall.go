package designer

import (
	"errors"
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/mapper"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"reflect"
)

// 组件对象函数调用

func methodNameToSet(name string) string {
	name = tool.FirstToUpper(name)
	return "Set" + name
}

// 更新当前组件属性
func (m *DesigningComponent) UpdateComponentProperty(nodeData *vtedit.TEditNodeData) {
	logs.Debug("更新组件:", m.object.ToString(), "属性:", nodeData.EditNodeData.Name)
	m.drag.Hide()
	lcl.RunOnMainThreadAsync(func(id uint32) {
		ref := &reflector{object: m.originObject, data: nodeData}
		result, err := ref.callMethod()
		_ = result
		if err != nil {
			logs.Error("更新组件属性失败,", err.Error())
		}
		m.drag.Show()
	})
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

func (m *reflector) convertArgs() (args []any) {
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
		args = append(args, m.data.EditNodeData.Checked)
	case vtedit.PdtCheckBoxList:
		// TSet 集合
		dataList := m.data.EditNodeData.CheckBoxValue
		set := types.NewSet()
		for _, item := range dataList {
			if item.Checked {
				if v := mapper.GetLCL(item.Name); v == nil {
					logs.Error("更新组件属性失败, TSet集合取types值不存在 常量名:", item.Name)
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
		logs.Error("更新组件属性失败, 未实现的类型:", m.data.EditNodeData.Type)
		return nil
	}
	return
}

func (m *reflector) methodName() string {
	var methodName string
	switch m.data.EditNodeData.Type {
	case vtedit.PdtCheckBox:
		node := m.data.AffiliatedNode.ToGo()
		parentNode := node.Parent
		// 有父节点 PdtCheckBoxList
		if pData := vtedit.GetPropertyNodeData(parentNode); pData != nil {
			methodName = pData.EditNodeData.Name
		} else {
			methodName = m.data.EditNodeData.Name
		}
	default:
		methodName = m.data.EditNodeData.Name
	}
	methodName = methodNameToSet(methodName)
	return methodName
}

// 调用方法
func (m *reflector) callMethod() ([]any, error) {
	methodName := m.methodName()

	val := reflect.ValueOf(m.object)

	method := m.findMethod(val, methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("方法 %v 未找到", methodName)
	}

	args := m.convertArgs()

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

	logs.Debug("调用方法:", methodName, "参数:", args)

	// 调用方法
	results := method.Call(in)

	// 转换结果
	out := make([]any, len(results))
	for i, result := range results {
		out[i] = result.Interface()
	}

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
