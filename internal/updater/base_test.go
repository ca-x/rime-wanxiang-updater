package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/types"
)

func TestHasUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	recordPath := filepath.Join(tmpDir, "test_record.json")

	// Create config manager with temp config
	configPath := filepath.Join(tmpDir, "config.json")
	cfg := &config.Manager{
		ConfigPath: configPath,
		Config: &types.Config{
			SchemeType: "base",
			UseMirror:  true,
		},
	}

	updater := NewBaseUpdater(cfg)

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	tests := []struct {
		name         string
		updateInfo   *types.UpdateInfo
		existingTime *time.Time
		want         bool
	}{
		{
			name:         "No update info",
			updateInfo:   nil,
			existingTime: nil,
			want:         false,
		},
		{
			name: "No local record",
			updateInfo: &types.UpdateInfo{
				UpdateTime: now,
			},
			existingTime: nil,
			want:         true,
		},
		{
			name: "Remote is newer",
			updateInfo: &types.UpdateInfo{
				UpdateTime: tomorrow,
			},
			existingTime: &yesterday,
			want:         true,
		},
		{
			name: "Remote is older",
			updateInfo: &types.UpdateInfo{
				UpdateTime: yesterday,
			},
			existingTime: &tomorrow,
			want:         false,
		},
		{
			name: "Same time",
			updateInfo: &types.UpdateInfo{
				UpdateTime: now,
			},
			existingTime: &now,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up record file
			os.Remove(recordPath)

			// Create existing record if needed
			if tt.existingTime != nil {
				record := types.UpdateRecord{
					UpdateTime: *tt.existingTime,
				}
				data, _ := json.MarshalIndent(record, "", "  ")
				os.WriteFile(recordPath, data, 0644)
			}

			got := updater.HasUpdate(tt.updateInfo, recordPath)
			if got != tt.want {
				t.Errorf("HasUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLocalRecord(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	cfg := &config.Manager{
		ConfigPath: configPath,
		Config: &types.Config{
			SchemeType: "base",
			UseMirror:  true,
		},
	}

	updater := NewBaseUpdater(cfg)

	tests := []struct {
		name        string
		setupRecord func(string) error
		wantNil     bool
		checkRecord func(*types.UpdateRecord) error
	}{
		{
			name: "Valid record",
			setupRecord: func(path string) error {
				record := types.UpdateRecord{
					Name:       "test-scheme",
					UpdateTime: time.Now(),
					Tag:        "v1.0",
					SHA256:     "abc123",
				}
				data, err := json.MarshalIndent(record, "", "  ")
				if err != nil {
					return err
				}
				return os.WriteFile(path, data, 0644)
			},
			wantNil: false,
			checkRecord: func(r *types.UpdateRecord) error {
				if r.Name != "test-scheme" {
					return fmt.Errorf("Name mismatch: got %s, want test-scheme", r.Name)
				}
				if r.Tag != "v1.0" {
					return fmt.Errorf("Tag mismatch: got %s, want v1.0", r.Tag)
				}
				return nil
			},
		},
		{
			name: "Non-existent file",
			setupRecord: func(path string) error {
				return nil // Don't create file
			},
			wantNil: true,
		},
		{
			name: "Invalid JSON",
			setupRecord: func(path string) error {
				return os.WriteFile(path, []byte("invalid json"), 0644)
			},
			wantNil: true,
		},
		{
			name: "Empty file",
			setupRecord: func(path string) error {
				return os.WriteFile(path, []byte(""), 0644)
			},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recordPath := filepath.Join(tmpDir, fmt.Sprintf("record_%s.json", tt.name))

			// Setup
			if err := tt.setupRecord(recordPath); err != nil {
				t.Fatalf("Failed to setup record: %v", err)
			}

			// Test
			got := updater.GetLocalRecord(recordPath)

			// Check nil
			if (got == nil) != tt.wantNil {
				t.Errorf("GetLocalRecord() nil = %v, wantNil %v", got == nil, tt.wantNil)
				return
			}

			// Check record content if applicable
			if !tt.wantNil && tt.checkRecord != nil {
				if err := tt.checkRecord(got); err != nil {
					t.Errorf("Record validation failed: %v", err)
				}
			}
		})
	}
}

func TestSaveRecord(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	cfg := &config.Manager{
		ConfigPath: configPath,
		Config: &types.Config{
			SchemeType: "base",
			UseMirror:  true,
		},
	}

	updater := NewBaseUpdater(cfg)

	now := time.Now()
	updateInfo := &types.UpdateInfo{
		UpdateTime: now,
		Tag:        "v2.0",
		SHA256:     "def456",
		ID:         "123",
	}

	recordPath := filepath.Join(tmpDir, "saved_record.json")

	// Save record
	err := updater.SaveRecord(recordPath, "scheme", "test-scheme", updateInfo)
	if err != nil {
		t.Fatalf("SaveRecord() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(recordPath); os.IsNotExist(err) {
		t.Fatal("Record file was not created")
	}

	// Read and verify content
	data, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatalf("Failed to read saved record: %v", err)
	}

	var record types.UpdateRecord
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatalf("Failed to unmarshal saved record: %v", err)
	}

	// Verify fields
	if record.Name != "test-scheme" {
		t.Errorf("Name = %s, want test-scheme", record.Name)
	}
	if record.Tag != "v2.0" {
		t.Errorf("Tag = %s, want v2.0", record.Tag)
	}
	if record.SHA256 != "def456" {
		t.Errorf("SHA256 = %s, want def456", record.SHA256)
	}
	if record.CnbID != "123" {
		t.Errorf("CnbID = %s, want 123", record.CnbID)
	}
	// Allow small time differences due to JSON marshaling
	timeDiff := record.UpdateTime.Sub(now)
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}
	if timeDiff > time.Second {
		t.Errorf("UpdateTime difference too large: %v", timeDiff)
	}
}

func TestDownloadFileWithValidation(t *testing.T) {
	testContent := []byte("test file content with specific size")
	expectedSize := int64(len(testContent))

	tests := []struct {
		name         string
		content      []byte
		expectedSize int64
		wantErr      bool
		errContains  string
	}{
		{
			name:         "Valid size",
			content:      testContent,
			expectedSize: expectedSize,
			wantErr:      false,
		},
		{
			name:         "Size mismatch - smaller",
			content:      []byte("short"),
			expectedSize: expectedSize,
			wantErr:      true,
			errContains:  "文件大小不匹配",
		},
		{
			name:         "Size mismatch - larger",
			content:      []byte("this is a much longer content than expected for validation test"),
			expectedSize: expectedSize,
			wantErr:      true,
			errContains:  "文件大小不匹配",
		},
		{
			name:         "No validation (size=0)",
			content:      []byte("any content"),
			expectedSize: 0,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write(tt.content)
			}))
			defer server.Close()

			// Create temp directory and config
			tmpDir := t.TempDir()
			destFile := filepath.Join(tmpDir, "test_validation")
			configPath := filepath.Join(tmpDir, "config.json")

			cfg := &config.Manager{
				ConfigPath: configPath,
				Config: &types.Config{
					SchemeType: "base",
					UseMirror:  true,
				},
			}

			updater := NewBaseUpdater(cfg)

			// Test download with validation
			err := updater.DownloadFileWithValidation(
				server.URL,
				destFile,
				"test.zip",
				"test-source",
				tt.expectedSize,
				nil, // No progress callback
			)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadFileWithValidation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.errContains != "" {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("Error message = %v, should contain %s", err, tt.errContains)
				}
			}

			// If no error expected, verify file exists and has correct size
			if !tt.wantErr {
				fileInfo, err := os.Stat(destFile)
				if err != nil {
					t.Errorf("Downloaded file does not exist: %v", err)
					return
				}

				if tt.expectedSize > 0 && fileInfo.Size() != tt.expectedSize {
					t.Errorf("File size = %d, want %d", fileInfo.Size(), tt.expectedSize)
				}
			}

			// If error expected with size mismatch, file should be deleted
			if tt.wantErr && tt.errContains == "文件大小不匹配" {
				if _, err := os.Stat(destFile); !os.IsNotExist(err) {
					t.Error("Corrupted file should have been deleted")
				}
			}
		})
	}
}

func TestDownloadFileWithProgress(t *testing.T) {
	testContent := make([]byte, 1024*100) // 100KB
	for i := range testContent {
		testContent[i] = byte(i % 256)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(testContent)))
		w.WriteHeader(http.StatusOK)
		w.Write(testContent)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	destFile := filepath.Join(tmpDir, "test_progress")
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &config.Manager{
		ConfigPath: configPath,
		Config: &types.Config{
			SchemeType: "base",
			UseMirror:  true,
		},
	}

	updater := NewBaseUpdater(cfg)

	// Track progress calls
	progressCalled := false
	var lastPercent float64

	progressFunc := func(msg string, percent float64, source, fileName string, downloaded, totalSize int64, speed float64, isDownloading bool) {
		progressCalled = true
		lastPercent = percent
		t.Logf("Progress: %s (%.2f%%)", msg, percent*100)
	}

	err := updater.DownloadFile(server.URL, destFile, "test.dat", "test-source", progressFunc)
	if err != nil {
		t.Fatalf("DownloadFile() error = %v", err)
	}

	if !progressCalled {
		t.Error("Progress callback was not called")
	}

	if lastPercent != 1.0 {
		t.Errorf("Final progress = %.2f, want 1.0", lastPercent)
	}

	// Verify file
	fileInfo, err := os.Stat(destFile)
	if err != nil {
		t.Fatalf("Downloaded file does not exist: %v", err)
	}

	if fileInfo.Size() != int64(len(testContent)) {
		t.Errorf("File size = %d, want %d", fileInfo.Size(), len(testContent))
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || anyIndex(s, substr) >= 0)
}

func anyIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
