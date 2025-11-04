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
	"reflect"
	"strings"
)

// TMethod 存储方法信息
type TMethod struct {
	Name       string
	StructName string
	Level      int
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
