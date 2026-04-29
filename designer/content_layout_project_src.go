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
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/energye/designer/event"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
)

// 项目源码

var gProjectSrcTree = &TProjectSrcTree{}

type TProjectSrcTree struct {
	lock        sync.Mutex
	cancel      context.CancelFunc
	ticker      *time.Ticker
	projectPath string
	loadedDirs  map[string]bool
	nodePaths   map[uintptr]string
	pathNodes   map[string]lcl.ITreeNode
	pathIsDir   map[string]bool
	snapshots   map[string]string
}

type TProjectSrcEntry struct {
	Name        string
	Path        string
	IsDir       bool
	HasChildren bool
}

func (m *TProjectSrcTree) scanProjectSrc() {
	projectPath := projectSrcRootPath()
	if projectPath == "" {
		m.stop()
		m.clearOnMainThread()
		return
	}
	m.ensureProject(projectPath)

	rootHasChildren := m.hasVisibleChildren(projectPath)
	loadedDirs := m.snapshotLoadedDirs(projectPath)
	dirEntries := make(map[string][]TProjectSrcEntry)
	dirSnapshots := make(map[string]string)
	for _, dir := range loadedDirs {
		entries, snapshot := m.readDirEntries(dir)
		dirEntries[dir] = entries
		dirSnapshots[dir] = snapshot
	}

	lcl.RunOnMainThreadAsync(func(id uint32) {
		if projectSrcRootPath() != projectPath {
			return
		}
		m.applyScan(projectPath, rootHasChildren, dirEntries, dirSnapshots)
	})
}

func (m *TProjectSrcTree) start() {
	projectPath := projectSrcRootPath()
	if projectPath == "" {
		m.stop()
		m.clearOnMainThread()
		return
	}

	m.stop()
	m.ensureProject(projectPath)
	m.scanProjectSrc()

	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(time.Second)
	m.lock.Lock()
	m.cancel = cancel
	m.ticker = ticker
	m.lock.Unlock()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.scanProjectSrc()
			}
		}
	}()
}

func (m *TProjectSrcTree) stop() {
	m.lock.Lock()
	if m.ticker != nil {
		m.ticker.Stop()
		m.ticker = nil
	}
	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}
	m.lock.Unlock()
}

func (m *TProjectSrcTree) reset() {
	m.lock.Lock()
	m.projectPath = ""
	m.loadedDirs = nil
	m.nodePaths = nil
	m.pathNodes = nil
	m.pathIsDir = nil
	m.snapshots = nil
	m.lock.Unlock()
}

func (m *TProjectSrcTree) ensureProject(projectPath string) {
	projectPath = filepath.Clean(projectPath)
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.projectPath == projectPath && m.loadedDirs != nil {
		return
	}
	m.projectPath = projectPath
	m.loadedDirs = make(map[string]bool)
	m.nodePaths = make(map[uintptr]string)
	m.pathNodes = make(map[string]lcl.ITreeNode)
	m.pathIsDir = make(map[string]bool)
	m.snapshots = make(map[string]string)
}

func (m *TProjectSrcTree) snapshotLoadedDirs(projectPath string) []string {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.projectPath != projectPath || m.loadedDirs == nil {
		return []string{projectPath}
	}
	result := make([]string, 0, len(m.loadedDirs))
	for dir := range m.loadedDirs {
		result = append(result, dir)
	}
	sort.Strings(result)
	return result
}

func (m *TProjectSrcTree) readDirEntries(dir string) ([]TProjectSrcEntry, string) {
	dir = filepath.Clean(dir)
	dirItems, err := os.ReadDir(dir)
	if err != nil {
		return nil, ""
	}
	entries := make([]TProjectSrcEntry, 0, len(dirItems))
	for _, item := range dirItems {
		name := item.Name()
		if m.isHiddenSrcEntry(name) {
			continue
		}
		entry := TProjectSrcEntry{
			Name:  name,
			Path:  filepath.Join(dir, name),
			IsDir: item.IsDir(),
		}
		if entry.IsDir {
			entry.HasChildren = m.hasVisibleChildren(entry.Path)
		}
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	var snapshot strings.Builder
	for _, entry := range entries {
		if entry.IsDir {
			snapshot.WriteString("d:")
			if entry.HasChildren {
				snapshot.WriteString("1:")
			} else {
				snapshot.WriteString("0:")
			}
		} else {
			snapshot.WriteString("f:")
		}
		snapshot.WriteString(entry.Path)
		snapshot.WriteByte('\n')
	}
	return entries, snapshot.String()
}

func (m *TProjectSrcTree) isHiddenSrcEntry(name string) bool {
	return strings.HasPrefix(name, ".")
}

func (m *TProjectSrcTree) applyScan(projectPath string, rootHasChildren bool, dirEntries map[string][]TProjectSrcEntry, dirSnapshots map[string]string) {
	root := ProjectTreeSrcTreeNode()
	if root == nil || !root.IsValid() {
		return
	}
	m.ensureProject(projectPath)
	m.bindNode(projectPath, root, true)
	root.SetHasChildren(rootHasChildren || root.Count() > 0)
	ProjectTreeBeginUpdate()
	defer ProjectTreeEndUpdate()
	for dir, entries := range dirEntries {
		m.lock.Lock()
		oldSnapshot := m.snapshots[dir]
		newSnapshot := dirSnapshots[dir]
		node := m.pathNodes[dir]
		m.lock.Unlock()
		if node == nil || !node.IsValid() || oldSnapshot == newSnapshot {
			continue
		}
		m.reloadChildren(node, dir, entries)
		m.lock.Lock()
		m.snapshots[dir] = newSnapshot
		m.lock.Unlock()
	}
}

func (m *TProjectSrcTree) loadNode(node lcl.ITreeNode) {
	if node == nil || !node.IsValid() {
		return
	}
	m.lock.Lock()
	dir := m.nodePaths[node.Instance()]
	m.lock.Unlock()
	if dir == "" {
		root := ProjectTreeSrcTreeNode()
		if root == nil || !root.IsValid() || root.Instance() != node.Instance() {
			return
		}
		dir = projectSrcRootPath()
		m.bindNode(dir, node, true)
	}
	entries, snapshot := m.readDirEntries(dir)
	ProjectTreeBeginUpdate()
	m.reloadChildren(node, dir, entries)
	ProjectTreeEndUpdate()
	m.lock.Lock()
	m.loadedDirs[dir] = true
	m.snapshots[dir] = snapshot
	m.lock.Unlock()
}

func (m *TProjectSrcTree) reloadChildren(parent lcl.ITreeNode, parentPath string, entries []TProjectSrcEntry) {
	parentPath = filepath.Clean(parentPath)
	parentExpanded := parent.Expanded()
	existing := m.childNodesByPath(parent)
	nextPaths := make(map[string]TProjectSrcEntry, len(entries))
	items := ProjectTree().Items()

	for _, entry := range entries {
		entry.Path = filepath.Clean(entry.Path)
		nextPaths[entry.Path] = entry
	}

	for path, node := range existing {
		if _, ok := nextPaths[path]; !ok {
			m.removePathMappings(path, true)
			if node != nil && node.IsValid() {
				node.Delete()
			}
		}
	}

	for index, entry := range entries {
		entry.Path = filepath.Clean(entry.Path)
		child := existing[entry.Path]
		oldIsDir := m.isDirPath(entry.Path)
		if child == nil || !child.IsValid() {
			child = items.AddChild(parent, entry.Name)
		} else {
			if child.Text() != entry.Name {
				child.SetText(entry.Name)
			}
			if oldIsDir && !entry.IsDir {
				child.DeleteChildren()
				m.removePathMappings(entry.Path, false)
			}
		}
		m.setupNode(child, entry)
		m.bindNode(entry.Path, child, entry.IsDir)
		child.SetIndex(int32(index))
	}

	parent.SetHasChildren(len(entries) > 0)
	if parentExpanded && len(entries) > 0 {
		parent.SetExpanded(true)
	}
}

func (m *TProjectSrcTree) childNodesByPath(parent lcl.ITreeNode) map[string]lcl.ITreeNode {
	result := make(map[string]lcl.ITreeNode)
	if parent == nil || !parent.IsValid() {
		return result
	}
	for i := int32(0); i < parent.Count(); i++ {
		child := parent.Items(i)
		path := m.pathByNode(child)
		if path != "" {
			result[path] = child
		}
	}
	return result
}

func (m *TProjectSrcTree) setupNode(node lcl.ITreeNode, entry TProjectSrcEntry) {
	if entry.IsDir {
		node.SetImageIndex(projectSrcFolderIcon())
		node.SetSelectedIndex(projectSrcFolderIcon())
		node.SetHasChildren(entry.HasChildren)
	} else {
		node.SetImageIndex(projectSrcFileIcon(entry))
		node.SetSelectedIndex(projectSrcFileIcon(entry))
		node.SetHasChildren(false)
	}
}

func (m *TProjectSrcTree) removePathMappings(path string, includeSelf bool) {
	path = filepath.Clean(path)
	prefix := path + string(os.PathSeparator)
	m.lock.Lock()
	defer m.lock.Unlock()
	for currentPath, node := range m.pathNodes {
		if currentPath == path && !includeSelf {
			continue
		}
		if currentPath == path || strings.HasPrefix(currentPath, prefix) {
			delete(m.pathNodes, currentPath)
			delete(m.pathIsDir, currentPath)
			delete(m.snapshots, currentPath)
			delete(m.loadedDirs, currentPath)
			if node != nil {
				delete(m.nodePaths, node.Instance())
			}
		}
	}
}

func (m *TProjectSrcTree) bindNode(path string, node lcl.ITreeNode, isDir bool) {
	if node == nil || !node.IsValid() {
		return
	}
	path = filepath.Clean(path)
	m.lock.Lock()
	m.pathNodes[path] = node
	m.nodePaths[node.Instance()] = path
	m.pathIsDir[path] = isDir
	m.lock.Unlock()
}

func (m *TProjectSrcTree) isDirPath(path string) bool {
	path = filepath.Clean(path)
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.pathIsDir[path]
}

func (m *TProjectSrcTree) pathByNode(node lcl.ITreeNode) string {
	if node == nil || !node.IsValid() {
		return ""
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.nodePaths[node.Instance()]
}

func (m *TProjectSrcTree) nodeByPath(path string) lcl.ITreeNode {
	path = filepath.Clean(path)
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.pathNodes[path]
}

func (m *TProjectSrcTree) hasVisibleChildren(dir string) bool {
	dirItems, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, item := range dirItems {
		if !m.isHiddenSrcEntry(item.Name()) {
			return true
		}
	}
	return false
}

func (m *TProjectSrcTree) clearOnMainThread() {
	lcl.RunOnMainThreadAsync(func(id uint32) {
		ProjectTreeClearSrcTreeNode()
	})
}

func projectSrcFolderIcon() int32 {
	return imageComponents.ImageIndex("folder.png")
}

func projectSrcFileIcon(entry TProjectSrcEntry) int32 {
	fmt.Println("projectSrcFileIcon", entry)
	return imageComponents.ImageIndex("tfilenameedit.png")
}

func getFileIconIndex(fileName string) int32 {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".go":
		return imageComponents.ImageIndex("file_go.png")
	case ".mod":
		return imageComponents.ImageIndex("file_mod.png")
	case ".sum":
		return imageComponents.ImageIndex("file_sum.png")
	case ".json":
		return imageComponents.ImageIndex("file_json.png")
	case ".yaml", ".yml":
		return imageComponents.ImageIndex("file_yaml.png")
	case ".toml":
		return imageComponents.ImageIndex("file_toml.png")
	case ".md":
		return imageComponents.ImageIndex("file_md.png")
	case ".txt":
		return imageComponents.ImageIndex("file_txt.png")
	case ".xml":
		return imageComponents.ImageIndex("file_xml.png")
	case ".html", ".htm":
		return imageComponents.ImageIndex("file_html.png")
	case ".css":
		return imageComponents.ImageIndex("file_css.png")
	case ".js":
		return imageComponents.ImageIndex("file_js.png")
	case ".ts":
		return imageComponents.ImageIndex("file_ts.png")
	default:
		return imageComponents.ImageIndex("file.png")
	}
}

func projectSrcRootPath() string {
	if bean.GPath == "" {
		return ""
	}
	return filepath.Clean(bean.GPath)
}

func (m *ContentLayoutProject) TreeOnExpanding(sender lcl.IObject, node lcl.ITreeNode, allowExpansion *bool) {
	gProjectSrcTree.loadNode(node)
}

func ProjectSrcTreeNodePath(node lcl.ITreeNode) string {
	return gProjectSrcTree.pathByNode(node)
}

func ProjectSrcTreePathNode(path string) lcl.ITreeNode {
	return gProjectSrcTree.nodeByPath(path)
}

func initProjectSrcEvent() {
	logs.Println("启动项目 SRC Tree 监听")
	event.On(event.ListenProjectSrcFileChange, func(trigger event.TTrigger) {
		payload, ok := trigger.Payload.(event.TPayload)
		if ok {
			switch payload.Type {
			case event.ProjectSrcScan:
				gProjectSrcTree.start()
			}
		}
	}, func() {
		logs.Println("停止项目 SRC Tree 监听")
		gProjectSrcTree.stop()
	})
}
