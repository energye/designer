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

// collectMethodNames 递归收集类型t及其嵌套结构体的所有方法名
func methodNames(t reflect.Type, processed map[reflect.Type]bool) *ArrayMap[string] {
	// 避免重复处理同一类型
	if processed[t] {
		return nil
	}
	processed[t] = true
	methods := NewArrayMap[string]()
	// 若为指针类型，先处理其指向的元素类型
	if t.Kind() == reflect.Ptr {
		elem := t.Elem()
		subMethods := methodNames(elem, processed)
		if subMethods != nil {
			subMethods.Iterate(func(key string, value string) bool {
				methods.Add(strings.ToLower(key), value)
				return false
			})
		}
	}
	// 收集当前类型的所有导出方法
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		methods.Add(strings.ToLower(method.Name), method.Name)
	}

	// 若为结构体类型，处理嵌套的匿名结构体字段
	if t.Kind() == reflect.Struct {
		// 遍历所有字段，处理匿名嵌套结构体
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			// 匿名字段
			if field.Anonymous {
				fieldType := field.Type
				subMethods := methodNames(fieldType, processed)
				if subMethods != nil {
					subMethods.Iterate(func(key string, value string) bool {
						methods.Add(strings.ToLower(key), value)
						return false
					})
				}
			}
		}

		// 处理结构体的指针类型（*T），以收集指针接收器的方法
		ptrType := reflect.PtrTo(t)
		subMethods := methodNames(ptrType, processed)
		if subMethods != nil {
			subMethods.Iterate(func(key string, value string) bool {
				methods.Add(strings.ToLower(key), value)
				return false
			})
		}
	}
	return methods
}

// GetObjectMethodNames 获取指定对象的所有方法名称
// 该函数通过反射机制遍历对象类型及其嵌套类型的所有方法，并返回一个包含所有方法名称的ArrayMap
// 参数:
//   v: 任意类型的对象实例，用于获取其方法信息
// 返回值:
//   *ArrayMap[string]: 包含所有方法名称的ArrayMap指针，如果输入为nil或无法获取类型信息则返回nil

func GetObjectMethodNames(v any) *ArrayMap[string] {
	if v == nil {
		return nil
	}
	t := reflect.TypeOf(v)
	if t == nil {
		return nil
	}
	processed := make(map[reflect.Type]bool)
	methods := methodNames(t, processed)
	return methods
}
