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

package mapper

import "github.com/energye/designer/pkg/dast"

// 获取映射的类型值
func GetLCL(name string) any {
	// TODO 这个文件需要通过配置动态获取
	fileName := "C:\\app\\workspace\\gen\\gout\\lcl\\go\\types\\lcl.go"
	val := dast.GetConstValue(fileName, name)
	return val
}
