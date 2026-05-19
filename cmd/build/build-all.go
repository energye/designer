package build

import (
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
)

// RunAll 执行全平台构建逻辑。
// 该函数会检查全局项目实例的有效性，并根据构建选项处理 CGO 环境变量的配置，
// 以支持在禁用 CGO 的情况下构建其他平台的目标。
func RunAll() {
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Build-All - project GProject is nil")
		return
	}
	/*
	 * 当项目配置为禁用 CGO 且启用了跨平台构建时，执行全平台编译流程。
	 * 首先备份当前的 CGO_ENABLED 环境变量并强制设为 "0"，同时设置全平台构建标记。
	 * 随后遍历预定义的目标平台列表，动态切换 GOOS 和 GOARCH 环境变量，
	 * 依次触发各平台的构建任务。函数退出前会自动恢复原始的环境变量状态。
	 */
	isBuildOtherPlatform := !proj.BuildOption.BuildCGOEnabled && proj.BuildOption.BuildOtherPlatform
	if isBuildOtherPlatform {
		for _, env := range buildPlatformENVs {
			//RunGoCleanCacheCMD()
			event.ConsoleWriteInfo("Build - GOOS:", env.OS, ", GOARCH:", env.ARCH)
			envs := Envs{
				"GOOS":           env.OS,
				"GOARCH":         env.ARCH,
				"CGO_ENABLED":    "0",
				buildAllPlatform: "true",
			}
			env.Run(envs)
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
	OS   string              // windows darwin linux
	ARCH string              // amd64 386 arm64 arm
	Run  func(env Envs) bool // build func
}
