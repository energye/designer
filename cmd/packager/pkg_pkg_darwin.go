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

//go:build darwin

package packager

import (
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"path/filepath"
)

func packager() {
	proj := bean.GProject
	if proj == nil {
		logs.Error("项目对象 GProject 为 nil")
		return
	}
	logs.Info("打包项目, 检查配置选项")
	option := proj.BuildOption

	output := option.Output
	if !filepath.IsAbs(option.Output) {
		output = filepath.Join(bean.GPath, output)
	}
	outputFilename := filepath.Join(output, option.BuildFileName)
	logs.Info("Packaging", outputFilename)
}
