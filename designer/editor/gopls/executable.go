package gopls

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// GoplsPathEnv 允许用户显式指定 gopls 可执行文件路径。
const GoplsPathEnv = "DESIGNER_GOPLS_PATH"

var (
	goEnvOnce  sync.Once
	goEnvCache map[string]string
)

// goEnvValue 先读进程环境变量，缺失时回退到 `go env -json`。
// designer 的 InitGoEnv 已跑过一次 go env，但结果只存在 config 包里；
// 这里自带一次带缓存的查询，避免 gopls 包反向依赖配置包。
func goEnvValue(name string) string {
	if v := strings.TrimSpace(os.Getenv(name)); v != "" {
		return v
	}
	goEnvOnce.Do(func() {
		goEnvCache = make(map[string]string)
		out, err := exec.Command("go", "env", "-json").Output()
		if err != nil {
			return
		}
		_ = json.Unmarshal(out, &goEnvCache)
	})
	return strings.TrimSpace(goEnvCache[name])
}

func goplsBinaryName() string {
	if runtime.GOOS == "windows" {
		return "gopls.exe"
	}
	return "gopls"
}

// resolveGoplsExecutable 按顺序定位 gopls：
// DESIGNER_GOPLS_PATH → go env GOBIN → GOPATH/bin → ~/go/bin → PATH。
// 返回的错误里带上所有尝试过的路径，便于用户排查。
func resolveGoplsExecutable() (string, error) {
	name := goplsBinaryName()
	var searched []string

	if p := strings.TrimSpace(os.Getenv(GoplsPathEnv)); p != "" {
		searched = append(searched, p)
		if isExecutableFile(p) {
			return p, nil
		}
	}

	if dir := goEnvValue("GOBIN"); dir != "" {
		p := filepath.Join(dir, name)
		searched = append(searched, p)
		if isExecutableFile(p) {
			return p, nil
		}
	}

	for _, dir := range filepath.SplitList(goEnvValue("GOPATH")) {
		if dir == "" {
			continue
		}
		p := filepath.Join(dir, "bin", name)
		searched = append(searched, p)
		if isExecutableFile(p) {
			return p, nil
		}
	}

	if home, err := os.UserHomeDir(); err == nil {
		p := filepath.Join(home, "go", "bin", name)
		searched = append(searched, p)
		if isExecutableFile(p) {
			return p, nil
		}
	}

	if p, err := exec.LookPath("gopls"); err == nil {
		return p, nil
	}
	searched = append(searched, "PATH/gopls")

	return "", fmt.Errorf(
		"未找到 gopls：已尝试 %s。请执行 `go install golang.org/x/tools/gopls@latest`，或用 %s 指定路径",
		strings.Join(searched, "、"), GoplsPathEnv)
}

func isExecutableFile(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && !fi.IsDir()
}
