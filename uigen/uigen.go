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
// See the License for the specific language governing permissions and
// limitations under the License.

package uigen

import (
	"encoding/json"
	"github.com/energye/designer/designer"
	"os"
	"path/filepath"
)

// UI 布局文件生成, JSON 格式, 自动时时生成
// 依赖 designer 生成 JSON UI文件(form[n].ui)
// 只存放被修改的组件属性
// xxx.ui 文件内容是 tree JSON 结构, 数据格式为组件[变更的属性列表]
// 生成触发条件: 即时触发 防抖

// 生成UI文件
func generateUIFile(formComponent *designer.TDesigningComponent, filePath string) error {
	// 构建UI树结构(组件子树序列化, 与剪贴板/撤销快照共用同一实现)
	uiTree := designer.BuildUIComponentTree(formComponent)

	// 序列化为JSON
	data, err := json.MarshalIndent(uiTree, "", "  ")
	//data, err := json.Marshal(uiTree)
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filePath, data, 0644)
}
