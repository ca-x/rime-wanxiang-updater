package base16

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

// Scheme 表示一个 base16 配色方案
type Scheme struct {
	Scheme string `yaml:"scheme"`
	Author string `yaml:"author"`
	Base00 string `yaml:"base00"` // 默认背景
	Base01 string `yaml:"base01"` // 较浅背景（状态栏）
	Base02 string `yaml:"base02"` // 选择背景
	Base03 string `yaml:"base03"` // 注释、不可见字符
	Base04 string `yaml:"base04"` // 深色前景（状态栏）
	Base05 string `yaml:"base05"` // 默认前景
	Base06 string `yaml:"base06"` // 较浅前景
	Base07 string `yaml:"base07"` // 最浅前景
	Base08 string `yaml:"base08"` // 变量、删除、错误
	Base09 string `yaml:"base09"` // 整数、布尔、常量
	Base0A string `yaml:"base0A"` // 类、搜索高亮
	Base0B string `yaml:"base0B"` // 字符串、成功
	Base0C string `yaml:"base0C"` // 正则、转义、引用
	Base0D string `yaml:"base0D"` // 函数、方法
	Base0E string `yaml:"base0E"` // 关键字、存储
	Base0F string `yaml:"base0F"` // 弃用、嵌入语言
}

// LoadFromYAML 从 YAML 文件加载配色方案
func LoadFromYAML(path string) (*Scheme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var scheme Scheme
	if err := yaml.Unmarshal(data, &scheme); err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	return &scheme, nil
}

// SaveToYAML 保存配色方案到 YAML 文件
func (s *Scheme) SaveToYAML(path string) error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal scheme: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Color 返回指定 base 的 lipgloss.Color
func (s *Scheme) Color(base int) lipgloss.Color {
	hex := s.Hex(base)
	return lipgloss.Color("#" + hex)
}

// Hex 返回指定 base 的十六进制颜色值（不含 #）
func (s *Scheme) Hex(base int) string {
	switch base {
	case 0x00:
		return s.Base00
	case 0x01:
		return s.Base01
	case 0x02:
		return s.Base02
	case 0x03:
		return s.Base03
	case 0x04:
		return s.Base04
	case 0x05:
		return s.Base05
	case 0x06:
		return s.Base06
	case 0x07:
		return s.Base07
	case 0x08:
		return s.Base08
	case 0x09:
		return s.Base09
	case 0x0A:
		return s.Base0A
	case 0x0B:
		return s.Base0B
	case 0x0C:
		return s.Base0C
	case 0x0D:
		return s.Base0D
	case 0x0E:
		return s.Base0E
	case 0x0F:
		return s.Base0F
	default:
		return s.Base05
	}
}

// AdaptiveColor 创建自适应颜色（用于明暗模式）
func AdaptiveColor(light, dark *Scheme, base int) lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: "#" + light.Hex(base),
		Dark:  "#" + dark.Hex(base),
	}
}
