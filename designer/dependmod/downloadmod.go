// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package dependmod

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/tool/command"
	"path/filepath"
	"strings"
)

// 根据 designer/resources/config.json 配置依赖模块下载模块
func downloadMod(dir *modCacheDir) bool {
	logs.Println("Download dependencies from designer/resources/config.json")
	designerCfg := config.DesignerConfig
	dependencies := designerCfg.Dependencies
	lclVer := dependencies.Get(config.LCLModPath)
	cefVer := dependencies.Get(config.CEFModPath)
	wvVer := dependencies.Get(config.WVModPath)
	energyVer := dependencies.Get(config.ENERGYModPath)
	dependenciesInfo := fmt.Sprintf("Dependencies LCL: %s, CEF: %s, WV: %s, ENERGY: %s", lclVer, cefVer, wvVer, energyVer)
	event.ConsoleWriteInfo(dependenciesInfo)

	// 处理依赖模块路径
	modCachePath := config.GGoEnv.Get("GOMODCACHE")
	event.ConsoleWriteInfo("Download modCachePath:", modCachePath)
	if modCachePath != "" && tool.IsExist(modCachePath) {
		dir.lclDir = filepath.Join(modCachePath, fmt.Sprintf("%s@%s", config.LCLModPath, lclVer))
		dir.cefDir = filepath.Join(modCachePath, fmt.Sprintf("%s@%s", config.CEFModPath, cefVer))
		dir.wvDir = filepath.Join(modCachePath, fmt.Sprintf("%s@%s", config.WVModPath, wvVer))
		dir.engDir = filepath.Join(modCachePath, fmt.Sprintf("%s@%s", config.ENERGYModPath, energyVer))
	}
	ok := true

	// os.Setenv("GOPROXY", "https://goproxy.io,direct")
	// os.Setenv("GO111MODULE", "on")

	runDownloadModCache := func(modPath string) (result []string) {
		version := dependencies[modPath]
		// go mod download -json github.com/energye/energy/v3@v3.0.0
		modPath = fmt.Sprintf("%s@%s", modPath, version)
		event.ConsoleWriteInfo("Download module cache:", modPath)
		cmdArgs := []string{"mod", "download", "-json", modPath}
		cmd := command.NewCMD()
		cmd.IsPrint = false
		cmd.HideWindow = true
		cmd.Console = func(data string, level command.Level) {
			logs.Println(data)
			if level == command.LError {
				ok = false
				event.ConsoleWriteError(data)
			} else {
				result = append(result, data)
			}
		}
		cmd.Command("go", cmdArgs...)
		event.ConsoleWriteInfo("Download module cache:", modPath, "END")
		return
	}
	event.ConsoleWriteInfo("Download module cache start")
	parserModCacheDir := func(modPath string) (string, error) {
		modJSON := runDownloadModCache(modPath)
		if !ok {
			return "", errors.New("-")
		}
		tmpDir := struct {
			Dir string `json:"Dir"`
		}{}
		jsonStr := strings.Join(modJSON, "")
		err := json.Unmarshal([]byte(jsonStr), &tmpDir)
		if err != nil {
			event.ConsoleWriteError("Parser Mod Cache:", jsonStr, err.Error())
		}
		return tmpDir.Dir, err
	}
	if !tool.IsExist(dir.lclDir) {
		if tmpModDir, err := parserModCacheDir(config.LCLModPath); err == nil {
			dir.lclDir = tmpModDir
		} else {
			return false
		}
	}
	if !tool.IsExist(dir.cefDir) {
		if tmpModDir, err := parserModCacheDir(config.CEFModPath); err == nil {
			dir.cefDir = tmpModDir
		} else {
			return false
		}
	}
	if !tool.IsExist(dir.wvDir) {
		if tmpModDir, err := parserModCacheDir(config.WVModPath); err == nil {
			dir.wvDir = tmpModDir
		} else {
			return false
		}
	}
	if !tool.IsExist(dir.engDir) {
		if tmpModDir, err := parserModCacheDir(config.ENERGYModPath); err == nil {
			dir.engDir = tmpModDir
		} else {
			return false
		}
	}
	event.ConsoleWriteInfo("Download module cache end")
	return ok
}
