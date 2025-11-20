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

package logs

import (
	"github.com/energye/lcl/api"
	"log"
)

// 简单的日志输出

// 定义日志级别
const (
	LevelDebug = iota // 0：调试信息
	LevelInfo         // 1：普通信息
	LevelWarn         // 2：警告
	LevelError        // 3：错误
)

var Level = LevelInfo // 例如：设置为 INFO，只输出 INFO 及以上级别

func Debug(v ...any) {
	if Level <= LevelDebug {
		s := []any{"TID:", GetTID(), "[DEBUG]"}
		s = append(s, v...)
		log.Println(s...)
	}
}

func Info(v ...any) {
	if Level <= LevelInfo {
		s := []any{"TID:", GetTID(), "[INFO]"}
		s = append(s, v...)
		log.Println(s...)
	}
}

func Println(v ...any) {
	log.Println(v...)
}

func Warn(v ...any) {
	if Level <= LevelWarn {
		s := []any{"TID:", GetTID(), "[WARN]"}
		s = append(s, v...)
		log.Println(s...)
	}
}

func Error(v ...any) {
	if Level <= LevelError {
		s := []any{"TID:", GetTID(), "[ERROR]"}
		s = append(s, v...)
		log.Println(s...)
	}
}

func GetTID() uintptr {
	return api.CurrentThreadId()
}
