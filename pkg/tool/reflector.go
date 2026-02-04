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
	"github.com/energye/designer/consts"
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
func methodNames(t reflect.Type) *ArrayMap[string, *TMethod] {
	methods := NewArrayMap[string, *TMethod]()
	collectMethods(t, 0, methods)
	if t.Kind() != reflect.Ptr {
		ptrType := reflect.PtrTo(t)
		collectMethods(ptrType, 0, methods)
	}
	return methods
}

// collectMethods 递归收集方法
func collectMethods(t reflect.Type, level int, methods *ArrayMap[string, *TMethod]) {
	originalType := t
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
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
func GetObjectMethodNames(v any) *ArrayMap[string, *TMethod] {
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

// FixPropInfo 修复属性信息，将属性名与对象方法名进行匹配和修正
// @param methods 对象方法列表，用于查找和匹配属性对应的方法
// @param prop 组件属性信息，需要被修正的属性对象
func FixPropInfo(methods *ArrayMap[string, *TMethod], prop *lcl.ComponentProperties) {
	if methods == nil {
		return
	}
	name := strings.ToLower(prop.Name)
	if methods.ContainsKey(name) {
		// 当前属性名不存在于对象的方法列表中
		// 原因: 1. 完全不存在, 2. 属性名与对象方法名不一致
		// 当为原因2时需要将属性名改为实际的方法名
		prop.Name = methods.Get(name).Name
		return
	} else if prop.Kind == consts.TkMethod {
		// 属性类型为事件
		onName := "set" + name
		if methods.ContainsKey(onName) {
			// SetOnXxx
			onName = methods.Get(onName).Name
			// 去除 Set 前缀
			onName = onName[3:]
			prop.Name = onName
			return
		} else {
			logs.Warn("属性类型为事件, 但对象方法不存在:", onName)
			return
		}
	}
	//logs.Warn("属性和对象方法不匹配, 当前属性名:", prop.Name, "属性类型:", prop.Type)
	// 遍历对象方法列表, 匹配出所有属性名
	type_ := strings.ToLower(RemoveT(prop.Type))
	name_ := strings.ToLower(name)
	// 最后的尝试, 当Go里的对象方法名, 同时包含了组件属性名+类型时,
	// 这种情况属于同名不同参的方法, 因为Go不支持重载方法
	// 例如： BorderStyleToFormStyle
	methods.Iterate(func(methodName string, value *TMethod) bool {
		if strings.Contains(methodName, name_) && strings.Contains(methodName, type_) {
			prop.Name = value.Name
			return true
		}
		return false
	})
}
