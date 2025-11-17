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
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// IntToString 数字转字符串
func IntToString(value any) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case uintptr:
		return strconv.FormatUint(uint64(v), 10)
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

// SetToHashSet
// 该函数将Set类型转换为HashSet
func SetToHashSet(value any) *HashSet {
	hashSet := NewHashSet()
	set := strings.Split(SetToString(value), ",")
	for _, v := range set {
		if v != "" {
			hashSet.Add(v)
		}
	}
	return hashSet
}

// StrToBool 字符串转 bool "true"/"false"
func StrToBool(s string) (bool, error) {
	val, err := strconv.ParseBool(s)
	return val, err
}

// StrToInt 字符串转 int
func StrToInt(s string) (int, error) {
	val, err := strconv.ParseInt(s, 0, strconv.IntSize)
	return int(val), err
}

// StrToInt8 字符串转 int8（范围：-128 ~ 127）
func StrToInt8(s string) (int8, error) {
	val, err := strconv.ParseInt(s, 0, 8)
	return int8(val), err
}

// StrToInt16 字符串转 int16（范围：-32768 ~ 32767）
func StrToInt16(s string) (int16, error) {
	val, err := strconv.ParseInt(s, 0, 16)
	return int16(val), err
}

// StrToInt32 字符串转 int32（范围：-2147483648 ~ 2147483647）
func StrToInt32(s string) (int32, error) {
	val, err := strconv.ParseInt(s, 0, 32)
	return int32(val), err
}

// StrToInt64 字符串转 int64（范围：-9223372036854775808 ~ 9223372036854775807）
func StrToInt64(s string) (int64, error) {
	val, err := strconv.ParseInt(s, 0, 64)
	return val, err
}

// StrToUint 字符串转 uint
func StrToUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 0, strconv.IntSize)
	return uint(val), err
}

// StrToUint8 字符串转 uint8（范围：0 ~ 255）
func StrToUint8(s string) (uint8, error) {
	val, err := strconv.ParseUint(s, 0, 8)
	return uint8(val), err
}

// StrToUint16 字符串转 uint16（范围：0 ~ 65535）
func StrToUint16(s string) (uint16, error) {
	val, err := strconv.ParseUint(s, 0, 16)
	if err != nil && len(s) == 1 {
		return uint16(s[0]), nil
	}
	return uint16(val), err
}

// StrToUint32 字符串转 uint32（范围：0 ~ 4294967295）
func StrToUint32(s string) (uint32, error) {
	val, err := strconv.ParseUint(s, 0, 32)
	return uint32(val), err
}

// StrToUint64 字符串转 uint64（范围：0 ~ 18446744073709551615）
func StrToUint64(s string) (uint64, error) {
	val, err := strconv.ParseUint(s, 0, 64)
	return val, err
}

// StrToUintptr 字符串转 uintptr
func StrToUintptr(s string) (uintptr, error) {
	val, err := strconv.ParseUint(s, 0, strconv.IntSize)
	return uintptr(val), err
}

// StrToFloat32 字符串转 float32（范围：±3.4e38，精度约 6-7 位小数）
func StrToFloat32(s string) (float32, error) {
	val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	if val > math.MaxFloat32 {
		return 0, errors.New("value exceeds float32 maximum range")
	}
	return float32(val), nil
}

// StrToFloat64 字符串转 float64（范围：±1.8e308，精度约 15-17 位小数）
func StrToFloat64(s string) (float64, error) {
	val, err := strconv.ParseFloat(s, 64)
	return val, err
}

// ValueToTargetType 基本类型转换函数, 将值转换为目标类型
func ValueToTargetType(value any, targetType reflect.Type) (any, error) {
	sourceType := reflect.TypeOf(value)
	if sourceType == targetType {
		return value, nil
	}
	if value == nil {
		return getZeroValue(targetType), nil
	}
	strValue, err := toString(value)
	if err != nil {
		return nil, fmt.Errorf("ValueToTargetType 转换到 string 类型失败: %v", err)
	}
	return fromString(strValue, targetType)
}

// toString 将任意基本类型转换为字符串
func toString(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32, float64:
		return fmt.Sprintf("%f", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	default:
		return "", fmt.Errorf("unsupported source type: %T", value)
	}
}

// fromString 将字符串转换为目标类型
func fromString(strValue string, targetType reflect.Type) (any, error) {
	switch targetType.Kind() {
	case reflect.String:
		return strValue, nil
	case reflect.Uintptr:
		return StrToUintptr(strValue)
	case reflect.Int:
		return StrToInt(strValue)
	case reflect.Int8:
		return StrToInt8(strValue)
	case reflect.Int16:
		return StrToInt16(strValue)
	case reflect.Int32:
		return StrToInt32(strValue)
	case reflect.Int64:
		return StrToInt64(strValue)
	case reflect.Uint:
		return StrToUint(strValue)
	case reflect.Uint8:
		return StrToUint8(strValue)
	case reflect.Uint16:
		return StrToUint16(strValue)
	case reflect.Uint32:
		return StrToUint32(strValue)
	case reflect.Uint64:
		return StrToUint64(strValue)
	case reflect.Float32:
		return StrToFloat32(strValue)
	case reflect.Float64:
		return StrToFloat64(strValue)
	case reflect.Bool:
		return StrToBool(strValue)
	default:
		return nil, fmt.Errorf("不支持的目标类型转换: %v", targetType)
	}
}

// getZeroValue 获取目标类型的零值
func getZeroValue(targetType reflect.Type) any {
	switch targetType.Kind() {
	case reflect.String:
		return ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint(0)
	case reflect.Float32, reflect.Float64:
		return 0.0
	case reflect.Bool:
		return false
	default:
		return nil
	}
}
