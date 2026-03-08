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

package build

import (
	"gen/tool"
	"github.com/energye/designer/pkg/logs"
	"strings"
)

func Run() {
	logs.Info("CMD-build-run")
	build()
}

func mergeTags(defaultTags, customTags string) (mergedTags []string) {
	defaultTagsArr := strings.FieldsFunc(defaultTags, func(r rune) bool {
		return r == ',' || r == ' '
	})
	customTagsArr := strings.FieldsFunc(customTags, func(r rune) bool {
		return r == ',' || r == ' '
	})
	tagMap := tool.NewArrayMap[string]()
	for _, tag := range defaultTagsArr {
		tag = strings.TrimSpace(tag)
		if tag == "" || tagMap.ContainsKey(tag) {
			continue
		}
		tagMap.Add(tag, tag)
	}
	for _, tag := range customTagsArr {
		tag = strings.TrimSpace(tag)
		if tag == "" || tagMap.ContainsKey(tag) {
			continue
		}
		tagMap.Add(tag, tag)
	}
	tagMap.Iterate(func(_ string, value string) {
		mergedTags = append(mergedTags, value)
	})
	return
}

func mergeLdflags(defaultLdflags, customLdflags string) (mergedLdFlags []string) {
	flagMap := tool.NewArrayMap[string]()
	defaultLdflagsArr := strings.Fields(defaultLdflags)
	for _, flag := range defaultLdflagsArr {
		flag = strings.TrimSpace(flag)
		if flag == "" || flagMap.ContainsKey(flag) {
			continue
		}
		flagMap.Add(flag, flag)
	}
	customLdflagsArr := strings.Fields(customLdflags)
	for _, flag := range customLdflagsArr {
		flag = strings.TrimSpace(flag)
		if flag == "" || flagMap.ContainsKey(flag) {
			continue
		}
		flagMap.Add(flag, flag)
	}
	flagMap.Iterate(func(_ string, value string) {
		mergedLdFlags = append(mergedLdFlags, value)
	})
	return
}
