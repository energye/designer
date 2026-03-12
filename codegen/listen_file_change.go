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

package codegen

import (
	"context"
	"github.com/energye/designer/designer"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"go/ast"
	"os"
	"path/filepath"
	"time"
)

var (
	gCtx      context.Context
	gCancel   context.CancelFunc
	gTicker   *time.Ticker
	gFileInfo *tool.HashMap[string, *TListenFileInfo]
)

type TListenFileInfo struct {
	Id      int       // 窗体 id
	Path    string    // 文件路径
	ModTime time.Time // 修改时间
	Exist   bool      // 是否存在
}

func stopListenFileChange() {
	if gCtx != nil {
		logs.Println("ListenFileChange 停止运行监听文件改变任务")
		gCancel()
		gCtx = nil
		gCancel = nil
	}
	if gTicker != nil {
		logs.Println("ListenFileChange stop")
		gTicker.Stop()
		gTicker = nil
	}
}

// 监听文件改变任务
func runListenFileChange() {
	stopListenFileChange()
	logs.Println("ListenFileChange 启动监听文件改变任务")
	gCtx, gCancel = context.WithCancel(context.Background())
	gTicker = time.NewTicker(time.Second) // 1 秒？
	// 创建(重置)
	gFileInfo = tool.NewHashMap[string, *TListenFileInfo]()
	go func() {
		defer func() {
			logs.Println("ListenFileChange done.")
		}()
		for {
			select {
			case <-gCtx.Done():
				return
			case <-gTicker.C:
				detectFileChange()
			}
		}
	}()
}

// DetectFileChange 检测文件
func detectFileChange() {
	projPath := bean.GPath
	proj := bean.GProject
	if proj == nil || projPath == "" {
		return
	}
	// 重置文件存在状态, 默认 false
	gFileInfo.Iterate(func(key string, value *TListenFileInfo) bool {
		value.Exist = false
		return false
	})

	// 重新定位文件信息
	appCodePath := bean.CodePath()
	for _, form := range bean.GProject.UIForms {
		userFile := filepath.Join(appCodePath, form.GOUserFile)
		if fi := gFileInfo.Get(form.Name); fi != nil {
			fi.Exist = true
		} else {
			gFileInfo.Add(form.Name, &TListenFileInfo{Id: form.Id, Path: userFile, ModTime: time.Now(), Exist: true})
		}
	}

	// 待删除不存在的文件列表
	var removeNotExistList []string
	// 验证文件状态
	gFileInfo.Iterate(func(formName string, file *TListenFileInfo) bool {
		if !file.Exist {
			// 不存在的添加到待删除列表
			removeNotExistList = append(removeNotExistList, formName)
			return false
		}
		fi, err := os.Stat(file.Path)
		if err != nil {
			return false
		}
		// 使用修改时间判断
		if !fi.ModTime().Equal(file.ModTime) {
			file.ModTime = fi.ModTime()
			// 当文件改变后处理对应逻辑
			// 功能:
			//   ast 获得窗体引用接收者的所有方法
			//   其它待添加... 如: 代码编辑检测
			logs.Debug("ListenFileChange 文件被修改 窗体名:", formName, "窗体ID", file.Id, "文件:", file.Path)

			// ast 加载窗体引用接收者的所有方法
			loadFormRecvMethods(file.Path, formName, file.Id)
		}
		return false
	})

	for _, name := range removeNotExistList {
		logs.Debug("ListenFileChange 在检测列表删除不存在的:", name)
		gFileInfo.Remove(name)
	}
}

// ast 加载窗体引用接收者的所有方法
func loadFormRecvMethods(filePath, formName string, formId int) {
	formName = "T" + formName
	var funcs []*dast.TFuncInfo
	dast.FindRecvMethod(filePath, formName, func(funcDecl *ast.FuncDecl) {
		name := funcDecl.Name.Name
		params := funcDecl.Type.Params
		results := funcDecl.Type.Results
		funcInfo := &dast.TFuncInfo{
			Name:    name,
			Params:  dast.ParseFields(params),
			Results: dast.ParseFields(results),
		}
		funcs = append(funcs, funcInfo)
	})
	designer.SetRecvMethods(formId, funcs)
}
