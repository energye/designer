//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

package config

import (
	"bufio"
	"encoding/json"
	"github.com/energye/lcl/tool/command"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// 仅提取 Go 相关环境变量

var goEnvVars = map[string]bool{
	"GOROOT":      true,
	"GOPATH":      true,
	"GOBIN":       true,
	"GO111MODULE": true,
	"GONOSUMDB":   true,
	"GONOPROXY":   true,
	"GOPRIVATE":   true,
	"GOFLAGS":     true,
	"GOSUMDB":     true,
	"GOPROXY":     true,
	"GOMODCACHE":  true,
	"GOWORK":      true,
	"PATH":        true,
}

// InitGoEnv 初始化 Go 语言环境变量。
//
// 该函数主要在非 Windows 平台（macOS/Linux）执行，用于确保 GUI 应用程序能够正确读取
// 系统 shell 中配置的 Go 环境变量。它会先应用预定义的运行时环境变量，随后通过
// 执行 `go env -json` 命令获取完整的 Go 环境配置并缓存到全局变量 GGoEnv 中。
func InitGoEnv() {
	env := BuildRuntimeEnv()
	for name, value := range env {
		if goEnvVars[name] {
			println(name, "=", value)
			_ = os.Setenv(name, value)
		}
	}
	var (
		goEnvCMDOk = true
		goEnvLines []string
	)
	goEnvCMD := command.NewCMD()
	goEnvCMD.HideWindow = true
	goEnvCMD.Console = func(data string, level command.Level) {
		if level == command.LError {
			goEnvCMDOk = false
		} else {
			goEnvLines = append(goEnvLines, data)
		}
	}
	goEnvCMD.Command("go", "env", "-json")
	if goEnvCMDOk {
		goEnvData := strings.Join(goEnvLines, "")
		e := json.Unmarshal([]byte(goEnvData), &GGoEnv)
		if e != nil {
			println("Init go env err:", e.Error())
		}
	}
}

func BuildRuntimeEnv() map[string]string {
	env := envSliceToMap(os.Environ())
	userEnv := readGoEnvFromConfig()
	for k, v := range userEnv {
		if k == "PATH" {
			env["PATH"] = mergePath(v, env["PATH"])
		} else {
			env[k] = v
		}
	}
	return env
}

func readGoEnvFromConfig() map[string]string {
	result := make(map[string]string)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return result
	}
	files := shellConfigFiles(homeDir)
	for _, file := range files {
		if !exists(file) {
			continue
		}
		parseEnvFromShellFile(file, result)
	}
	return result
}

func shellConfigFiles(home string) []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{
			filepath.Join(home, ".zprofile"),
			filepath.Join(home, ".zshrc"),
			filepath.Join(home, ".profile"),
		}
	case "linux":
		return []string{
			filepath.Join(home, ".profile"),
			filepath.Join(home, ".bash_profile"),
			filepath.Join(home, ".bashrc"),
			filepath.Join(home, ".zshrc"),
		}
	case "windows":
		return []string{}
	default:
		return []string{}
	}
}

func parseEnvFromShellFile(filePath string, env map[string]string) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	exportRegex := regexp.MustCompile(`^\s*export\s+([A-Za-z_][A-Za-z0-9_]*)=(?:"([^"]*)"|'([^']*)'|([^\s#]*))`)
	assignRegex := regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)=(?:"([^"]*)"|'([^']*)'|([^\s#]*))`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = stripComment(line)
		var matches []string
		if matches = exportRegex.FindStringSubmatch(line); len(matches) > 0 {
			setParsedValue(matches, env)
			continue
		}
		if matches = assignRegex.FindStringSubmatch(line); len(matches) > 0 {
			setParsedValue(matches, env)
			continue
		}
	}
}

func setParsedValue(matches []string, env map[string]string) {
	key := matches[1]
	if !goEnvVars[key] {
		return
	}
	value := firstNonEmpty(matches[2], matches[3], matches[4])
	value = strings.TrimSpace(value)
	value = stripTrailingInlineComment(value)
	value = expandEnvVarsWithMap(value, env)
	env[key] = value
}

func envSliceToMap(env []string) map[string]string {
	m := make(map[string]string)
	for _, line := range env {
		pos := strings.Index(line, "=")
		if pos <= 0 {
			continue
		}
		m[line[:pos]] = line[pos+1:]
	}
	return m
}

func expandEnvVarsWithMap(value string, env map[string]string) string {
	return os.Expand(value, func(key string) string {
		if v, ok := env[key]; ok {
			return v
		}
		return os.Getenv(key)
	})
}

func mergePath(newPath, oldPath string) string {
	sep := string(os.PathListSeparator)
	items := []string{}
	seen := map[string]bool{}
	add := func(path string) {
		for _, p := range strings.Split(path, sep) {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			if !seen[p] {
				seen[p] = true
				items = append(items, p)
			}
		}
	}
	add(newPath)
	add(oldPath)
	return strings.Join(items, sep)
}

func firstNonEmpty(v ...string) string {
	for _, s := range v {
		if s != "" {
			return s
		}
	}
	return ""
}

func stripComment(s string) string {
	pos := strings.Index(s, "#")
	if pos >= 0 {
		return strings.TrimSpace(s[:pos])
	}
	return s
}

func stripTrailingInlineComment(s string) string {
	pos := strings.Index(s, " #")
	if pos >= 0 {
		return strings.TrimSpace(s[:pos])
	}
	return s
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
