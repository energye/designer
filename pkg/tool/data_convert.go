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
	"fmt"
	"strconv"
	"strings"
)

// IntToString 数字转字符串
func IntToString(value any) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// BoolToString 布尔转字符串
func BoolToString(value any) string {
	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v)
	case string:
		// 如果已经是字符串，直接返回
		return v
	default:
		// 其他类型尝试转换为字符串
		return fmt.Sprintf("%v", v)
	}
}

// FloatToString 浮点转字符串
func FloatToString(value any) string {
	switch v := value.(type) {
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// SetToString 集合转字符串
// [Xxx,Xxx,Xxx] > Xxx,Xxx,Xxx
func SetToString(value any) string {
	val := fmt.Sprintf("%v", value)
	return strings.Trim(val, "[]")
}
