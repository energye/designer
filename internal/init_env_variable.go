//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

//go:build darwin || linux

package internal

import (
	"bufio"
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

func init() {
	env := BuildRuntimeEnv()
	for name, value := range env {
		if goEnvVars[name] {
			println(name, "=", value)
			_ = os.Setenv(name, value)
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

		// 去掉行尾注释（简单处理）
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
