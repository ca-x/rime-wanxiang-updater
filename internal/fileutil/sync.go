package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SyncDirectory 同步目录内容到目标目录
// sourceDir: 源目录
// targetDir: 目标目录
// excludePatterns: 排除的文件模式（glob 或正则表达式）
func SyncDirectory(sourceDir, targetDir string, excludePatterns []string) error {
	// 确保目标目录存在
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 遍历源目录
	return filepath.Walk(sourceDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(sourceDir, srcPath)
		if err != nil {
			return err
		}

		// 跳过根目录自身
		if relPath == "." {
			return nil
		}

		// 目标路径
		dstPath := filepath.Join(targetDir, relPath)

		// 检查是否应该排除：仅在目标位置已存在时才跳过匹配排除模式的文件。
		// 这确保首次安装时默认配置文件能被正确部署，
		// 同时在后续更新时保留用户的自定义修改。
		if shouldExclude(relPath, excludePatterns) && fileExistsAtPath(dstPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 处理目录
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// 复制文件
		return copyFile(srcPath, dstPath)
	})
}

// shouldExclude 检查文件是否应该被排除
func shouldExclude(relPath string, excludePatterns []string) bool {
	// 始终排除 build 目录
	if strings.HasPrefix(relPath, "build"+string(filepath.Separator)) || relPath == "build" {
		return true
	}

	// 检查用户配置的排除规则
	for _, pattern := range excludePatterns {
		// 简单的 glob 匹配
		if matched, _ := filepath.Match(pattern, filepath.Base(relPath)); matched {
			return true
		}

		// 检查是否匹配完整路径
		if matched, _ := filepath.Match(pattern, relPath); matched {
			return true
		}

		// 通配符匹配
		if strings.Contains(pattern, "*") {
			// 处理 **/*.ext 这样的模式
			if strings.HasPrefix(pattern, "**/") {
				suffix := strings.TrimPrefix(pattern, "**/")
				if matched, _ := filepath.Match(suffix, filepath.Base(relPath)); matched {
					return true
				}
			}
		}
	}

	return false
}

// fileExistsAtPath 检查文件或目录是否存在于指定路径
// 返回 true 表示文件存在（或存在但无法访问），此时应保留用户文件不覆盖
// 返回 false 表示文件确定不存在，可以安全部署新文件
func fileExistsAtPath(path string) bool {
	_, err := os.Stat(path)
	// 只有当文件确定不存在时才返回 false
	// 其他情况（包括存在但无权限访问）都返回 true，以保护用户文件
	return !os.IsNotExist(err)
}
