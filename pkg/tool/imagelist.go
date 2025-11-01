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
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"strings"
)

// 图片资源

type ImageList struct {
	imageIndex   map[string]int32
	imageList100 lcl.IImageList
	imageList150 lcl.IImageList
	imageList200 lcl.IImageList
}

type ImageRect struct {
	Image100 types.TSize
	Image150 types.TSize
	Image200 types.TSize
}

func NewImageList(owner lcl.IComponent, dirName string, rect ImageRect) *ImageList {
	m := new(ImageList)
	m.imageIndex = make(map[string]int32)
	imageList := resources.GetImageFileList(dirName)
	var (
		images100 []string
		images150 []string
		images200 []string
	)
	for _, name := range imageList {
		is48 := strings.LastIndex(name, "_200.png") != -1
		is36 := strings.LastIndex(name, "_150.png") != -1
		if is48 {
			images200 = append(images200, name)
		} else if is36 {
			images150 = append(images150, name)
		} else {
			images100 = append(images100, name)
		}
	}
	loadImage := func(images []string, w, h int32) lcl.IImageList {
		resultImageList := LoadImageList(owner, images, w, h)
		for index, name := range images {
			name = strings.ToLower(strings.Replace(name, dirName+"/", "", 1))
			m.imageIndex[name] = int32(index)
		}
		return resultImageList
	}
	m.imageList100 = loadImage(images100, rect.Image100.Cx, rect.Image100.Cy)
	m.imageList150 = loadImage(images150, rect.Image150.Cx, rect.Image150.Cy)
	m.imageList200 = loadImage(images200, rect.Image200.Cx, rect.Image200.Cy)
	return m
}

func (m *ImageList) ImageIndex(name string) int32 {
	index, ok := m.imageIndex[strings.ToLower(name)]
	if ok {
		return index
	}
	return 0
}

func (m *ImageList) ImageList100() lcl.IImageList {
	return m.imageList100
}

func (m *ImageList) ImageList150() lcl.IImageList {
	return m.imageList150
}

func (m *ImageList) ImageList200() lcl.IImageList {
	return m.imageList200
}
