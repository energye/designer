package tool

import (
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/lcl"
	"strings"
)

// 图片资源

type ImageList struct {
	imageIndex  map[string]int32
	imageList24 lcl.IImageList
	imageList36 lcl.IImageList
	imageList48 lcl.IImageList
}

func NewImageList(owner lcl.IComponent, dirName string) *ImageList {
	m := new(ImageList)
	imageList := resources.GetImageFileList(dirName)
	var (
		images24 []string
		images36 []string
		images48 []string
	)
	for _, name := range imageList {
		is48 := strings.LastIndex(name, "_200.png") != -1
		is36 := strings.LastIndex(name, "_150.png") != -1
		if is48 {
			images48 = append(images48, name)
		} else if is36 {
			images36 = append(images36, name)
		} else {
			images24 = append(images24, name)
		}
	}
	loadImage := func(images []string, w, h int32) lcl.IImageList {
		resultImageList := LoadImageList(owner, imageList, w, h)
		m.imageIndex = map[string]int32{}
		for index, name := range images {
			name = strings.Replace(name, dirName+"/", "", 1)
			m.imageIndex[name] = int32(index)
		}
		return resultImageList
	}
	m.imageList24 = loadImage(images24, 24, 24)
	m.imageList36 = loadImage(images36, 36, 36)
	m.imageList48 = loadImage(images48, 48, 48)
	return m
}

func (m *ImageList) ImageIndex(name string) int32 {
	return m.imageIndex[name]
}

func (m *ImageList) ImageList24() lcl.IImageList {
	return m.imageList24
}

func (m *ImageList) ImageList36() lcl.IImageList {
	return m.imageList36
}

func (m *ImageList) ImageList48() lcl.IImageList {
	return m.imageList48
}
