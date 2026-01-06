package fileutil

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

// ExtractZip 解压 ZIP 文件，支持排除模式
func ExtractZip(src, dest string, excludeFiles []string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("打开压缩包失败: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if shouldExclude(f.Name, excludeFiles) {
			continue
		}

		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return fmt.Errorf("创建目录失败: %w", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return fmt.Errorf("创建父目录失败: %w", err)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("创建文件失败: %w", err)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("打开压缩文件失败: %w", err)
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return fmt.Errorf("解压文件失败: %w", err)
		}
	}
	return nil
}

// shouldExclude 检查文件是否应该被排除
func shouldExclude(filename string, excludePatterns []string) bool {
	for _, pattern := range excludePatterns {
		matched, _ := regexp.MatchString(pattern, filename)
		if matched {
			return true
		}
	}
	return false
}

// GetZipFileList 获取 ZIP 文件中的文件列表
func GetZipFileList(zipPath string) ([]string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("打开压缩包失败: %w", err)
	}
	defer r.Close()

	var files []string
	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
			files = append(files, f.Name)
		}
	}
	return files, nil
}

// HandleCNBNestedDir 处理 CNB 镜像解压后的嵌套目录问题
// CNB 镜像解压后可能会有额外的一层嵌套目录，例如：
// temp_dir/base-dicts/base-dicts/files... (需要处理)
// 而不是：temp_dir/base-dicts/files... (正常情况)
func HandleCNBNestedDir(extractPath, zipFileName string) error {
	// 去掉 .zip 后缀获取目录名
	dirName := zipFileName[:len(zipFileName)-4]

	// 检查是否存在嵌套目录
	nestedPath := filepath.Join(extractPath, dirName)
	if !FileExists(nestedPath) {
		// 没有嵌套，直接返回
		return nil
	}

	// 检查嵌套目录中的内容
	entries, err := os.ReadDir(nestedPath)
	if err != nil {
		return fmt.Errorf("读取目录失败: %w", err)
	}

	// 如果嵌套目录中只有一个同名目录，则将其内容移动到上层
	if len(entries) == 1 && entries[0].IsDir() && entries[0].Name() == dirName {
		innerPath := filepath.Join(nestedPath, dirName)

		// 临时目录
		tempPath := filepath.Join(extractPath, "_temp_"+dirName)

		// 先移动到临时目录
		if err := os.Rename(innerPath, tempPath); err != nil {
			return fmt.Errorf("移动到临时目录失败: %w", err)
		}

		// 删除原嵌套目录
		if err := os.RemoveAll(nestedPath); err != nil {
			return fmt.Errorf("删除嵌套目录失败: %w", err)
		}

		// 将临时目录重命名为正确的目录名
		if err := os.Rename(tempPath, nestedPath); err != nil {
			return fmt.Errorf("重命名目录失败: %w", err)
		}
	}

	return nil
}
