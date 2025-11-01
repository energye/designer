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

package main

import "github.com/energye/designer/pkg/scantypes/mappergen"

// 生成类型映射.go
// 用于在动态设置属性时

func main() {
	mappergen.LCLMapper()
	mappergen.CEFMapper()
	mappergen.WV2Mapper()
}
