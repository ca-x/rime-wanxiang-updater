package fileutil

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadFile 下载文件，支持断点续传
func DownloadFile(url, dest string, client *http.Client) error {
	// 检查是否支持断点续传
	var downloaded int64 = 0
	if info, err := os.Stat(dest); err == nil {
		downloaded = info.Size()
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置断点续传
	if downloaded > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", downloaded))
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查是否支持断点续传
	var out *os.File
	if resp.StatusCode == http.StatusPartialContent {
		// 支持断点续传，追加模式打开
		out, err = os.OpenFile(dest, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		// 不支持断点续传，重新下载
		downloaded = 0
		out, err = os.Create(dest)
	}

	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	// 下载文件
	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, err := out.Write(buf[:n]); err != nil {
				return fmt.Errorf("写入文件失败: %w", err)
			}
			downloaded += int64(n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取数据失败: %w", err)
		}
	}

	return nil
}
