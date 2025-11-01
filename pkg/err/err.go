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

package err

// 简单的错误管理

// 返回状态类型
type ResultStatus int32

const (
	RsSuccess       ResultStatus = iota // 成功 通用
	RsFail                              // 失败 通用
	RsIgnoreProp                        // 忽略的属性
	RsNotValid                          // 对象无效
	RsDuplicateName                     // 组件名重复
	RsNoModify                          // 不做变更
)

// 检测 err
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
