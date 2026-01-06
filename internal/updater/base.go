package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"rime-wanxiang-updater/internal/api"
	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/deployer"
	"rime-wanxiang-updater/internal/fileutil"
	"rime-wanxiang-updater/internal/types"
)

// BaseUpdater 更新器基类
type BaseUpdater struct {
	Config    *config.Manager
	APIClient *api.Client
	Deployer  deployer.Deployer
}

// NewBaseUpdater 创建基础更新器
func NewBaseUpdater(cfg *config.Manager) *BaseUpdater {
	return &BaseUpdater{
		Config:    cfg,
		APIClient: api.NewClient(cfg.Config),
		Deployer:  deployer.GetDeployer(cfg.Config),
	}
}

// HasUpdate 检查是否有更新
func (b *BaseUpdater) HasUpdate(updateInfo *types.UpdateInfo, recordPath string) bool {
	if updateInfo == nil {
		return false
	}

	localTime := b.GetLocalTime(recordPath)
	if localTime == nil {
		return true
	}

	return updateInfo.UpdateTime.After(*localTime)
}

// GetLocalTime 获取本地记录的更新时间
func (b *BaseUpdater) GetLocalTime(recordPath string) *time.Time {
	record := b.GetLocalRecord(recordPath)
	if record == nil {
		return nil
	}
	return &record.UpdateTime
}

// GetLocalRecord 获取本地更新记录
func (b *BaseUpdater) GetLocalRecord(recordPath string) *types.UpdateRecord {
	if !fileutil.FileExists(recordPath) {
		return nil
	}

	data, err := os.ReadFile(recordPath)
	if err != nil {
		return nil
	}

	var record types.UpdateRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil
	}

	return &record
}

// SaveRecord 保存更新记录
func (b *BaseUpdater) SaveRecord(recordPath string, propertyType, propertyName string, info *types.UpdateInfo) error {
	record := types.UpdateRecord{
		Name:       propertyName,
		UpdateTime: info.UpdateTime,
		Tag:        info.Tag,
		ApplyTime:  time.Now(),
		SHA256:     info.SHA256,
		CnbID:      info.ID,
	}

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化记录失败: %w", err)
	}

	return os.WriteFile(recordPath, data, 0644)
}

// DownloadFile 下载文件
func (b *BaseUpdater) DownloadFile(url, dest, fileName, source string, progress types.ProgressFunc) error {
	// 使用 API 客户端的 HTTP 客户端进行下载
	req, err := b.APIClient.Get(url)
	if err != nil {
		return fmt.Errorf("创建下载请求失败: %w", err)
	}
	defer req.Body.Close()

	// 获取文件总大小
	totalSize := req.ContentLength

	// 检查是否支持断点续传
	var downloaded int64 = 0
	if info, err := os.Stat(dest); err == nil {
		downloaded = info.Size()
	}

	// 检查服务器响应
	var out *os.File
	if req.StatusCode == http.StatusPartialContent {
		out, err = os.OpenFile(dest, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		downloaded = 0
		out, err = os.Create(dest)
	}

	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	// 下载文件并报告进度
	buf := make([]byte, 32*1024)
	startTime := time.Now()
	lastUpdate := time.Now()

	for {
		n, err := req.Body.Read(buf)
		if n > 0 {
			if _, err := out.Write(buf[:n]); err != nil {
				return fmt.Errorf("写入文件失败: %w", err)
			}
			downloaded += int64(n)

			// 每 100ms 更新一次进度
			if progress != nil && time.Since(lastUpdate) > 100*time.Millisecond {
				elapsed := time.Since(startTime).Seconds()
				speed := float64(downloaded) / elapsed / 1024 / 1024 // MB/s

				var percent float64
				if totalSize > 0 {
					percent = float64(downloaded) / float64(totalSize)
				}

				// 格式化消息
				downloadedMB := float64(downloaded) / 1024 / 1024
				var msg string
				if totalSize > 0 {
					totalMB := float64(totalSize) / 1024 / 1024
					msg = fmt.Sprintf("下载中: %.2f MB / %.2f MB (%.2f MB/s)",
						downloadedMB, totalMB, speed)
				} else {
					msg = fmt.Sprintf("下载中: %.2f MB (%.2f MB/s)", downloadedMB, speed)
				}

				progress(msg, percent, source, fileName, downloaded, totalSize, speed, true)
				lastUpdate = time.Now()
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取数据失败: %w", err)
		}
	}

	// 下载完成
	if progress != nil {
		downloadedMB := float64(downloaded) / 1024 / 1024
		elapsed := time.Since(startTime).Seconds()
		speed := float64(downloaded) / elapsed / 1024 / 1024
		msg := fmt.Sprintf("下载完成: %.2f MB (平均 %.2f MB/s)", downloadedMB, speed)
		progress(msg, 1.0, "", "", 0, 0, 0, false)
	}

	return nil
}

// DownloadFileWithValidation 下载文件并验证大小
func (b *BaseUpdater) DownloadFileWithValidation(url, dest, fileName, source string, expectedSize int64, progress types.ProgressFunc) error {
	// 调用下载方法
	if err := b.DownloadFile(url, dest, fileName, source, progress); err != nil {
		return err
	}

	// 验证文件大小
	if expectedSize > 0 {
		fileInfo, err := os.Stat(dest)
		if err != nil {
			return fmt.Errorf("获取文件信息失败: %w", err)
		}

		actualSize := fileInfo.Size()
		if actualSize != expectedSize {
			// 删除损坏的文件
			os.Remove(dest)
			return fmt.Errorf("文件大小不匹配：期望 %d 字节，实际 %d 字节，已删除损坏的文件", expectedSize, actualSize)
		}
	}

	return nil
}

// ExtractZip 解压文件
func (b *BaseUpdater) ExtractZip(src, dest string) error {
	return fileutil.ExtractZip(src, dest, b.Config.Config.ExcludeFiles)
}

// CompareHash 比较文件哈希
func (b *BaseUpdater) CompareHash(remoteHash, filePath string) bool {
	if remoteHash == "" || !fileutil.FileExists(filePath) {
		return false
	}

	localHash, err := fileutil.CalculateSHA256(filePath)
	if err != nil {
		return false
	}

	return remoteHash == localHash
}

// CleanOldFiles 清理旧文件
func (b *BaseUpdater) CleanOldFiles(oldZip, newZip, extractPath string, isDict bool) error {
	if !fileutil.FileExists(oldZip) {
		return nil
	}

	oldFiles, err := fileutil.GetZipFileList(oldZip)
	if err != nil {
		return fmt.Errorf("获取旧文件列表失败: %w", err)
	}

	var newFiles []string
	if newZip != "" && fileutil.FileExists(newZip) {
		newFiles, _ = fileutil.GetZipFileList(newZip)
	}

	// 找出需要删除的文件
	toDelete := difference(oldFiles, newFiles)

	// 删除文件
	for _, file := range toDelete {
		fullPath := filepath.Join(extractPath, file)
		if fileutil.FileExists(fullPath) {
			os.Remove(fullPath)
		}
	}

	return nil
}

// difference 返回在 a 中但不在 b 中的元素
func difference(a, b []string) []string {
	mb := make(map[string]bool, len(b))
	for _, x := range b {
		mb[x] = true
	}
	var diff []string
	for _, x := range a {
		if !mb[x] {
			diff = append(diff, x)
		}
	}
	return diff
}

// TerminateProcesses 终止进程
func (b *BaseUpdater) TerminateProcesses() error {
	return b.Deployer.TerminateProcesses()
}

// Deploy 部署
func (b *BaseUpdater) Deploy() error {
	return b.Deployer.Deploy()
}
