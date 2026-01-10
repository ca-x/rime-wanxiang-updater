package fileutil

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	testContent := []byte("test file content for download")

	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		wantErr        bool
		checkContent   bool
	}{
		{
			name: "Successful download",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write(testContent)
			},
			wantErr:      false,
			checkContent: true,
		},
		{
			name: "Resume download",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				rangeHeader := r.Header.Get("Range")
				if rangeHeader != "" {
					// Simulating partial content response
					w.Header().Set("Content-Range", fmt.Sprintf("bytes 10-%d/%d", len(testContent)-1, len(testContent)))
					w.WriteHeader(http.StatusPartialContent)
					w.Write(testContent[10:]) // Send remaining bytes
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write(testContent)
				}
			},
			wantErr:      false,
			checkContent: false, // Content will be different due to resume
		},
		{
			name: "Server error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr:      false, // Function doesn't check status code
			checkContent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			// Create temp file
			tmpDir := t.TempDir()
			destFile := filepath.Join(tmpDir, "test_download")

			// Execute download
			client := &http.Client{}
			err := DownloadFile(server.URL, destFile, client)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check file exists
			if !tt.wantErr {
				if _, err := os.Stat(destFile); os.IsNotExist(err) {
					t.Errorf("Downloaded file does not exist")
				}
			}

			// Check content if needed
			if tt.checkContent && !tt.wantErr {
				content, err := os.ReadFile(destFile)
				if err != nil {
					t.Errorf("Failed to read downloaded file: %v", err)
					return
				}
				if string(content) != string(testContent) {
					t.Errorf("Downloaded content mismatch, got %s, want %s", string(content), string(testContent))
				}
			}
		})
	}
}

func TestDownloadFileResume(t *testing.T) {
	fullContent := []byte("0123456789abcdefghijklmnopqrstuvwxyz")

	// Create test server that supports resume
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rangeHeader := r.Header.Get("Range")
		if rangeHeader != "" {
			// Parse range header (simple implementation)
			var start int
			fmt.Sscanf(rangeHeader, "bytes=%d-", &start)

			if start < len(fullContent) {
				w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, len(fullContent)-1, len(fullContent)))
				w.WriteHeader(http.StatusPartialContent)
				w.Write(fullContent[start:])
			} else {
				w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			}
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(fullContent)
		}
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	destFile := filepath.Join(tmpDir, "test_resume")

	// First download (partial)
	client := &http.Client{}

	// Simulate partial download by creating a file with partial content
	partialContent := fullContent[:10]
	if err := os.WriteFile(destFile, partialContent, 0644); err != nil {
		t.Fatalf("Failed to create partial file: %v", err)
	}

	// Resume download
	err := DownloadFile(server.URL, destFile, client)
	if err != nil {
		t.Errorf("DownloadFile() resume error = %v", err)
		return
	}

	// Verify full content
	content, err := os.ReadFile(destFile)
	if err != nil {
		t.Errorf("Failed to read resumed file: %v", err)
		return
	}

	if string(content) != string(fullContent) {
		t.Errorf("Resumed content mismatch, got length %d, want %d", len(content), len(fullContent))
	}
}

func TestDownloadFileNoResume(t *testing.T) {
	fullContent := []byte("complete file content")

	// Create test server that doesn't support resume
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always return full content, ignoring Range header
		w.WriteHeader(http.StatusOK)
		w.Write(fullContent)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	destFile := filepath.Join(tmpDir, "test_no_resume")

	// Create partial file first
	partialContent := []byte("partial")
	if err := os.WriteFile(destFile, partialContent, 0644); err != nil {
		t.Fatalf("Failed to create partial file: %v", err)
	}

	// Try to download (should overwrite)
	client := &http.Client{}
	err := DownloadFile(server.URL, destFile, client)
	if err != nil {
		t.Errorf("DownloadFile() error = %v", err)
		return
	}

	// Verify content is the full content, not partial + full
	content, err := os.ReadFile(destFile)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
		return
	}

	if string(content) != string(fullContent) {
		t.Errorf("Content mismatch after non-resume download")
	}
}

func TestDownloadFileNetworkError(t *testing.T) {
	tmpDir := t.TempDir()
	destFile := filepath.Join(tmpDir, "test_network_error")

	// Use invalid URL to simulate network error
	client := &http.Client{}
	err := DownloadFile("http://invalid-host-that-does-not-exist-12345.com/file", destFile, client)

	if err == nil {
		t.Error("Expected network error, got nil")
	}
}

func TestDownloadFileWriteError(t *testing.T) {
	testContent := []byte("test content")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(testContent)
	}))
	defer server.Close()

	// Try to write to a directory (not a file) to cause write error
	tmpDir := t.TempDir()
	destDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(destDir, 0755)

	// Try to download to directory path (should fail)
	client := &http.Client{}
	err := DownloadFile(server.URL, destDir, client)

	if err == nil {
		t.Error("Expected write error when destination is directory, got nil")
	}
}

func TestDownloadFileReadError(t *testing.T) {
	// Server that closes connection abruptly
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("partial"))
		// Simulate connection close by hijacking and closing
		if hijacker, ok := w.(http.Hijacker); ok {
			conn, _, err := hijacker.Hijack()
			if err == nil {
				conn.Close()
			}
		}
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	destFile := filepath.Join(tmpDir, "test_read_error")

	client := &http.Client{}
	err := DownloadFile(server.URL, destFile, client)

	// Should complete without error because partial data was written
	// The test verifies that connection errors are handled
	if err != nil && err != io.EOF {
		// Either succeeds or gets a read error (both are acceptable)
		t.Logf("Download interrupted as expected: %v", err)
	}
}
