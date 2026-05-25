package build

import (
	"context"
	"github.com/energye/designer/cmd/env"
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
)

// RunAll 执行全平台构建逻辑。
// 该函数会检查全局项目实例的有效性，并根据构建选项处理 CGO 环境变量的配置，
// 以支持在禁用 CGO 的情况下构建其他平台的目标。
func RunAll(ctx context.Context) {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Build-All - project GProject is nil")
		return
	}
	// 当项目配置为禁用 CGO 且启用了跨平台构建时，执行全平台编译流程。
	// 依次触发各平台的构建任务。
	isBuildOtherPlatform := !proj.BuildOption.BuildCGOEnabled && proj.BuildOption.BuildOtherPlatform
	if isBuildOtherPlatform {
		defer env.Clear()
		for _, buildEnv := range buildPlatformENVs {
			//RunGoCleanCacheCMD()
			event.ConsoleWriteInfo("Build - GOOS:", buildEnv.OS, ", GOARCH:", buildEnv.ARCH)
			env.Put("GOOS", buildEnv.OS)
			env.Put("GOARCH", buildEnv.ARCH)
			env.Put("CGO_ENABLED", "0")
			env.Put(buildAllPlatform, "true")
			buildEnv.Run(ctx)
		}
	}
}

var buildPlatformENVs = []buildPlatformENV{
	{OS: "windows", ARCH: "amd64", Run: buildWindows},
	{OS: "windows", ARCH: "386", Run: buildWindows},
	{OS: "darwin", ARCH: "amd64", Run: buildDarwin},
	{OS: "darwin", ARCH: "arm64", Run: buildDarwin},
	{OS: "linux", ARCH: "amd64", Run: buildLinux},
	{OS: "linux", ARCH: "386", Run: buildLinux},
	{OS: "linux", ARCH: "arm", Run: buildLinux},
	{OS: "linux", ARCH: "arm64", Run: buildLinux},
}

type buildPlatformENV struct {
	OS   string                         // windows darwin linux
	ARCH string                         // amd64 386 arm64 arm
	Run  func(ctx context.Context) bool // build func
}
