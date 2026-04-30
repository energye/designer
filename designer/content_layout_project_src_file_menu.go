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

package designer

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
)

// 源码文件右键菜单

type TSrcFileMenu struct {
	treePopupMenu lcl.IPopupMenu
	// 文件菜单项
	openFile     lcl.IMenuItem // 打开文件（在编辑器中）
	openExternal lcl.IMenuItem // 用外部编辑器打开
	copyPath     lcl.IMenuItem // 复制绝对路径
	copyRelPath  lcl.IMenuItem // 复制相对路径
	separator1   lcl.IMenuItem // 分隔线
	separator2   lcl.IMenuItem // 分隔线
	refresh      lcl.IMenuItem // 刷新
	// 文件夹菜单项
	newFile   lcl.IMenuItem // 新建文件
	newFolder lcl.IMenuItem // 新建文件夹
}

// 初始化源码文件右键菜单
func initSrcFileMenu(owner lcl.IComponent) *TSrcFileMenu {
	menu := new(TSrcFileMenu)
	menu.treePopupMenu = lcl.NewPopupMenu(owner)
	menu.treePopupMenu.SetImages(imageActions.ImageList100())
	menu.treePopupMenu.SetParent(owner)
	items := menu.treePopupMenu.Items()

	// 打开文件
	openFile := lcl.NewMenuItem(owner)
	openFile.SetCaption("打开文件")
	openFile.SetImageIndex(imageActions.ImageIndex("laz_open.png"))
	openFile.SetOnClick(menu.OnOpenFile)
	menu.openFile = openFile
	items.Add(openFile)

	// 用外部编辑器打开
	openExternal := lcl.NewMenuItem(owner)
	openExternal.SetCaption("用外部编辑器打开")
	openExternal.SetImageIndex(imageActions.ImageIndex("laz_edit.png"))
	openExternal.SetOnClick(menu.OnOpenExternal)
	menu.openExternal = openExternal
	items.Add(openExternal)

	// 分隔线1
	separator1 := lcl.NewMenuItem(owner)
	separator1.SetCaption("-")
	menu.separator1 = separator1
	items.Add(separator1)

	// 新建文件
	newFile := lcl.NewMenuItem(owner)
	newFile.SetCaption("新建文件")
	newFile.SetOnClick(menu.OnNewFile)
	menu.newFile = newFile
	items.Add(newFile)

	// 新建文件夹
	newFolder := lcl.NewMenuItem(owner)
	newFolder.SetCaption("新建文件夹")
	newFolder.SetOnClick(menu.OnNewFolder)
	menu.newFolder = newFolder
	items.Add(newFolder)

	// 分隔线2
	separator2 := lcl.NewMenuItem(owner)
	separator2.SetCaption("-")
	menu.separator2 = separator2
	items.Add(separator2)

	// 复制路径
	copyPath := lcl.NewMenuItem(owner)
	copyPath.SetCaption("复制路径")
	copyPath.SetOnClick(menu.OnCopyPath)
	menu.copyPath = copyPath
	items.Add(copyPath)

	// 复制相对路径
	copyRelPath := lcl.NewMenuItem(owner)
	copyRelPath.SetCaption("复制相对路径")
	copyRelPath.SetOnClick(menu.OnCopyRelPath)
	menu.copyRelPath = copyRelPath
	items.Add(copyRelPath)

	// 刷新
	refresh := lcl.NewMenuItem(owner)
	refresh.SetCaption("刷新")
	refresh.SetOnClick(menu.OnRefresh)
	menu.refresh = refresh
	items.Add(refresh)

	return menu
}

// SetupForFile 为文件节点设置菜单项可见性
func (m *TSrcFileMenu) SetupForFile() {
	m.openFile.SetVisible(true)
	m.openExternal.SetVisible(true)
	m.newFile.SetVisible(false)
	m.newFolder.SetVisible(false)
	m.separator1.SetVisible(true)
	m.separator2.SetVisible(true)
}

// SetupForFolder 为文件夹节点设置菜单项可见性
func (m *TSrcFileMenu) SetupForFolder() {
	m.openFile.SetVisible(false)
	m.openExternal.SetVisible(false)
	m.newFile.SetVisible(true)
	m.newFolder.SetVisible(true)
	m.separator1.SetVisible(true)
	m.separator2.SetVisible(true)
}

// selectedSrcPath 获取当前选中源码节点的文件路径
func (m *TSrcFileMenu) selectedSrcPath() string {
	node := ProjectTreeSelectNode()
	if node == nil || !node.IsValid() {
		return ""
	}
	return ProjectSrcTreeNodePath(node)
}

// OnOpenFile 在编辑器中打开文件
func (m *TSrcFileMenu) OnOpenFile(sender lcl.IObject) {
	filePath := m.selectedSrcPath()
	if filePath == "" {
		return
	}
	// 检查是否为目录
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return
	}
	logs.Debug("打开文件: ", filePath)
}

// OnOpenExternal 用外部编辑器打开文件
func (m *TSrcFileMenu) OnOpenExternal(sender lcl.IObject) {
	filePath := m.selectedSrcPath()
	if filePath == "" {
		return
	}
	logs.Debug("用外部编辑器打开: ", filePath)
	openFileWithDefaultApp(filePath)
}

// OnCopyPath 复制绝对路径到剪贴板
func (m *TSrcFileMenu) OnCopyPath(sender lcl.IObject) {
	filePath := m.selectedSrcPath()
	if filePath == "" {
		return
	}
	lcl.Clipboard.SetAsText(filePath)
	logs.Debug("复制路径: ", filePath)
}

// OnCopyRelPath 复制相对路径到剪贴板
func (m *TSrcFileMenu) OnCopyRelPath(sender lcl.IObject) {
	filePath := m.selectedSrcPath()
	if filePath == "" {
		return
	}
	relPath, err := filepath.Rel(bean.GPath, filePath)
	if err != nil {
		relPath = filePath
	}
	lcl.Clipboard.SetAsText(relPath)
	logs.Debug("复制相对路径: ", relPath)
}

// OnNewFile 新建文件
func (m *TSrcFileMenu) OnNewFile(sender lcl.IObject) {
	dirPath := m.selectedSrcPath()
	if dirPath == "" {
		return
	}
	// 如果选中的是文件，取其父目录
	info, err := os.Stat(dirPath)
	if err != nil {
		return
	}
	if !info.IsDir() {
		dirPath = filepath.Dir(dirPath)
	}
	logs.Debug("新建文件目录: ", dirPath)
	// TODO: 弹出输入框让用户输入文件名，然后创建文件
}

// OnNewFolder 新建文件夹
func (m *TSrcFileMenu) OnNewFolder(sender lcl.IObject) {
	dirPath := m.selectedSrcPath()
	if dirPath == "" {
		return
	}
	// 如果选中的是文件，取其父目录
	info, err := os.Stat(dirPath)
	if err != nil {
		return
	}
	if !info.IsDir() {
		dirPath = filepath.Dir(dirPath)
	}
	logs.Debug("新建文件夹目录: ", dirPath)
	// TODO: 弹出输入框让用户输入文件夹名，然后创建文件夹
}

// OnRefresh 刷新目录树
func (m *TSrcFileMenu) OnRefresh(sender lcl.IObject) {
	logs.Debug("刷新项目源码树")
	gProjectSrcTree.scanProjectSrc()
}

// openFileWithDefaultApp 使用系统默认程序打开文件
func openFileWithDefaultApp(filePath string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	default:
		cmd = exec.Command("xdg-open", filePath)
	}
	if cmd != nil {
		cmd.Start()
	}
}
