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

func Debug(v ...interface{}) {
	if Level <= LevelDebug {
		log.Println("[DEBUG]", v)
	}
}

func Info(v ...interface{}) {
	if Level <= LevelInfo {
		log.Println("[INFO]", v)
	}
}

func Warn(v ...interface{}) {
	if Level <= LevelWarn {
		log.Println("[WARN]", v)
	}
}

func Error(v ...interface{}) {
	if Level <= LevelError {
		log.Println("[ERROR]", v)
	}
}
