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

package metadata

var (
	CEFFrameworkOptionCEFConfig string
)

var (
	CEFFormStatusLabelCaption                      string
	CEFFormStatusLabelCaptionInvalid               string
	CEFFormStatusLabelCaptionSelectCEF             string
	CEFFormStatusLabelCaptionURLNotFound           string
	CEFFormStatusLabelCaptionFailedCreateDirectory string
	CEFFormStatusLabelCaptionCaptionFailedPreDown  string
	CEFFormStatusLabelCaptionCaptionPaused         string
	CEFFormOsLabelText                             string
	CEFFormArchLabelText                           string
)

func init() {
	GI18n.Bind("MenuSettingCEFFrameworkOption.Caption", &CEFFrameworkOptionCEFConfig)

	GI18n.Bind("ChromiumDirFormStatusLabel.Caption", &CEFFormStatusLabelCaption)
	GI18n.Bind("ChromiumDirFormStatusLabel.CaptionInvalid", &CEFFormStatusLabelCaptionInvalid)
	GI18n.Bind("ChromiumDirFormStatusLabel.CaptionSelectCEF", &CEFFormStatusLabelCaptionSelectCEF)
	GI18n.Bind("ChromiumDirFormStatusLabel.CaptionURLNotFound", &CEFFormStatusLabelCaptionURLNotFound)
	GI18n.Bind("ChromiumDirFormStatusLabel.CaptionFailedCreateDirectory", &CEFFormStatusLabelCaptionFailedCreateDirectory)
	GI18n.Bind("ChromiumDirFormStatusLabel.CaptionFailedPreDown", &CEFFormStatusLabelCaptionCaptionFailedPreDown)
	GI18n.Bind("ChromiumDirFormStatusLabel.CaptionPaused", &CEFFormStatusLabelCaptionCaptionPaused)
	GI18n.Bind("ChromiumDirFormOSText.Caption", &CEFFormOsLabelText)
	GI18n.Bind("ChromiumDirFormARCHText.Caption", &CEFFormArchLabelText)
}
