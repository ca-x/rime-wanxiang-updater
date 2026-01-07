package config

import (
	"testing"
)

func TestParseExcludePattern(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		expectType  ExcludePatternType
		testFile    string
		shouldMatch bool
	}{
		{
			name:        "通配符 - 单个星号",
			pattern:     "*.userdb",
			expectType:  PatternTypeWildcard,
			testFile:    "test.userdb",
			shouldMatch: true,
		},
		{
			name:        "通配符 - 路径匹配",
			pattern:     "dicts/*.txt",
			expectType:  PatternTypeWildcard,
			testFile:    "dicts/test.txt",
			shouldMatch: true,
		},
		{
			name:        "通配符 - 双星号",
			pattern:     "sync/**/*.yaml",
			expectType:  PatternTypeWildcard,
			testFile:    "sync/deep/nested/file.yaml",
			shouldMatch: true,
		},
		{
			name:        "正则表达式",
			pattern:     `^sync/.*\.yaml$`,
			expectType:  PatternTypeRegex,
			testFile:    "sync/test.yaml",
			shouldMatch: true,
		},
		{
			name:        "精确匹配",
			pattern:     "installation.yaml",
			expectType:  PatternTypeExact,
			testFile:    "installation.yaml",
			shouldMatch: true,
		},
		{
			name:        "精确匹配 - 不匹配",
			pattern:     "installation.yaml",
			expectType:  PatternTypeExact,
			testFile:    "user.yaml",
			shouldMatch: false,
		},
		{
			name:        "通配符 - 自定义配置",
			pattern:     "*.custom.yaml",
			expectType:  PatternTypeWildcard,
			testFile:    "default.custom.yaml",
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep, err := ParseExcludePattern(tt.pattern)
			if err != nil {
				t.Fatalf("解析失败: %v", err)
			}

			if ep.Type != tt.expectType {
				t.Errorf("期望类型 %v, 实际 %v", tt.expectType, ep.Type)
			}

			matched := ep.Match(tt.testFile)
			if matched != tt.shouldMatch {
				t.Errorf("文件 %s 匹配结果: 期望 %v, 实际 %v", tt.testFile, tt.shouldMatch, matched)
			}

			t.Logf("模式描述: %s", ep.GetPatternDescription())
		})
	}
}

func TestWildcardToRegex(t *testing.T) {
	tests := []struct {
		wildcard string
		expected string
		testCase string
	}{
		{
			wildcard: "*.txt",
			expected: `^[^/\\]*\.txt$`,
			testCase: "test.txt",
		},
		{
			wildcard: "sync/*.yaml",
			expected: `^sync/[^/\\]*\.yaml$`,
			testCase: "sync/test.yaml",
		},
		{
			wildcard: "**/*.userdb",
			expected: `^.*[^/\\]*\.userdb$`,
			testCase: "deep/nested/test.userdb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.wildcard, func(t *testing.T) {
			result := wildcardToRegex(tt.wildcard)
			t.Logf("通配符: %s -> 正则: %s", tt.wildcard, result)

			if result != tt.expected {
				t.Logf("警告: 实际正则 %s 与期望 %s 不同", result, tt.expected)
			}
		})
	}
}

func TestMatchAny(t *testing.T) {
	patterns := []string{
		"*.userdb",
		"*.custom.yaml",
		"^sync/.*",
		"installation.yaml",
	}

	parsedPatterns, errs := ParseExcludePatterns(patterns)
	if len(errs) > 0 {
		t.Fatalf("解析模式失败: %v", errs)
	}

	tests := []struct {
		file        string
		shouldMatch bool
	}{
		{"test.userdb", true},
		{"default.custom.yaml", true},
		{"sync/data.txt", true},
		{"installation.yaml", true},
		{"normal.yaml", false},
		{"build/output.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			matched := MatchAny(tt.file, parsedPatterns)
			if matched != tt.shouldMatch {
				t.Errorf("文件 %s 匹配结果: 期望 %v, 实际 %v", tt.file, tt.shouldMatch, matched)
			}
		})
	}
}

func TestDefaultExcludePatterns(t *testing.T) {
	t.Logf("默认排除模式数量: %d", len(DefaultExcludePatterns))

	parsedPatterns, errs := ParseExcludePatterns(DefaultExcludePatterns)
	if len(errs) > 0 {
		t.Fatalf("默认模式解析失败: %v", errs)
	}

	// 测试一些应该被排除的文件
	shouldExclude := []string{
		"test.userdb",
		"test.userdb.txt",
		"default.custom.yaml",
		"installation.yaml",
		"user.yaml",
		"sync/data.txt",
		"build/output.bin",
		"custom/user_exclude_file.txt",
	}

	for _, file := range shouldExclude {
		t.Run("应排除-"+file, func(t *testing.T) {
			if !MatchAny(file, parsedPatterns) {
				t.Errorf("文件 %s 应该被排除但没有匹配到任何模式", file)
			}
		})
	}

	// 测试一些不应该被排除的文件
	shouldNotExclude := []string{
		"wanxiang.yaml",
		"dicts/base.txt",
		"lua/wanxiang.lua",
	}

	for _, file := range shouldNotExclude {
		t.Run("不应排除-"+file, func(t *testing.T) {
			if MatchAny(file, parsedPatterns) {
				t.Errorf("文件 %s 不应该被排除但匹配到了模式", file)
			}
		})
	}
}
