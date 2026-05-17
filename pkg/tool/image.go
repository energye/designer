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
	"bytes"
	"errors"
	"github.com/energye/designer/pkg/draw"
	"github.com/energye/designer/pkg/logs"
	"image"
	"image/png"
)

// Scale 将PNG格式的图片数据缩放到指定尺寸
//
//	data: 原始PNG图片的字节数据
//	targetW: 目标宽度
//	targetH: 目标高度
//	[]byte: 缩放后的PNG图片字节数据，如果处理失败则返回nil
func Scale(data []byte, targetW, targetH int) []byte {
	pngBuf := &bytes.Buffer{}
	if _, err := pngBuf.Write(data); err != nil {
		logs.Error("Image Scale PNG Write Buffer:", err.Error())
		return nil
	}
	// 解码 png 到 image
	pngImg, err := png.Decode(pngBuf)
	if err != nil {
		logs.Error("Image Scale PNG Decode:", err.Error())
		return nil
	}
	pngBounds := pngImg.Bounds()
	// 存放缩放后的图像，使用传入的目标尺寸而非硬编码
	scaledImg := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	draw.CatmullRom.Scale(scaledImg, scaledImg.Bounds(), pngImg, pngBounds, draw.Over, nil)
	// 最后保存缩放 png
	scalePngBuf := &bytes.Buffer{}
	if err := png.Encode(scalePngBuf, scaledImg); err != nil {
		logs.Error("Image Scale PNG Encode Save Buffer:", err.Error())
		return nil
	}
	data = scalePngBuf.Bytes()
	return data
}

// 图片格式的魔数签名
var magicTable = map[string]string{
	"\xff\xd8\xff":      "jpeg",
	"\x89PNG\r\n\x1a\n": "png",
	"GIF87a":            "gif",
	"GIF89a":            "gif",
	"BM":                "bmp",
	"\x00\x00\x01\x00":  "ico",
}

// DetectImageFormatByte 检测图片真实格式
func DetectImageFormatByte(imageData []byte) (string, error) {
	if len(imageData) < 16 {
		return "", errors.New("图片太小")
	}
	buffer := imageData[:16]

	// 特殊处理ICO中嵌套PNG的情况
	if bytes.HasPrefix(buffer, []byte("\x00\x00\x01\x00")) {
		// 如果ICO文件内嵌PNG，从偏移量6开始检测
		if bytes.HasPrefix(buffer[6:], []byte("\x89PNG")) {
			return "png", nil
		}
		return "ico", nil
	}

	// 检查其他格式
	for magic, format := range magicTable {
		if bytes.HasPrefix(buffer, []byte(magic)) {
			return format, nil
		}
	}

	return "", errors.New("不支持的格式")
}
