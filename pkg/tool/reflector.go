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

package tool

import (
	"gen/tool"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"reflect"
	"strings"
)

// TMethod 存储方法信息
type TMethod struct {
	Name       string // 方法名
	StructName string // 所属结构
	Level      int    // 等级, 0~[1...] 数字越小表示优先级高
}

// GetAllMethods 获取结构体及其嵌入结构体的所有导出方法
func methodNames(t reflect.Type) *ArrayMap[*TMethod] {
	methods := NewArrayMap[*TMethod]()
	collectMethods(t, 0, methods)
	if t.Kind() != reflect.Ptr {
		ptrType := reflect.PtrTo(t)
		collectMethods(ptrType, 0, methods)
	}
	return methods
}

// collectMethods 递归收集方法
func collectMethods(t reflect.Type, level int, methods *ArrayMap[*TMethod]) {
	originalType := t
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	structName := t.Name()
	// 收集当前类型的Method
	numMethod := originalType.NumMethod()
	for i := 0; i < numMethod; i++ {
		method := originalType.Method(i)
		// 只处理导出的方法并且避免重复
		methods.Add(
			strings.ToLower(method.Name),
			&TMethod{Name: method.Name, StructName: structName, Level: level},
		)
	}
	// 遍历嵌入字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			// 指针类型的嵌入字段
			if field.Type.Kind() != reflect.Ptr {
				ptrFieldType := reflect.PtrTo(field.Type)
				collectMethods(ptrFieldType, level+1, methods)
			} else {
				collectMethods(field.Type, level+1, methods)
			}
		}
	}
}

// GetObjectMethodNames 获取指定对象的所有方法名称
// 该函数通过反射机制遍历对象类型及其嵌套类型的所有方法，并返回一个包含所有方法名称的ArrayMap
// 参数:
//
//	v: 任意类型的对象实例，用于获取其方法信息
//
// 返回值:
//
//	*ArrayMap[string]: 包含所有方法名称的ArrayMap指针，如果输入为nil或无法获取类型信息则返回nil
func GetObjectMethodNames(v any) *ArrayMap[*TMethod] {
	if v == nil {
		return nil
	}
	t := reflect.TypeOf(v)
	if t == nil {
		return nil
	}
	methods := methodNames(t)
	return methods
}

// 修复属性信息
func FixPropInfo(methods *ArrayMap[*TMethod], prop *lcl.ComponentProperties) {
	if methods == nil {
		return
	}
	name := strings.ToLower(prop.Name)
	if methods.ContainsKey(name) {
		// 当前属性名不存在于对象的方法列表中
		// 原因: 1. 完全不存在, 2. 属性名与对象方法名不一致
		// 当为原因2时需要将属性名改为实际的方法名
		prop.Name = methods.Get(name).Name
	} else if prop.Kind == "tkMethod" {
		// TODO on event
	} else {
		logs.Warn("属性和对象方法不匹配, 当前属性名:", prop.Name, "属性类型:", prop.Type)
		// 遍历对象方法列表, 匹配出所有属性名
		type_ := strings.ToLower(tool.RemoveT(prop.Type))
		name_ := strings.ToLower(name)
		methods.Iterate(func(methodName string, value *TMethod) bool {
			if strings.Contains(methodName, name_) && strings.Contains(methodName, type_) {
				prop.Name = value.Name
				return true
			}
			return false
		})
	}
	// 获取属性默认值
}
