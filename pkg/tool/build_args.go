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

package tool

import (
	"regexp"
	"strings"
)

func MergeTags(defaultTags, customTags string) (mergedTags []string) {
	defaultTagsArr := strings.FieldsFunc(defaultTags, func(r rune) bool {
		return r == ',' || r == ' '
	})
	customTagsArr := strings.FieldsFunc(customTags, func(r rune) bool {
		return r == ',' || r == ' '
	})
	tagMap := NewArrayMap[string, string]()
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
	tagMap.Iterate(func(_ string, value string) bool {
		mergedTags = append(mergedTags, value)
		return false
	})
	return
}

func MergeLdflags(defaultLdflags, customLdflags string) (mergedLdFlags []string) {
	flagMap := NewArrayMap[string, string]()
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
	flagMap.Iterate(func(_ string, value string) bool {
		mergedLdFlags = append(mergedLdFlags, value)
		return false
	})
	return
}

// ExtractOtherBuildArgs 从构建参数字符串中提取其他构建参数
// 该函数会移除 -tags 和 -ldflags 相关的参数，并清理格式
//
//	buildArgs: 原始的构建参数字符串，可能包含各种 go build 选项
//	string: 清理后的构建参数字符串，移除了 -tags 和 -ldflags 参数，
//	        所有单引号已转换为双引号，多余空格已被清理
func ExtractOtherBuildArgs(buildArgs string) []string {
	buildArgs = strings.ReplaceAll(buildArgs, "'", "\"")
	reTags := regexp.MustCompile(`-tags\s+("?[^"\s-]+"?|\s*"[^"]+")`)
	reLdflags := regexp.MustCompile(`-ldflags\s+"[^"]+"`)
	cleanArgs := reTags.ReplaceAllString(buildArgs, "")
	cleanArgs = reLdflags.ReplaceAllString(cleanArgs, "")
	cleanArgs = regexp.MustCompile(`\s+`).ReplaceAllString(cleanArgs, " ")
	cleanArgs = strings.TrimSpace(cleanArgs)
	return strings.Fields(cleanArgs)
}
