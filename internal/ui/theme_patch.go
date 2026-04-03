package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/deployer"
	"rime-wanxiang-updater/internal/types"

	"gopkg.in/yaml.v3"
)

type themePatchDefinition struct {
	Key         string
	DisplayName string
	Values      map[string]any
}

var deployThemePatch = func(cfg *types.Config) error {
	return deployer.GetDeployer(cfg).Deploy()
}

func (m Model) themePatchFilePath() (string, error) {
	engineName, fileName, ok := themePatchTargetForPlatform(runtime.GOOS, m.Cfg.Config.InstalledEngines)
	if !ok {
		return "", fmt.Errorf("当前平台没有可用的主题 patch 目标")
	}

	rootDir := config.GetEngineDataDir(engineName)
	if rootDir == "" {
		return "", fmt.Errorf("未找到 %s 的用户目录", engineName)
	}

	return filepath.Join(rootDir, fileName), nil
}

func normalizeStringMap(value any) (map[string]any, bool) {
	switch typed := value.(type) {
	case map[string]any:
		return typed, true
	case map[any]any:
		result := make(map[string]any, len(typed))
		for key, entry := range typed {
			keyString, ok := key.(string)
			if !ok {
				return nil, false
			}
			result[keyString] = entry
		}
		return result, true
	default:
		return nil, false
	}
}

func readThemePatchDocument(path string) (map[string]any, error) {
	document := make(map[string]any)

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("读取现有 patch 文件失败: %w", err)
	}

	if len(data) > 0 {
		if err := yaml.Unmarshal(data, &document); err != nil {
			return nil, fmt.Errorf("解析现有 patch 文件失败: %w", err)
		}
	}

	return document, nil
}

func themePatchSection(document map[string]any) map[string]any {
	patch, ok := normalizeStringMap(document["patch"])
	if !ok || patch == nil {
		patch = map[string]any{}
	}

	return patch
}

func writeThemePatchDocument(path string, document map[string]any) error {
	encoded, err := yaml.Marshal(document)
	if err != nil {
		return fmt.Errorf("序列化主题 patch 失败: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("创建主题 patch 目录失败: %w", err)
	}

	if err := os.WriteFile(path, encoded, 0644); err != nil {
		return fmt.Errorf("写入主题 patch 文件失败: %w", err)
	}

	return nil
}

func writeThemePatchPresets(path string, definitions []themePatchDefinition) error {
	document, err := readThemePatchDocument(path)
	if err != nil {
		return err
	}

	patch := themePatchSection(document)
	for _, definition := range definitions {
		patch["preset_color_schemes/"+definition.Key] = definition.Values
	}
	document["patch"] = patch

	return writeThemePatchDocument(path, document)
}

func readThemePatchSelections(path string) (map[string]bool, error) {
	document, err := readThemePatchDocument(path)
	if err != nil {
		return nil, err
	}

	patch := themePatchSection(document)
	selections := make(map[string]bool)
	for key := range patch {
		if !strings.HasPrefix(key, "preset_color_schemes/") {
			continue
		}

		themeKey := strings.TrimPrefix(key, "preset_color_schemes/")
		if _, ok := themePatchDefinitionByKey(themeKey); ok {
			selections[themeKey] = true
		}
	}

	return selections, nil
}

func syncThemePatchPresets(path string, selections map[string]bool) error {
	document, err := readThemePatchDocument(path)
	if err != nil {
		return err
	}

	patch := themePatchSection(document)
	for _, definition := range themePatchDefinitions() {
		key := "preset_color_schemes/" + definition.Key
		if selections[definition.Key] {
			patch[key] = definition.Values
			continue
		}
		delete(patch, key)
	}

	document["patch"] = patch
	return writeThemePatchDocument(path, document)
}

func writeThemePatchDefault(path, themeKey string) error {
	document, err := readThemePatchDocument(path)
	if err != nil {
		return err
	}

	patch := themePatchSection(document)
	patch["style/color_scheme"] = themeKey
	patch["style/color_scheme_dark"] = themeKey
	document["patch"] = patch

	return writeThemePatchDocument(path, document)
}

func themePatchDefinitionByKey(key string) (themePatchDefinition, bool) {
	for _, definition := range themePatchDefinitions() {
		if definition.Key == key {
			return definition, true
		}
	}

	return themePatchDefinition{}, false
}

func themePatchDefinitions() []themePatchDefinition {
	return []themePatchDefinition{
		{
			Key:         "jianchun",
			DisplayName: "简纯",
			Values: map[string]any{
				"name":                         "简纯",
				"author":                       "amzxyz",
				"back_color":                   "0xf2f2f2",
				"border_color":                 "0xCE7539",
				"text_color":                   "0x3c647e",
				"hilited_text_color":           "0x3c647e",
				"hilited_back_color":           "0x797954",
				"hilited_comment_text_color":   "0xffffff",
				"hilited_candidate_text_color": "0xffffff",
				"hilited_candidate_back_color": "0xCE7539",
				"hilited_label_color":          "0xdedede",
				"candidate_text_color":         "0x000000",
				"comment_text_color":           "0x000000",
				"label_color":                  "0x91897e",
			},
		},
		{
			Key:         "win11_light",
			DisplayName: "Win11浅色 / Win11light",
			Values: map[string]any{
				"name":                         "Win11浅色 / Win11light",
				"text_color":                   "0x191919",
				"label_color":                  "0x191919",
				"hilited_label_color":          "0x191919",
				"back_color":                   "0xf9f9f9",
				"border_color":                 "0x009e5a00",
				"hilited_mark_color":           "0xc06700",
				"hilited_candidate_back_color": "0xf0f0f0",
				"shadow_color":                 "0x20000000",
			},
		},
		{
			Key:         "win11_dark",
			DisplayName: "Win11暗色 / Win11Dark",
			Values: map[string]any{
				"name":                         "Win11暗色 / Win11Dark",
				"text_color":                   "0xf9f9f9",
				"label_color":                  "0xf9f9f9",
				"back_color":                   "0x2C2C2C",
				"hilited_label_color":          "0xf9f9f9",
				"border_color":                 "0x002C2C2C",
				"hilited_mark_color":           "0xFFC24C",
				"hilited_candidate_back_color": "0x383838",
				"shadow_color":                 "0x20000000",
			},
		},
		{
			Key:         "mac_light",
			DisplayName: "Mac 白",
			Values: map[string]any{
				"name":                         "Mac 白",
				"text_color":                   "0x000000",
				"back_color":                   "0xffffff",
				"border_color":                 "0xe9e9e9",
				"label_color":                  "0x999999",
				"hilited_text_color":           "0x000000",
				"hilited_back_color":           "0xffffff",
				"candidate_text_color":         "0x000000",
				"comment_text_color":           "0x999999",
				"hilited_candidate_text_color": "0xffffff",
				"hilited_comment_text_color":   "0xdddddd",
				"hilited_candidate_back_color": 16740656,
				"hilited_label_color":          "0xffffff",
			},
		},
		{
			Key:         "wechat",
			DisplayName: "微信 / Wechat",
			Values: map[string]any{
				"name":                         "微信／Wechat",
				"text_color":                   "0x424242",
				"label_color":                  "0x999999",
				"back_color":                   "0xFFFFFF",
				"border_color":                 "0xFFFFFF",
				"comment_text_color":           "0x999999",
				"candidate_text_color":         "0x3c3c3c",
				"hilited_comment_text_color":   "0xFFFFFF",
				"hilited_back_color":           "0x79af22",
				"hilited_text_color":           "0xFFFFFF",
				"hilited_label_color":          "0xFFFFFF",
				"hilited_candidate_back_color": "0x79af22",
				"shadow_color":                 "0x20000000",
			},
		},
		{
			Key:         "Lumk_light",
			DisplayName: "鹿鸣 / Lumk light",
			Values: map[string]any{
				"name":                          "鹿鸣／Lumk light",
				"author":                        "Lumk X <x@xx.cc>",
				"back_color":                    "0xF9F9F9",
				"border_color":                  "0xE2E7F5",
				"candidate_text_color":          "0x121212",
				"comment_text_color":            "0x8E8E8E",
				"hilited_candidate_back_color":  "0xECE4FC",
				"hilited_candidate_label_color": "0xB18FF4",
				"hilited_candidate_text_color":  "0x7A40EC",
				"hilited_label_color":           "0xA483EC",
				"hilited_mark_color":            "0x7A40EC",
				"label_color":                   "0x888785",
				"text_color":                    "0x8100EB",
				"shadow_color":                  "0x20000000",
			},
		},
		{
			Key:         "amber-7",
			DisplayName: "淡白 / weasel",
			Values: map[string]any{
				"name":                           "淡白/weasel",
				"author":                         "五笔小筑 <wubixiaozhu@126.com>",
				"back_color":                     "0xffffff",
				"border_color":                   "0xE99321",
				"shadow_color":                   "0x00000000",
				"text_color":                     "0xE99321",
				"hilited_text_color":             "0x2238dc",
				"hilited_back_color":             "0xffffff",
				"hilited_shadow_color":           "0x00000000",
				"nextpage_color":                 "0x0000FF",
				"prevpage_color":                 "0x0000FF",
				"hilited_mark_color":             "0x00000000",
				"hilited_label_color":            "0x2021FF",
				"hilited_candidate_text_color":   "0x2021FF",
				"hilited_comment_text_color":     "0x000000",
				"hilited_candidate_back_color":   "0xECF1FC",
				"hilited_candidate_border_color": "0xCFDCFD",
				"hilited_candidate_shadow_color": "0x00000000",
				"label_color":                    "0xE99321",
				"candidate_text_color":           "0xE99321",
				"comment_text_color":             "0x000000",
				"candidate_back_color":           "0xffffff",
				"candidate_border_color":         "0xffffff",
				"candidate_shadow_color":         "0x00000000",
			},
		},
		{
			Key:         "win10gray",
			DisplayName: "win10灰 / win10gray",
			Values: map[string]any{
				"name":                           "win10灰／win10gray",
				"author":                         "五笔小筑 <wubixiaozhu@126.com>",
				"back_color":                     "0xfff4f4f4",
				"shadow_color":                   "0xf7606060",
				"border_color":                   "0xff305c3c",
				"text_color":                     "0xff000000",
				"hilited_text_color":             "0xff000000",
				"hilited_back_color":             "0x4ff4f4f4",
				"hilited_shadow_color":           "0x00000000",
				"hilited_mark_color":             "0xffd77800",
				"hilited_label_color":            "0xff555555",
				"hilited_candidate_text_color":   "0xff000000",
				"hilited_comment_text_color":     "0xff555555",
				"hilited_candidate_back_color":   "0x4fcccccc",
				"hilited_candidate_border_color": "0x00000000",
				"hilited_candidate_shadow_color": "0x00000000",
				"label_color":                    "0xff888888",
				"candidate_text_color":           "0xff222222",
				"comment_text_color":             "0xff888888",
				"candidate_back_color":           "0x00000000",
				"candidate_border_color":         "0x00000000",
				"candidate_shadow_color":         "0x00000000",
			},
		},
		{
			Key:         "mint_light_blue",
			DisplayName: "蓝水鸭 / Mint Light Blue",
			Values: map[string]any{
				"name":                          "蓝水鸭／Mint Light Blue",
				"author":                        "Mintimate <Mintimate's Blog>",
				"text_color":                    "0x6495ed",
				"back_color":                    "0xefefef",
				"label_color":                   "0xcac9c8",
				"border_color":                  "0xefefef",
				"shadow_color":                  "0xb4000000",
				"comment_text_color":            "0xcac9c8",
				"candidate_text_color":          "0x424242",
				"hilited_text_color":            "0xed9564",
				"hilited_back_color":            "0xefefef",
				"hilited_candidate_back_color":  "0xed9564",
				"hilited_candidate_text_color":  "0xefefef",
				"hilited_candidate_label_color": "0xcac9c8",
				"hilited_label_color":           "0xcac9c8",
				"hilited_comment_text_color":    "0xefefef",
				"nextpage_color":                "0x0000FF",
				"prevpage_color":                "0x0000FF",
			},
		},
		{
			Key:         "ayaya",
			DisplayName: "文文 / Ayaya",
			Values: map[string]any{
				"name":                          "文文／Ayaya",
				"author":                        "Lufs X <i@isteed.cc>",
				"font_face":                     "LXGWWenKai-Regular, PingFangSC",
				"font_point":                    16.5,
				"label_font_face":               "LXGWWenKai-Regular, PingFangSC",
				"label_font_point":              12,
				"candidate_format":              "[label]\u2005[candidate] [comment]",
				"candidate_list_layout":         "linear",
				"text_orientation":              "horizontal",
				"inline_preedit":                true,
				"corner_radius":                 5,
				"hilited_corner_radius":         0,
				"border_height":                 0,
				"border_width":                  0,
				"alpha":                         0.95,
				"shadow_size":                   0,
				"color_space":                   "display_p3",
				"back_color":                    "0xFFFFFF",
				"border_color":                  "0xECE4FC",
				"candidate_text_color":          "0x121212",
				"comment_text_color":            "0x8E8E8E",
				"label_color":                   "0x888785",
				"hilited_candidate_back_color":  "0xECE4FC",
				"hilited_candidate_text_color":  "0x7A40EC",
				"hilited_comment_text_color":    "0x8E8E8E",
				"hilited_candidate_label_color": "0xB18FF4",
				"text_color":                    "0x8100EB",
				"hilited_text_color":            "0xD8000000",
			},
		},
		{
			Key:         "ayaya_dark",
			DisplayName: "文文 / Ayaya / 深色",
			Values: map[string]any{
				"name":                          "文文／Ayaya／深色",
				"author":                        "Lufs X <i@isteed.cc>",
				"font_face":                     "LXGWWenKai-Regular, PingFangSC",
				"font_point":                    16.5,
				"label_font_face":               "LXGWWenKai-Regular, PingFangSC",
				"label_font_point":              12,
				"candidate_format":              "[label]\u2005[candidate] [comment]",
				"candidate_list_layout":         "linear",
				"text_orientation":              "horizontal",
				"inline_preedit":                true,
				"corner_radius":                 5,
				"hilited_corner_radius":         0,
				"border_height":                 0,
				"border_width":                  0,
				"alpha":                         0.95,
				"shadow_size":                   0,
				"color_space":                   "display_p3",
				"back_color":                    "0x000000",
				"border_color":                  "0xECE4FC",
				"candidate_text_color":          "0xD2D2D2",
				"comment_text_color":            "0x8E8E8E",
				"label_color":                   "0x888785",
				"hilited_candidate_back_color":  "0x2C1E3C",
				"hilited_candidate_text_color":  "0x7036E2",
				"hilited_comment_text_color":    "0x8E8E8E",
				"hilited_candidate_label_color": "0x7036E2",
				"text_color":                    "0x8100EB",
				"hilited_text_color":            "0xD8000000",
			},
		},
		{
			Key:         "reimu",
			DisplayName: "灵梦 / Reimu",
			Values: map[string]any{
				"name":                          "灵梦／Reimu",
				"author":                        "Lufs X <i@isteed.cc>",
				"font_face":                     "LXGWWenKai-Regular, PingFangSC",
				"font_point":                    17,
				"label_font_face":               "LXGWWenKai-Regular, PingFangSC",
				"label_font_point":              14,
				"candidate_format":              "[label]\u2005[candidate] [comment]",
				"candidate_list_layout":         "linear",
				"text_orientation":              "horizontal",
				"inline_preedit":                true,
				"corner_radius":                 7,
				"hilited_corner_radius":         6,
				"border_height":                 1,
				"border_width":                  1,
				"alpha":                         0.95,
				"shadow_size":                   2,
				"color_space":                   "display_p3",
				"back_color":                    "0xF5FCFD",
				"candidate_text_color":          "0x282C32",
				"comment_text_color":            "0x717172",
				"label_color":                   "0x888785",
				"hilited_candidate_back_color":  "0xF5FCFD",
				"hilited_candidate_text_color":  "0x4F00E5",
				"hilited_comment_text_color":    "0x9F9CF2",
				"hilited_candidate_label_color": "0x4F00E5",
				"text_color":                    "0x6B54E9",
				"hilited_text_color":            "0xD8000000",
			},
		},
		{
			Key:         "reimu_dark",
			DisplayName: "灵梦 / Reimu / 深色",
			Values: map[string]any{
				"name":                          "灵梦／Reimu／深色",
				"author":                        "Lufs X <i@isteed.cc>",
				"font_face":                     "LXGWWenKai-Regular, PingFangSC",
				"font_point":                    17,
				"label_font_face":               "LXGWWenKai-Regular, PingFangSC",
				"label_font_point":              14,
				"candidate_format":              "[label]\u2005[candidate] [comment]",
				"candidate_list_layout":         "linear",
				"text_orientation":              "horizontal",
				"inline_preedit":                true,
				"corner_radius":                 7,
				"hilited_corner_radius":         6,
				"border_height":                 1,
				"border_width":                  1,
				"alpha":                         0.95,
				"shadow_size":                   2,
				"color_space":                   "display_p3",
				"back_color":                    "0x020A00",
				"border_color":                  "0x020A00",
				"candidate_text_color":          "0xC0C0C0",
				"comment_text_color":            "0x717172",
				"label_color":                   "0x717172",
				"hilited_candidate_back_color":  "0x0C140A",
				"hilited_candidate_text_color":  "0x3100C7",
				"hilited_comment_text_color":    "0x7772AF",
				"hilited_candidate_label_color": "0x3100C7",
				"text_color":                    "0x6B54E9",
				"hilited_text_color":            "0xD8000000",
			},
		},
		{
			Key:         "apathy",
			DisplayName: "冷漠 / Apathy",
			Values: map[string]any{
				"name":                          "冷漠／Apathy",
				"author":                        "LIANG Hai",
				"font_face":                     "PingFangSC-Regular",
				"font_point":                    16,
				"label_font_face":               "STHeitiSC-Light",
				"label_font_point":              12,
				"candidate_format":              "[label]\u2005[candidate] [comment]",
				"candidate_list_layout":         "linear",
				"text_orientation":              "horizontal",
				"inline_preedit":                true,
				"corner_radius":                 5,
				"alpha":                         0.95,
				"color_space":                   "srgb",
				"back_color":                    "0xFFFFFF",
				"candidate_text_color":          "0xD8000000",
				"comment_text_color":            "0x999999",
				"label_color":                   "0xE5555555",
				"hilited_candidate_back_color":  "0xFFF0E4",
				"hilited_candidate_text_color":  "0xEE6E00",
				"hilited_comment_text_color":    "0x999999",
				"hilited_candidate_label_color": "0xF4994C",
				"text_color":                    "0x424242",
				"hilited_text_color":            "0xD8000000",
			},
		},
		{
			Key:         "win10",
			DisplayName: "WIN10",
			Values: map[string]any{
				"name":                          "WIN10",
				"author":                        "Lufs X <i@isteed.cc>",
				"font_face":                     "LXGWWenKai-Regular, PingFangSC",
				"font_point":                    16.5,
				"label_font_face":               "LXGWWenKai-Regular, PingFangSC",
				"label_font_point":              12,
				"candidate_list_layout":         "linear",
				"text_orientation":              "horizontal",
				"inline_preedit":                true,
				"corner_radius":                 0,
				"hilited_corner_radius":         -6,
				"border_height":                 7,
				"border_width":                  6,
				"spacing":                       10,
				"color_space":                   "srgb",
				"back_color":                    "0xFFFFFF",
				"candidate_text_color":          "0x000000",
				"comment_text_color":            "0x888888",
				"label_color":                   "0x888888",
				"hilited_candidate_back_color":  "0xCC8F29",
				"hilited_candidate_text_color":  "0xFFFFFF",
				"hilited_comment_text_color":    "0xFFFFFF",
				"hilited_candidate_label_color": "0xEEDAB8",
				"text_color":                    "0x000000",
				"hilited_back_color":            "0xFFFFFF",
				"hilited_text_color":            "0x000000",
			},
		},
		{
			Key:         "win10_ayaya",
			DisplayName: "WIN10 / 文文 / Ayaya",
			Values: map[string]any{
				"name":                          "WIN10／文文／Ayaya",
				"author":                        "Lufs X <i@isteed.cc>",
				"font_face":                     "LXGWWenKai-Regular, PingFangSC",
				"font_point":                    16.5,
				"label_font_face":               "LXGWWenKai-Regular, PingFangSC",
				"label_font_point":              12,
				"candidate_list_layout":         "linear",
				"text_orientation":              "horizontal",
				"inline_preedit":                true,
				"corner_radius":                 0,
				"hilited_corner_radius":         -6,
				"border_height":                 7,
				"border_width":                  6,
				"spacing":                       10,
				"color_space":                   "display_p3",
				"back_color":                    "0xFFFFFF",
				"border_color":                  "0xFFFFFF",
				"candidate_text_color":          "0x121212",
				"comment_text_color":            "0x8E8E8E",
				"label_color":                   "0x888785",
				"hilited_candidate_back_color":  "0xECE4FC",
				"hilited_candidate_text_color":  "0x7A40EC",
				"hilited_comment_text_color":    "0x8E8E8E",
				"hilited_candidate_label_color": "0xB18FF4",
				"text_color":                    "0x8100EB",
				"hilited_text_color":            "0xD8000000",
			},
		},
		{
			Key:         "macos12_light",
			DisplayName: "高仿亮色 macOS",
			Values: map[string]any{
				"name":                          "高仿亮色 macOS",
				"author":                        "Lufs X <i@isteed.cc>",
				"font_face":                     "PingFangSC-Regular",
				"font_point":                    15,
				"label_font_face":               "PingFangSC-Regular",
				"label_font_point":              12,
				"comment_font_face":             "PingFangSC-Regular",
				"comment_font_point":            13,
				"candidate_format":              "[label]\u2004[candidate] [comment]",
				"candidate_list_layout":         "linear",
				"text_orientation":              "horizontal",
				"inline_preedit":                true,
				"corner_radius":                 7,
				"hilited_corner_radius":         6,
				"border_width":                  2,
				"line_spacing":                  1,
				"color_space":                   "display_p3",
				"back_color":                    "0xFFFFFF",
				"border_color":                  "0xFFFFFF",
				"candidate_text_color":          "0xD8000000",
				"comment_text_color":            "0x3F000000",
				"label_color":                   "0x7F7F7F",
				"hilited_candidate_back_color":  "0xD05925",
				"hilited_candidate_text_color":  "0xFFFFFF",
				"hilited_comment_text_color":    "0x808080",
				"hilited_candidate_label_color": "0xFFFFFF",
				"text_color":                    "0x3F000000",
				"hilited_text_color":            "0xD8000000",
			},
		},
		{
			Key:         "macos12_dark",
			DisplayName: "高仿暗色 macOS",
			Values: map[string]any{
				"name":                          "高仿暗色 macOS",
				"author":                        "Lufs X <i@isteed.cc>",
				"font_face":                     "PingFangSC-Regular",
				"font_point":                    15,
				"label_font_face":               "PingFangSC-Regular",
				"label_font_point":              12,
				"comment_font_face":             "PingFangSC-Regular",
				"comment_font_point":            13,
				"candidate_format":              "[label]\u2004[candidate] [comment]",
				"candidate_list_layout":         "linear",
				"text_orientation":              "horizontal",
				"inline_preedit":                true,
				"corner_radius":                 7,
				"hilited_corner_radius":         6,
				"border_width":                  2,
				"line_spacing":                  1,
				"color_space":                   "display_p3",
				"back_color":                    "0x1E1F24",
				"border_color":                  "0x1E1F24",
				"candidate_text_color":          "0xE8E8E8",
				"comment_text_color":            "0x3F000000",
				"label_color":                   "0x7C7C7C",
				"hilited_candidate_back_color":  "0xDA6203",
				"hilited_candidate_text_color":  "0xFFFFFF",
				"hilited_comment_text_color":    "0x808080",
				"hilited_candidate_label_color": "0xFFE7D6",
				"text_color":                    "0x3F000000",
				"hilited_text_color":            "0xD8000000",
			},
		},
		{
			Key:         "wechat_dark",
			DisplayName: "高仿暗色微信输入法",
			Values: map[string]any{
				"name":                          "高仿暗色微信输入法",
				"author":                        "Lufs X <i@isteed.cc>",
				"font_face":                     "PingFangSC-Regular",
				"font_point":                    15,
				"label_font_face":               "PingFangSC-Regular",
				"label_font_point":              12,
				"comment_font_face":             "PingFangSC-Regular",
				"comment_font_point":            12,
				"candidate_format":              "[label]\u2005[candidate] [comment]",
				"candidate_list_layout":         "linear",
				"text_orientation":              "horizontal",
				"inline_preedit":                true,
				"corner_radius":                 7,
				"hilited_corner_radius":         7,
				"border_height":                 -2,
				"color_space":                   "display_p3",
				"back_color":                    "0x151515",
				"border_color":                  "0x151515",
				"candidate_text_color":          "0xB9B9B9",
				"comment_text_color":            "0x8E8E8E",
				"label_color":                   "0x888785",
				"hilited_candidate_back_color":  "0x74A54B",
				"hilited_candidate_text_color":  "0xFFFFFF",
				"hilited_comment_text_color":    "0xF0F0F0",
				"hilited_candidate_label_color": "0xFFFFFF",
				"text_color":                    "0xFFFFFF",
				"hilited_text_color":            "0x777777",
			},
		},
		{
			Key:         "macos14",
			DisplayName: "高仿 macOS 14",
			Values: map[string]any{
				"name":                          "高仿 macOS 14",
				"author":                        "Lufs X <i@isteed.cc>",
				"font_face":                     "PingFangSC-Regular",
				"font_point":                    16,
				"label_font_face":               "PingFangSC-Regular",
				"label_font_point":              8,
				"comment_font_face":             "PingFangSC-Regular",
				"comment_font_point":            13,
				"candidate_format":              "[label]\u2003\u2003\u2003[candidate] [comment]\u2004\u2001",
				"candidate_list_layout":         "linear",
				"text_orientation":              "horizontal",
				"inline_preedit":                true,
				"corner_radius":                 7,
				"hilited_corner_radius":         -1,
				"border_height":                 -4,
				"border_width":                  2,
				"color_space":                   "display_p3",
				"back_color":                    "0xE7E8EA",
				"border_color":                  "0xE7E8EA",
				"candidate_text_color":          "0x464647",
				"comment_text_color":            "0x3F000000",
				"label_color":                   "0x7F7F7F",
				"hilited_candidate_back_color":  "0xD05925",
				"hilited_candidate_text_color":  "0xFFFFFF",
				"hilited_comment_text_color":    "0xDCDCDC",
				"hilited_candidate_label_color": "0xFFFFFF",
				"text_color":                    "0x3F000000",
				"hilited_text_color":            "0xD8000000",
			},
		},
		{
			Key:         "macos14_dark",
			DisplayName: "高仿暗色 macOS 14",
			Values: map[string]any{
				"name":                          "高仿暗色 macOS 14",
				"author":                        "Lufs X <i@isteed.cc>",
				"font_face":                     "PingFangSC-Regular",
				"font_point":                    16,
				"label_font_face":               "PingFangSC-Regular",
				"label_font_point":              8,
				"comment_font_face":             "PingFangSC-Regular",
				"comment_font_point":            13,
				"candidate_format":              "[label]\u2003\u2003\u2003[candidate] [comment]\u2004\u2001",
				"candidate_list_layout":         "linear",
				"text_orientation":              "horizontal",
				"inline_preedit":                true,
				"corner_radius":                 7,
				"hilited_corner_radius":         -1,
				"border_height":                 -4,
				"border_width":                  2,
				"color_space":                   "display_p3",
				"back_color":                    "0x555557",
				"border_color":                  "0x0C0C0C",
				"candidate_text_color":          "0xEEEEEE",
				"comment_text_color":            "0x80FFFFFF",
				"label_color":                   "0x7C7C7C",
				"hilited_candidate_back_color":  "0xCA5824",
				"hilited_candidate_text_color":  "0xFFFFFF",
				"hilited_comment_text_color":    "0xF0F0F0",
				"hilited_candidate_label_color": "0xFFFFFF",
				"text_color":                    "0x3F000000",
				"hilited_text_color":            "0xD8000000",
			},
		},
	}
}
