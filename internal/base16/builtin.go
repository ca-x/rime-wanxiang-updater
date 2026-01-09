package base16

// Nord 配色方案
func Nord() *Scheme {
	return &Scheme{
		Scheme: "Nord",
		Author: "arcticicestudio",
		Base00: "2E3440",
		Base01: "3B4252",
		Base02: "434C5E",
		Base03: "4C566A",
		Base04: "D8DEE9",
		Base05: "E5E9F0",
		Base06: "ECEFF4",
		Base07: "8FBCBB",
		Base08: "BF616A",
		Base09: "D08770",
		Base0A: "EBCB8B",
		Base0B: "A3BE8C",
		Base0C: "88C0D0",
		Base0D: "81A1C1",
		Base0E: "B48EAD",
		Base0F: "5E81AC",
	}
}

// Dracula 配色方案
func Dracula() *Scheme {
	return &Scheme{
		Scheme: "Dracula",
		Author: "Zeno Rocha",
		Base00: "282a36",
		Base01: "44475a",
		Base02: "44475a",
		Base03: "6272a4",
		Base04: "f8f8f2",
		Base05: "f8f8f2",
		Base06: "f8f8f2",
		Base07: "f8f8f2",
		Base08: "ff5555",
		Base09: "ffb86c",
		Base0A: "f1fa8c",
		Base0B: "50fa7b",
		Base0C: "8be9fd",
		Base0D: "bd93f9",
		Base0E: "ff79c6",
		Base0F: "ff5555",
	}
}

// GruvboxDark 配色方案
func GruvboxDark() *Scheme {
	return &Scheme{
		Scheme: "Gruvbox Dark",
		Author: "morhetz",
		Base00: "282828",
		Base01: "3c3836",
		Base02: "504945",
		Base03: "665c54",
		Base04: "bdae93",
		Base05: "d5c4a1",
		Base06: "ebdbb2",
		Base07: "fbf1c7",
		Base08: "fb4934",
		Base09: "fe8019",
		Base0A: "fabd2f",
		Base0B: "b8bb26",
		Base0C: "8ec07c",
		Base0D: "83a598",
		Base0E: "d3869b",
		Base0F: "d65d0e",
	}
}

// GruvboxLight 配色方案
func GruvboxLight() *Scheme {
	return &Scheme{
		Scheme: "Gruvbox Light",
		Author: "morhetz",
		Base00: "fbf1c7",
		Base01: "ebdbb2",
		Base02: "d5c4a1",
		Base03: "bdae93",
		Base04: "665c54",
		Base05: "504945",
		Base06: "3c3836",
		Base07: "282828",
		Base08: "9d0006",
		Base09: "af3a03",
		Base0A: "b57614",
		Base0B: "79740e",
		Base0C: "427b58",
		Base0D: "076678",
		Base0E: "8f3f71",
		Base0F: "d65d0e",
	}
}

// Monokai 配色方案
func Monokai() *Scheme {
	return &Scheme{
		Scheme: "Monokai",
		Author: "Wimer Hazenberg",
		Base00: "272822",
		Base01: "383830",
		Base02: "49483e",
		Base03: "75715e",
		Base04: "a59f85",
		Base05: "f8f8f2",
		Base06: "f5f4f1",
		Base07: "f9f8f5",
		Base08: "f92672",
		Base09: "fd971f",
		Base0A: "f4bf75",
		Base0B: "a6e22e",
		Base0C: "a1efe4",
		Base0D: "66d9ef",
		Base0E: "ae81ff",
		Base0F: "cc6633",
	}
}

// TokyoNight 配色方案
func TokyoNight() *Scheme {
	return &Scheme{
		Scheme: "Tokyo Night",
		Author: "enkia",
		Base00: "1a1b26",
		Base01: "16161e",
		Base02: "2f3549",
		Base03: "444b6a",
		Base04: "787c99",
		Base05: "a9b1d6",
		Base06: "cbccd1",
		Base07: "d5d6db",
		Base08: "f7768e",
		Base09: "ff9e64",
		Base0A: "e0af68",
		Base0B: "9ece6a",
		Base0C: "7dcfff",
		Base0D: "7aa2f7",
		Base0E: "bb9af7",
		Base0F: "c0caf5",
	}
}

// TokyoNightLight 配色方案
func TokyoNightLight() *Scheme {
	return &Scheme{
		Scheme: "Tokyo Night Light",
		Author: "enkia",
		Base00: "d5d6db",
		Base01: "cbccd1",
		Base02: "dfe0e5",
		Base03: "9699a3",
		Base04: "4c505e",
		Base05: "343b59",
		Base06: "1a1b26",
		Base07: "1a1b26",
		Base08: "8c4351",
		Base09: "965027",
		Base0A: "8f5e15",
		Base0B: "485e30",
		Base0C: "166775",
		Base0D: "34548a",
		Base0E: "5a4a78",
		Base0F: "343b59",
	}
}

// CatppuccinMocha 配色方案
func CatppuccinMocha() *Scheme {
	return &Scheme{
		Scheme: "Catppuccin Mocha",
		Author: "Catppuccin",
		Base00: "1e1e2e",
		Base01: "181825",
		Base02: "313244",
		Base03: "45475a",
		Base04: "585b70",
		Base05: "cdd6f4",
		Base06: "f5e0dc",
		Base07: "b4befe",
		Base08: "f38ba8",
		Base09: "fab387",
		Base0A: "f9e2af",
		Base0B: "a6e3a1",
		Base0C: "94e2d5",
		Base0D: "89b4fa",
		Base0E: "cba6f7",
		Base0F: "f2cdcd",
	}
}

// CatppuccinLatte 配色方案（浅色）
func CatppuccinLatte() *Scheme {
	return &Scheme{
		Scheme: "Catppuccin Latte",
		Author: "Catppuccin",
		Base00: "eff1f5",
		Base01: "e6e9ef",
		Base02: "ccd0da",
		Base03: "bcc0cc",
		Base04: "acb0be",
		Base05: "4c4f69",
		Base06: "dc8a78",
		Base07: "7287fd",
		Base08: "d20f39",
		Base09: "fe640b",
		Base0A: "df8e1d",
		Base0B: "40a02b",
		Base0C: "179299",
		Base0D: "1e66f5",
		Base0E: "8839ef",
		Base0F: "dd7878",
	}
}

// OneDark 配色方案
func OneDark() *Scheme {
	return &Scheme{
		Scheme: "One Dark",
		Author: "Atom",
		Base00: "282c34",
		Base01: "353b45",
		Base02: "3e4451",
		Base03: "545862",
		Base04: "565c64",
		Base05: "abb2bf",
		Base06: "b6bdca",
		Base07: "c8ccd4",
		Base08: "e06c75",
		Base09: "d19a66",
		Base0A: "e5c07b",
		Base0B: "98c379",
		Base0C: "56b6c2",
		Base0D: "61afef",
		Base0E: "c678dd",
		Base0F: "be5046",
	}
}

// OneLight 配色方案
func OneLight() *Scheme {
	return &Scheme{
		Scheme: "One Light",
		Author: "Atom",
		Base00: "fafafa",
		Base01: "f0f0f1",
		Base02: "e5e5e6",
		Base03: "a0a1a7",
		Base04: "696c77",
		Base05: "383a42",
		Base06: "202227",
		Base07: "090a0b",
		Base08: "ca1243",
		Base09: "d75f00",
		Base0A: "c18401",
		Base0B: "50a14f",
		Base0C: "0184bc",
		Base0D: "4078f2",
		Base0E: "a626a4",
		Base0F: "986801",
	}
}

// SolarizedDark 配色方案
func SolarizedDark() *Scheme {
	return &Scheme{
		Scheme: "Solarized Dark",
		Author: "Ethan Schoonover",
		Base00: "002b36",
		Base01: "073642",
		Base02: "586e75",
		Base03: "657b83",
		Base04: "839496",
		Base05: "93a1a1",
		Base06: "eee8d5",
		Base07: "fdf6e3",
		Base08: "dc322f",
		Base09: "cb4b16",
		Base0A: "b58900",
		Base0B: "859900",
		Base0C: "2aa198",
		Base0D: "268bd2",
		Base0E: "6c71c4",
		Base0F: "d33682",
	}
}

// SolarizedLight 配色方案
func SolarizedLight() *Scheme {
	return &Scheme{
		Scheme: "Solarized Light",
		Author: "Ethan Schoonover",
		Base00: "fdf6e3",
		Base01: "eee8d5",
		Base02: "93a1a1",
		Base03: "839496",
		Base04: "657b83",
		Base05: "586e75",
		Base06: "073642",
		Base07: "002b36",
		Base08: "dc322f",
		Base09: "cb4b16",
		Base0A: "b58900",
		Base0B: "859900",
		Base0C: "2aa198",
		Base0D: "268bd2",
		Base0E: "6c71c4",
		Base0F: "d33682",
	}
}

// Cyberpunk 配色方案（赛博朋克风格，用于保持原有风格）
func Cyberpunk() *Scheme {
	return &Scheme{
		Scheme: "Cyberpunk",
		Author: "rime-wanxiang-updater",
		Base00: "0A0E27",
		Base01: "1A1F3A",
		Base02: "2A2F4A",
		Base03: "4C566A",
		Base04: "A9A9A9",
		Base05: "E5E9F0",
		Base06: "ECEFF4",
		Base07: "FFFFFF",
		Base08: "FF0040",
		Base09: "FF6600",
		Base0A: "FFFF00",
		Base0B: "00FF41",
		Base0C: "00FFFF",
		Base0D: "0080FF",
		Base0E: "FF00FF",
		Base0F: "B026FF",
	}
}

// CyberpunkLight 配色方案（赛博朋克风格浅色版）
func CyberpunkLight() *Scheme {
	return &Scheme{
		Scheme: "Cyberpunk Light",
		Author: "rime-wanxiang-updater",
		Base00: "F0F0F0",
		Base01: "E0E0E0",
		Base02: "D0D0D0",
		Base03: "A9A9A9",
		Base04: "505050",
		Base05: "303030",
		Base06: "202020",
		Base07: "000000",
		Base08: "DC143C",
		Base09: "FF4500",
		Base0A: "DAA520",
		Base0B: "008000",
		Base0C: "008B8B",
		Base0D: "0000CD",
		Base0E: "8B008B",
		Base0F: "6A0DAD",
	}
}

// ClassicLight 经典浅色主题
func ClassicLight() *Scheme {
	return &Scheme{
		Scheme: "Classic Light",
		Author: "Jason Heeris",
		Base00: "f5f5f5",
		Base01: "e0e0e0",
		Base02: "d0d0d0",
		Base03: "b0b0b0",
		Base04: "505050",
		Base05: "303030",
		Base06: "202020",
		Base07: "151515",
		Base08: "ac4142",
		Base09: "d28445",
		Base0A: "f4bf75",
		Base0B: "90a959",
		Base0C: "75b5aa",
		Base0D: "6a9fb5",
		Base0E: "aa759f",
		Base0F: "8f5536",
	}
}

// ClassicDark 经典深色主题
func ClassicDark() *Scheme {
	return &Scheme{
		Scheme: "Classic Dark",
		Author: "Jason Heeris",
		Base00: "151515",
		Base01: "202020",
		Base02: "303030",
		Base03: "505050",
		Base04: "b0b0b0",
		Base05: "d0d0d0",
		Base06: "e0e0e0",
		Base07: "f5f5f5",
		Base08: "ac4142",
		Base09: "d28445",
		Base0A: "f4bf75",
		Base0B: "90a959",
		Base0C: "75b5aa",
		Base0D: "6a9fb5",
		Base0E: "aa759f",
		Base0F: "8f5536",
	}
}

// BuiltinSchemes 返回所有内置配色方案
func BuiltinSchemes() map[string]*Scheme {
	return map[string]*Scheme{
		"nord":              Nord(),
		"dracula":           Dracula(),
		"gruvbox-dark":      GruvboxDark(),
		"gruvbox-light":     GruvboxLight(),
		"monokai":           Monokai(),
		"tokyo-night":       TokyoNight(),
		"tokyo-night-light": TokyoNightLight(),
		"catppuccin-mocha":  CatppuccinMocha(),
		"catppuccin-latte":  CatppuccinLatte(),
		"one-dark":          OneDark(),
		"one-light":         OneLight(),
		"solarized-dark":    SolarizedDark(),
		"solarized-light":   SolarizedLight(),
		"cyberpunk":         Cyberpunk(),
		"cyberpunk-light":   CyberpunkLight(),
		"classic-dark":      ClassicDark(),
		"classic-light":     ClassicLight(),
	}
}
