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
	"strings"

	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 对象检查器 搜索过滤
//
// 属性/事件树是 LazVirtualStringTree, 不属于 ICustomTreeView,
// TreeFilterEdit 无法直连, 这里通过节点 IsVisible 手动过滤:
// 属性名命中或子孙命中则保留, 否则隐藏, 命中的 class 属性保留整个子树并自动展开。

var gInspectorTarget *TDesigningComponent // 检查器当前展示的组件

// inspectorSetTarget 更新检查器目标组件, 并按当前搜索词应用过滤
func inspectorSetTarget(comp *TDesigningComponent) {
	gInspectorTarget = comp
	if MainWindow.contentLayout != nil && MainWindow.contentLayout.layoutInspector != nil {
		inspector := MainWindow.contentLayout.layoutInspector
		inspector.applyFilter()
		// 有展示目标时启用底部属性文档面板
		inspector.setDocVisible(comp != nil)
	}
}

func (m *ContentLayoutInspector) applyFilter() {
	if m == nil || m.searchEdit == nil || gInspectorTarget == nil {
		return
	}
	filter := strings.ToLower(strings.TrimSpace(m.searchEdit.Filter()))
	applyInspectorTreeFilter(gInspectorTarget.propertyTree, filter)
	applyInspectorTreeFilter(gInspectorTarget.eventTree, filter)
}

func applyInspectorTreeFilter(tree lcl.ILazVirtualStringTree, filter string) {
	if tree == nil || !tree.IsValid() {
		return
	}
	tree.EndEditNode()
	// 前序收集全部节点, 保证倒序处理时子节点先于父节点
	var nodes []types.PVirtualNode
	for node := tree.GetFirst(false); node != 0; node = tree.GetNext(node, false) {
		nodes = append(nodes, node)
	}
	if filter == "" {
		for _, node := range nodes {
			tree.SetIsVisible(node, true)
		}
		return
	}
	matched := make(map[types.PVirtualNode]bool, len(nodes))
	for _, node := range nodes {
		if data := vtedit.GetPropertyNodeData(node); data != nil {
			if strings.Contains(strings.ToLower(data.Name()), filter) {
				matched[node] = true
			}
		}
	}
	// 自底向上传播: 命中节点的祖先链保留
	visible := make(map[types.PVirtualNode]bool, len(nodes))
	for i := len(nodes) - 1; i >= 0; i-- {
		node := nodes[i]
		vis := matched[node]
		if !vis {
			for child := tree.GetFirstChild(node); child != 0; child = tree.GetNextSibling(child) {
				if visible[child] {
					vis = true
					break
				}
			}
		}
		visible[node] = vis
		tree.SetIsVisible(node, vis)
	}
	// 命中的 class 属性(如 Font)保留整个子树, 便于过滤状态下编辑子属性
	for _, node := range nodes {
		if !matched[node] || !visible[node] {
			continue
		}
		for child := tree.GetFirstChild(node); child != 0; child = tree.GetNextSibling(child) {
			applyFilterShowSubTree(tree, child, visible)
		}
	}
	// 展开所有带可见子孙的节点, 保证命中属性直接可见
	for _, node := range nodes {
		if !visible[node] {
			continue
		}
		for child := tree.GetFirstChild(node); child != 0; child = tree.GetNextSibling(child) {
			if visible[child] {
				tree.SetExpanded(node, true)
				break
			}
		}
	}
}

func applyFilterShowSubTree(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, visible map[types.PVirtualNode]bool) {
	if !visible[node] {
		visible[node] = true
		tree.SetIsVisible(node, true)
	}
	for child := tree.GetFirstChild(node); child != 0; child = tree.GetNextSibling(child) {
		applyFilterShowSubTree(tree, child, visible)
	}
}
