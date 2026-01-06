package fileutil

import (
	"fmt"
	"io"
	"os"
)

// MoveFile 移动文件，支持跨磁盘分区
// 优先使用 os.Rename（同磁盘快速），失败时使用 copy+delete（跨磁盘兼容）
func MoveFile(src, dst string) error {
	// 尝试直接重命名（适用于同磁盘分区）
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}

	// 如果 rename 失败（可能是跨磁盘），使用 copy+delete
	if err := copyFile(src, dst); err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	// 删除源文件
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("删除临时文件失败: %w", err)
	}

	return nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 获取源文件信息
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 复制内容
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// 同步到磁盘
	if err := destFile.Sync(); err != nil {
		return err
	}

	// 设置文件权限
	return os.Chmod(dst, sourceInfo.Mode())
}
