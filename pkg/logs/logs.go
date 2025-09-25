package logs

import (
	"log"
)

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
		s := []any{"[DEBUG]"}
		s = append(s, v...)
		log.Println(s...)
	}
}

func Info(v ...any) {
	if Level <= LevelInfo {
		s := []any{"[INFO]"}
		s = append(s, v...)
		log.Println(s...)
	}
}

func Warn(v ...any) {
	if Level <= LevelWarn {
		s := []any{"[WARN]"}
		s = append(s, v...)
		log.Println(s...)
	}
}

func Error(v ...any) {
	if Level <= LevelError {
		s := []any{"[ERROR]"}
		s = append(s, v...)
		log.Println(s...)
	}
}
