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
	"path/filepath"
	"strings"
)

// 图片资源

type ImageList struct {
	imageIndex   map[string]int32
	imageList50  lcl.IImageList
	imageList100 lcl.IImageList
	imageList150 lcl.IImageList
	imageList200 lcl.IImageList
}

type ImageRect struct {
	Image50  types.TSize
	Image100 types.TSize
	Image150 types.TSize
	Image200 types.TSize
}

func AppendSVGToImageList(imageList *ImageList, dirName string, rect ImageRect) {
	svgImageList := resources.GetImageFileList(dirName)
	loadImage := func(image lcl.IImageList, size types.TSize, scaleName string) {
		for _, name := range svgImageList {
			svgData := resources.Images(name)
			if svgData == nil {
				continue
			}
			pngData, err := SVGToPNG(svgData, int(size.Cx), int(size.Cy))
			if err != nil {
				continue
			}
			ImageListAddPng(image, pngData)
			_, name = filepath.Split(name) // strings.ToLower(strings.Replace(name, dirName+"/", "", 1))
			name = strings.TrimSuffix(name, ".svg")
			if scaleName != "" {
				name = name + "_" + scaleName
			}
			name += ".png"
			count := image.Count()
			imageList.imageIndex[name] = count - 1
		}
	}
	if CheckSizeZero(rect.Image50) {
		loadImage(imageList.imageList50, rect.Image50, "50")
	}
	if CheckSizeZero(rect.Image100) {
		loadImage(imageList.imageList100, rect.Image100, "")
	}
	if CheckSizeZero(rect.Image150) {
		loadImage(imageList.imageList150, rect.Image150, "150")
	}
	if CheckSizeZero(rect.Image200) {
		loadImage(imageList.imageList200, rect.Image200, "200")
	}
}

func NewImageListSVGToPNG(owner lcl.IComponent, dirName string, rect ImageRect) *ImageList {
	m := new(ImageList)
	m.imageIndex = make(map[string]int32)
	svgImageList := resources.GetImageFileList(dirName)
	loadImage := func(size types.TSize, scaleName string) lcl.IImageList {
		imageList := lcl.NewImageList(owner)
		if size.Cx > 0 {
			imageList.SetWidth(size.Cx)
		}
		if size.Cy > 0 {
			imageList.SetHeight(size.Cy)
		}
		for index, name := range svgImageList {
			svgData := resources.Images(name)
			if svgData == nil {
				continue
			}
			pngData, err := SVGToPNG(svgData, int(size.Cx), int(size.Cy))
			if err != nil {
				continue
			}
			ImageListAddPng(imageList, pngData)
			name = strings.ToLower(strings.Replace(name, dirName+"/", "", 1))
			if scaleName != "" {
				name = name + "_" + scaleName
			}
			m.imageIndex[name] = int32(index)
		}
		return imageList
	}
	if CheckSizeZero(rect.Image50) {
		m.imageList50 = loadImage(rect.Image50, "50")
	}
	if CheckSizeZero(rect.Image100) {
		m.imageList100 = loadImage(rect.Image100, "")
	}
	if CheckSizeZero(rect.Image150) {
		m.imageList150 = loadImage(rect.Image150, "150")
	}
	if CheckSizeZero(rect.Image200) {
		m.imageList200 = loadImage(rect.Image200, "200")
	}
	return m
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
	if CheckSizeZero(rect.Image50) {
		m.imageList50 = loadImage(images100, rect.Image50.Cx, rect.Image50.Cy)
	}
	if CheckSizeZero(rect.Image100) {
		m.imageList100 = loadImage(images100, rect.Image100.Cx, rect.Image100.Cy)
	}
	if CheckSizeZero(rect.Image150) {
		m.imageList150 = loadImage(images150, rect.Image150.Cx, rect.Image150.Cy)
	}
	if CheckSizeZero(rect.Image200) {
		m.imageList200 = loadImage(images200, rect.Image200.Cx, rect.Image200.Cy)
	}
	return m
}

func CheckSizeZero(size types.TSize) bool {
	if size.Cx == 0 || size.Cy == 0 {
		return false
	}
	return true
}

func (m *ImageList) ImageIndex(name string) int32 {
	index, ok := m.imageIndex[strings.ToLower(name)]
	if ok {
		return index
	}
	return 0
}

func (m *ImageList) ImageList50() lcl.IImageList {
	return m.imageList50
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
