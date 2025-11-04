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
