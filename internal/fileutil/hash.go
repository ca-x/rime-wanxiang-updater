package fileutil

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// CalculateSHA256 计算文件的 SHA256 哈希值
func CalculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("计算哈希失败: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// FileExists 检查文件是否存在
// 返回 true 表示文件存在（或存在但无法访问），此时应保留用户文件不覆盖
// 返回 false 表示文件确定不存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	// 只有当文件确定不存在时才返回 false
	return !os.IsNotExist(err)
}

// FormatBytes 格式化字节数为人类可读格式
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
