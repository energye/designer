package designer

import (
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
)

func LoadImageList(owner lcl.IComponent, imageList []string, width, height int32) lcl.IImageList {
	images := lcl.NewImageList(owner)
	images.SetWidth(width)
	images.SetHeight(height)
	for _, image := range imageList {
		tool.ImageListAddPng(images, image)
	}
	return images
}
