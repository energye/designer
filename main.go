package main

import (
	"github.com/energye/designer/designer"
	"github.com/energye/designer/pkg/logs"
	_ "github.com/energye/designer/pkg/syso"
	"github.com/energye/designer/resources"
	"github.com/energye/designer/uigen"
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/tool"
	"os"
	"path/filepath"
)

func main() {
	logs.Level = logs.LevelDebug
	libname.LibName = func() string {
		wd, _ := os.Getwd()
		return filepath.Join(wd, "../", "gen", "gout", "liblcl.dll")
	}()
	if !tool.IsExist(libname.LibName) {
		libname.LibName = resources.LibPath
	}
	lcl.Init(nil, nil)
	// 初始化UI生成
	uigen.InitUIGeneration()
	// 运行设计器
	designer.Run()
}
