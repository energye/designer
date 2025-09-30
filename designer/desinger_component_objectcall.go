package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/vtedit"
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
	methodName := nodeData.EditNodeData.Name
	methodName = methodNameToSet(methodName)
	result, err := gEmbeddingReflector.CallMethod(m.originObject, methodName)
	if err != nil {
		logs.Error("更新组件属性失败,", err.Error())
	}
	fmt.Println("result:", result)
}

type embeddingReflector struct{}

var gEmbeddingReflector = &embeddingReflector{}

// 查找方法（包含匿名嵌套字段的方法）
func (e *embeddingReflector) findMethod(val reflect.Value, methodName string) reflect.Value {
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
	return e.findMethodInEmbeddedFields(val, methodName)
}

// 在匿名嵌套字段中递归查找方法
func (e *embeddingReflector) findMethodInEmbeddedFields(val reflect.Value, methodName string) reflect.Value {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// 检查是否是匿名嵌套字段（嵌入字段）
		if field.Anonymous {
			embeddedField := val.Field(i)
			// 递归在嵌套字段中查找
			method := e.findMethod(embeddedField, methodName)
			if method.IsValid() {
				return method
			}
		}
	}
	return reflect.Value{}
}

// 调用方法
func (e *embeddingReflector) CallMethod(obj any, methodName string, args ...any) ([]any, error) {
	val := reflect.ValueOf(obj)

	method := e.findMethod(val, methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("方法 %v 未找到", methodName)
	}

	mType := method.Type()
	if mType.NumIn() != len(args) {
		return nil, fmt.Errorf("参数数量不匹配 需要: %v 实际: %v", mType.NumIn(), len(args))
	}
	// 转换参数类型

	// 准备参数
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
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
