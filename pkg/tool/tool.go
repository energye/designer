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
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"os"
	"strings"
)

// 加载图像到列表
func ImageListAddPng(imageList lcl.IImageList, filePath string) {
	data := resources.Images(filePath)
	if data != nil {
		pic := lcl.NewPicture()
		defer pic.Free()

		mem := lcl.NewMemoryStream()
		defer mem.Free()
		lcl.StreamHelper.WriteBuffer(mem, data)
		mem.SetPosition(0)
		pic.LoadFromStream(mem)
		imageList.Add(pic.Bitmap(), nil)
	} else {
		data = resources.Images("components/default_150.png")
		pic := lcl.NewPicture()
		defer pic.Free()

		mem := lcl.NewMemoryStream()
		defer mem.Free()
		lcl.StreamHelper.WriteBuffer(mem, data)
		mem.SetPosition(0)
		pic.LoadFromStream(mem)
		imageList.Add(pic.Bitmap(), nil)
	}
}

// 加载图片列表
func LoadImageList(owner lcl.IComponent, imageList []string, width, height int32) lcl.IImageList {
	images := lcl.NewImageList(owner)
	if width > 0 {
		images.SetWidth(width)
	}
	if height > 0 {
		images.SetHeight(height)
	}
	for _, image := range imageList {
		ImageListAddPng(images, image)
	}
	return images
}

// 判断字符串相等, 忽略大小写
func Equal(s1 string, s2 ...string) bool {
	s1 = strings.ToLower(s1)
	for _, s := range s2 {
		s = strings.ToLower(s)
		if s1 == s {
			return true
		}
	}
	return false
}

// 删除道字母 T
func RemoveT(name string) string {
	if name == "" {
		return ""
	}
	if name[0] == 'T' {
		return name[1:]
	}
	return name
}

// 第一个字母转为大写
func FirstToUpper(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}

// 判断文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		} else if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

// 字符串数组元素反转
func StringArrayReverse(array []string) {
	n := len(array)
	for i := 0; i < n/2; i++ {
		j := n - 1 - i
		array[i], array[j] = array[j], array[i]
	}
}

// 分割字符串并去除空字符串 aa,bb,,cc > [aa,bb,cc]
func Split(s, sep string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var result []string
	for _, v := range strings.Split(s, sep) {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		result = append(result, v)
	}
	return result[:]
}

// 字符串替换
func Replace(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

// 判断当前是否为主线程
func IsMainThread() bool {
	return api.MainThreadId() == api.CurrentThreadId()
}
