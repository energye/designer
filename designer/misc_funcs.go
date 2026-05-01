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

package designer

import (
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

func setFormDesignPageStyle(page *wg.TPage, icon []byte) {
	tabColor := colors.ClWhite
	btnColor := colors.RGBToColor(234, 239, 249)
	if icon != nil && len(icon) > 0 {
		page.Button().SetIconFavoriteFormBytes(icon)
	}
	page.Button().SetWidth(65)
	page.Button().Font().SetColor(colors.ClBlack)
	page.Button().RoundedCorner = types.NewSet(wg.RcLeftTop, wg.RcRightTop)
	page.Button().TextOffSetX = 10
	page.Button().SetBorderColor(wg.BbdNone, wg.LightenColor(btnColor, 0.8))
	page.Button().SetRadius(5)
	page.Button().SetColor(tabColor)
	page.Button().SetDownColor(wg.LightenColor(btnColor, 0.3), wg.LightenColor(btnColor, 0.5))
	page.Button().SetEnterColor(wg.LightenColor(btnColor, 0.1), wg.LightenColor(btnColor, 0.3))
	page.SetDefaultColor(tabColor)
	page.SetActiveColor(btnColor)
	page.Button().SetCursor(types.CrHandPoint)
	page.SetColor(tabColor)
}
