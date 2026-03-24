package run

import (
	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/tool/command"
	"path/filepath"
)

// Run 运行应用程序的入口函数
func Run(runCmd *command.CMD) bool {
	event.ConsoleWriteInfo("CMD-Run-run")
	proj := bean.GProject
	if proj == nil {
		event.ConsoleWriteError("Run - GProject is nil")
		return false
	}
	output := AppExecutable()
	if runCmd == nil {
		runCmd = command.NewCMD()
	}
	runCmd.IsPrint = false
	runCmd.HideWindow = true
	runCmd.Dir = bean.GPath
	runCmd.Console = func(data string, level command.Level) {
		event.ConsoleWriteInfo(data)
	}
	event.ConsoleWriteInfo("CMD-Run", output)
	runCmd.Command(output)
	return true
}

// AppExecutable 获取应用程序可执行文件的完整路径
//
//	根据项目构建选项和操作系统类型，计算并返回可执行文件的绝对路径
//	在 macOS 上会构建 .app 包的可执行文件路径，在 Windows/Linux 上直接返回构建文件路径
func AppExecutable() string {
	proj := bean.GProject
	option := proj.BuildOption
	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	if tool.IsWindows || tool.IsLinux {
		buildFileName := option.BuildFileName
		if tool.IsWindows && filepath.Ext(buildFileName) != ".exe" {
			buildFileName += ".exe"
		}
		buildFile := filepath.Join(output, buildFileName)
		return buildFile
	} else if tool.IsDarwin {
		packageName := option.PackageName + ".app"
		appRoot := filepath.Join(output, packageName)
		macOS := filepath.Join(appRoot, "Contents", "MacOS")
		buildFile := filepath.Join(macOS, proj.AppOption.MacOS.PList.CFBundleExecutable)
		return buildFile
	}
	return ""
}
