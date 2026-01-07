package config

import (
	"path/filepath"
	"regexp"
	"strings"
)

// DefaultExcludePatterns 默认排除文件模式
// 这些文件通常是用户自定义的配置，不应该被更新覆盖
var DefaultExcludePatterns = []string{
	`.*\.userdb$`,                     // 用户词库数据库
	`.*\.userdb\.txt$`,                // 用户词库文本
	`.*\.custom\.yaml$`,               // 用户自定义配置文件
	`^installation\.yaml$`,            // Rime 安装信息
	`^user\.yaml$`,                    // Rime 用户信息
	`^sync/.*`,                        // 同步目录下的所有文件
	`^build/.*`,                       // 构建目录下的所有文件
	`^custom/user_exclude_file\.txt$`, // 排除文件列表本身
}

// CommonExcludePatterns 常见的排除文件模式（用户可选）
var CommonExcludePatterns = map[string][]string{
	"用户数据": {
		`.*\.userdb$`,
		`.*\.userdb\.txt$`,
	},
	"自定义配置": {
		`.*\.custom\.yaml$`,
		`^default\.custom\.yaml$`,
	},
	"系统文件": {
		`^installation\.yaml$`,
		`^user\.yaml$`,
	},
	"临时文件": {
		`.*\.tmp$`,
		`.*\.bak$`,
		`.*~$`,
	},
	"同步和构建": {
		`^sync/.*`,
		`^build/.*`,
	},
}

// ExcludePatternType 排除模式类型
type ExcludePatternType int

const (
	PatternTypeWildcard ExcludePatternType = iota // 通配符模式 (*.yaml)
	PatternTypeRegex                              // 正则表达式模式 (^sync/.*)
	PatternTypeExact                              // 精确匹配模式 (user.yaml)
)

// ExcludePattern 排除模式结构
type ExcludePattern struct {
	Original string             // 原始模式字符串
	Type     ExcludePatternType // 模式类型
	Regex    *regexp.Regexp     // 编译后的正则表达式
}

// ParseExcludePattern 解析排除模式
// 支持三种格式：
//  1. 通配符: *.userdb, dicts/*.txt
//  2. 精确匹配: installation.yaml
//  3. 正则表达式: ^sync/.*$ (需要有正则特殊字符)
func ParseExcludePattern(pattern string) (*ExcludePattern, error) {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return nil, nil
	}

	ep := &ExcludePattern{
		Original: pattern,
	}

	// 检测模式类型
	// 优先检查正则表达式（因为正则可能包含 * 但不是通配符的含义）
	if hasRegexChars(pattern) {
		// 正则表达式模式
		ep.Type = PatternTypeRegex
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		ep.Regex = regex
	} else if strings.Contains(pattern, "*") || strings.Contains(pattern, "?") {
		// 通配符模式
		ep.Type = PatternTypeWildcard
		regexPattern := wildcardToRegex(pattern)
		regex, err := regexp.Compile(regexPattern)
		if err != nil {
			return nil, err
		}
		ep.Regex = regex
	} else {
		// 精确匹配模式
		ep.Type = PatternTypeExact
		// 对于精确匹配，也编译一个正则，用于统一处理
		regexPattern := "^" + regexp.QuoteMeta(pattern) + "$"
		regex, err := regexp.Compile(regexPattern)
		if err != nil {
			return nil, err
		}
		ep.Regex = regex
	}

	return ep, nil
}

// wildcardToRegex 将通配符模式转换为正则表达式
func wildcardToRegex(pattern string) string {
	// 转义所有正则特殊字符，除了 * 和 ?
	var result strings.Builder
	result.WriteString("^")

	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		switch c {
		case '*':
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				// ** 匹配任意层级目录
				result.WriteString(".*")
				i++ // 跳过下一个 *
			} else {
				// * 匹配单层目录内的任意字符（不包括路径分隔符）
				result.WriteString("[^/\\\\]*")
			}
		case '?':
			// ? 匹配单个字符（不包括路径分隔符）
			result.WriteString("[^/\\\\]")
		case '.', '+', '^', '$', '(', ')', '[', ']', '{', '}', '|', '\\':
			// 转义正则特殊字符
			result.WriteByte('\\')
			result.WriteByte(c)
		default:
			result.WriteByte(c)
		}
	}

	result.WriteString("$")
	return result.String()
}

// hasRegexChars 检查字符串是否包含正则表达式特殊字符
// 排除点号(.)，因为文件名中的点号很常见
func hasRegexChars(s string) bool {
	// 正则特殊字符（不包括 * 和 ? 因为它们被当作通配符处理）
	regexChars := []string{"^", "$", "[", "]", "(", ")", "{", "}", "|", "+", "\\"}
	for _, char := range regexChars {
		if strings.Contains(s, char) {
			return true
		}
	}
	// 特殊处理：如果有 \. 这种转义的点号，说明是正则
	if strings.Contains(s, `\.`) {
		return true
	}
	return false
}

// Match 检查文件路径是否匹配排除模式
func (ep *ExcludePattern) Match(filePath string) bool {
	if ep.Regex == nil {
		return false
	}

	// 标准化路径分隔符为正斜杠
	normalizedPath := filepath.ToSlash(filePath)

	// 同时检查完整路径和文件名
	baseName := filepath.Base(normalizedPath)

	return ep.Regex.MatchString(normalizedPath) || ep.Regex.MatchString(baseName)
}

// GetPatternDescription 获取模式的人类可读描述
func (ep *ExcludePattern) GetPatternDescription() string {
	switch ep.Type {
	case PatternTypeWildcard:
		return "通配符: " + ep.Original
	case PatternTypeRegex:
		return "正则: " + ep.Original
	case PatternTypeExact:
		return "精确: " + ep.Original
	default:
		return ep.Original
	}
}

// ParseExcludePatterns 批量解析排除模式
func ParseExcludePatterns(patterns []string) ([]*ExcludePattern, []error) {
	var result []*ExcludePattern
	var errors []error

	for _, pattern := range patterns {
		ep, err := ParseExcludePattern(pattern)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		if ep != nil {
			result = append(result, ep)
		}
	}

	return result, errors
}

// MatchAny 检查文件路径是否匹配任意一个排除模式
func MatchAny(filePath string, patterns []*ExcludePattern) bool {
	for _, pattern := range patterns {
		if pattern.Match(filePath) {
			return true
		}
	}
	return false
}
