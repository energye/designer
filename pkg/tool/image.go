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
		logs.Error("图标加载 PNG Write Buffer:", err.Error())
		return nil
	}
	// 解码 png 到 image
	pngImg, err := png.Decode(pngBuf)
	if err != nil {
		logs.Error("图标加载 PNG Decode:", err.Error())
		return nil
	}
	pngBounds := pngImg.Bounds()
	// 存放缩放后的图像，使用传入的目标尺寸而非硬编码
	scaledImg := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	draw.CatmullRom.Scale(scaledImg, scaledImg.Bounds(), pngImg, pngBounds, draw.Over, nil)
	// 最后保存缩放 png
	scalePngBuf := &bytes.Buffer{}
	if err := png.Encode(scalePngBuf, scaledImg); err != nil {
		logs.Error("图标加载 PNG Encode Save Buffer:", err.Error())
		return nil
	}
	data = scalePngBuf.Bytes()
	return data
}
