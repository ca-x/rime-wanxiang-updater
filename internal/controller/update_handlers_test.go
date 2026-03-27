package controller

import "testing"

func TestSuccessMessageForSingleUpdate(t *testing.T) {
	tests := []struct {
		name        string
		updateType  string
		localStatus string
		want        string
	}{
		{
			name:        "scheme first install",
			updateType:  "方案",
			localStatus: "未安装",
			want:        "方案安装完成！",
		},
		{
			name:        "scheme unknown version treated as install",
			updateType:  "方案",
			localStatus: "未知版本",
			want:        "方案安装完成！",
		},
		{
			name:        "scheme normal update",
			updateType:  "方案",
			localStatus: "v1.0.0",
			want:        "方案更新完成！",
		},
		{
			name:        "dict first install",
			updateType:  "词库",
			localStatus: "未安装",
			want:        "词库安装完成！",
		},
		{
			name:        "model normal update",
			updateType:  "模型",
			localStatus: "2025-01-01",
			want:        "模型更新完成！",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := successMessageForSingleUpdate(tt.updateType, tt.localStatus)
			if got != tt.want {
				t.Fatalf("successMessageForSingleUpdate(%q, %q) = %q, want %q", tt.updateType, tt.localStatus, got, tt.want)
			}
		})
	}
}

func TestSuccessMessageForAutoUpdate(t *testing.T) {
	tests := []struct {
		name        string
		components  []string
		versions    map[string]string
		wantMessage string
	}{
		{
			name:        "single install component",
			components:  []string{"方案"},
			versions:    map[string]string{"方案": "未安装"},
			wantMessage: "安装完成！",
		},
		{
			name:        "single updated component",
			components:  []string{"方案"},
			versions:    map[string]string{"方案": "v1.0.0"},
			wantMessage: "更新完成！",
		},
		{
			name:        "mixed install and update",
			components:  []string{"方案", "词库"},
			versions:    map[string]string{"方案": "未安装", "词库": "dict-nightly"},
			wantMessage: "安装和更新完成！",
		},
		{
			name:        "multiple installs only",
			components:  []string{"方案", "词库"},
			versions:    map[string]string{"方案": "未安装", "词库": "未知版本"},
			wantMessage: "安装完成！",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := successMessageForAutoUpdate(tt.components, tt.versions)
			if got != tt.wantMessage {
				t.Fatalf("successMessageForAutoUpdate(%v, %v) = %q, want %q", tt.components, tt.versions, got, tt.wantMessage)
			}
		})
	}
}
