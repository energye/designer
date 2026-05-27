# 日志

日志模块提供异步、带缓冲的日志系统，支持多种日志级别和全局/实例两种使用方式。

## 创建日志实例

```go
import "github.com/energye/energy/v3/logger"

log := logger.New(logger.Config{
    Level:      logger.InfoLevel,
    Output:     os.Stdout,
    Caller:     true,
    BufferSize: 1024,
    BatchSize:  100,
})
```

### Config 配置

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| Level | Level | InfoLevel | 最低日志级别 |
| Output | io.Writer | os.Stdout | 日志输出目标 |
| Caller | bool | false | 是否显示调用位置 |
| BufferSize | int | 65536 | 日志通道缓冲大小 |
| BatchSize | int | 256 | 批量写入大小 |

## 日志级别

| 级别 | 常量 | 说明 |
|------|------|------|
| Debug | logger.DebugLevel | 调试信息 |
| Info | logger.InfoLevel | 普通信息 |
| Warn | logger.WarnLevel | 警告信息 |
| Error | logger.ErrorLevel | 错误信息 |

## 使用方式

### 实例方法

```go
log := logger.New(logger.Config{Level: logger.DebugLevel})

log.Debug("调试信息", "key", value)
log.Info("普通信息", "key", value)
log.Warn("警告信息", "key", value)
log.Error("错误信息", "key", value)
```

### 全局日志

```go
// 设置全局日志实例
logger.SetDefault(logger.New(logger.Config{
    Level: logger.InfoLevel,
}))

// 使用全局日志
logger.Debug("调试信息", "key", value)
logger.Info("普通信息", "key", value)
logger.Warn("警告信息", "key", value)
logger.Error("错误信息", "key", value)
```

### 获取全局日志实例

```go
log := logger.L()
log.Info("通过 L() 获取的实例")
```

### 设置日志级别

```go
log.SetLevel(logger.DebugLevel)
```

### 关闭日志

```go
log.Close() // 关闭日志，等待缓冲区写入完成
```

## 异步写入

日志采用异步通道写入机制：

1. 日志写入通道（非阻塞）
2. 后台协程批量读取
3. 批量写入输出目标

### 优势

- 不阻塞主协程
- 批量写入减少 I/O 次数
- 通道缓冲处理突发日志

### 配置建议

| 场景 | BufferSize | BatchSize |
|------|------------|-----------|
| 开发调试 | 1024 | 32 |
| 生产环境 | 65536 | 256 |
| 高并发 | 131072 | 512 |

## 调用位置

启用 `Caller: true` 后，日志会显示调用位置：

```
2024-01-15 10:30:45 INFO main.go:25 用户登录成功
```

格式：`日期 时间 级别 文件名:行号 消息`

## 自定义输出

### 输出到文件

```go
file, err := os.Create("app.log")
if err != nil {
    panic(err)
}
defer file.Close()

log := logger.New(logger.Config{
    Level:  logger.InfoLevel,
    Output: file,
})
```

### 输出到多目标

```go
writer := io.MultiWriter(os.Stdout, file)
log := logger.New(logger.Config{
    Level:  logger.InfoLevel,
    Output: writer,
})
```

## 调试模式

使用 `api.SetDebug(true)` 启用调试模式，会输出更详细的框架内部日志：

```go
import "github.com/energye/lcl/api"

func main() {
    api.SetDebug(true) // 启用调试
    // ...
}
```

## 完整示例

```go
package main

import (
    "io"
    "os"
    "github.com/energye/energy/v3/logger"
)

func main() {
    // 创建文件日志
    file, _ := os.Create("app.log")
    defer file.Close()

    // 初始化全局日志
    logger.SetDefault(logger.New(logger.Config{
        Level:  logger.DebugLevel,
        Output: io.MultiWriter(os.Stdout, file),
        Caller: true,
    }))

    logger.Info("应用启动")
    logger.Debug("配置加载完成", "count", 10)

    // 业务逻辑
    if err := doSomething(); err != nil {
        logger.Error("操作失败", "error", err)
    }

    logger.Info("应用退出")
}
```
