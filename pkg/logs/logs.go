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
	"github.com/energye/energy/v3/application/pack"
	"github.com/energye/energy/v3/logger"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

// Window 当前用于调试窗口指针
var Window lcl.IForm

var (
	myLogger *logger.Logger
	LogPath  string
)

func init() {
	getLogFilePath := func() string {
		homeDir, _ := os.UserHomeDir()
		bundleId := pack.Info.Id
		if bundleId == "" {
			bundleId = "com.energye.designer.dev" // dev
		}
		if runtime.GOOS == "darwin" {
			logDir := filepath.Join(homeDir, "Library", "Logs", bundleId)
			_ = os.MkdirAll(logDir, 0700)
			return filepath.Join(logDir, "designer.log")
		} else {
			return filepath.Join(homeDir, ".energy", "designer.log")
		}
	}
	LogPath = getLogFilePath()
	file, err := os.OpenFile(LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		myLogger = logger.New(logger.Config{Level: logger.DebugLevel, Output: os.Stdout})
	} else {
		multiWriter := io.MultiWriter(os.Stdout, file)
		myLogger = logger.New(logger.Config{
			Level:  logger.DebugLevel,
			Output: multiWriter,
		})
	}
	api.SetOnReleaseCallback(func() {
		if file != nil {
			_ = file.Close()
		}
		myLogger.Close()
	})
}

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
		s := []any{GetTID(), windowIsValid()}
		s = append(s, v...)
		myLogger.Debug("-", s...)
	}
}

func Info(v ...any) {
	if Level <= LevelInfo {
		s := []any{GetTID(), windowIsValid()}
		s = append(s, v...)
		myLogger.Info("-", s...)
	}
}

func Println(v ...any) {
	log.Println(v...)
}

func Warn(v ...any) {
	if Level <= LevelWarn {
		s := []any{GetTID(), windowIsValid()}
		s = append(s, v...)
		myLogger.Warn("-", s...)
	}
}

func Error(v ...any) {
	if Level <= LevelError {
		s := []any{GetTID(), windowIsValid()}
		s = append(s, v...)
		myLogger.Error("-", s...)
	}
}

func GetTID() string {
	if api.LibLoaded() {
		return "[" + strconv.FormatUint(uint64(api.CurrentThreadId()), 10) + "]"
	}
	return "[0]"
}

func windowIsValid() string {
	if Window != nil && Window.IsValid() {
		return "window_valid=true"
	}
	return ""
}
