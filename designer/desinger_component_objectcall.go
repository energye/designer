package designer

import (
	"errors"
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/mapper"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
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
	data := nodeData.EditNodeData
	m.drag.Hide()
	lcl.RunOnMainThreadAsync(func(id uint32) {
		reflector := &embeddingReflector{object: m.originObject, data: data}
		result, err := reflector.CallMethod()
		_ = result
		if err != nil {
			logs.Error("更新组件属性失败,", err.Error())
		}
		m.drag.Show()
	})
}

type embeddingReflector struct {
	object any
	data   *vtedit.TEditLinkNodeData
}

// 查找方法（包含匿名嵌套字段的方法）
func (m *embeddingReflector) findMethod(val reflect.Value, methodName string) reflect.Value {
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
func (m *embeddingReflector) findMethodInEmbeddedFields(val reflect.Value, methodName string) reflect.Value {
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

func (m *embeddingReflector) convertArgs() (args []any) {
	switch m.data.Type {
	case vtedit.PdtText:
		// string
		args = append(args, m.data.StringValue)
	case vtedit.PdtInt, vtedit.PdtInt64:
		// int
		args = append(args, m.data.IntValue)
	case vtedit.PdtFloat:
		// float
		args = append(args, m.data.FloatValue)
	case vtedit.PdtCheckBox:
		// bool
		args = append(args, m.data.Checked)
	case vtedit.PdtCheckBoxList:
		// TSet
	case vtedit.PdtComboBox:
		// const
		args = append(args, m.data.StringValue)
	case vtedit.PdtColorSelect:
		// uint32
	default:
		logs.Error("更新组件属性失败, 未实现的类型:", m.data.Type)
		return nil
	}
	return
}

// 调用方法
func (m *embeddingReflector) CallMethod() ([]any, error) {
	methodName := m.data.Name
	methodName = methodNameToSet(methodName)

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
		if !argValue.Type().AssignableTo(targetType) {
			if convertValue, err := m.convertArgsType(arg, targetType); err != nil {
				return nil, fmt.Errorf("转换参数失败, index: %v 值: %v 需要类型: %v", i, arg, targetType.Name())
			} else {
				in[i] = convertValue
			}
		} else {
			in[i] = argValue
		}
		//fmt.Println("targetType:", targetType, targetType.String(), targetType.Name())
	}

	// 调用方法
	results := method.Call(in)

	// 转换结果
	out := make([]any, len(results))
	for i, result := range results {
		out[i] = result.Interface()
	}

	return out, nil
}

func (m *embeddingReflector) convertArgsType(value any, targetType reflect.Type) (reflect.Value, error) {
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
			return reflect.ValueOf(val), nil
		}
	}
	return reflect.Value{}, errors.New("参数类型转换失败")
}
