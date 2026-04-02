package api

import (
	"net/http"
	"testing"

	"rime-wanxiang-updater/internal/types"
)

func TestNewDownloadHTTPClientUsesNoTimeout(t *testing.T) {
	client := NewDownloadHTTPClient(&types.Config{})
	if client.Timeout != 0 {
		t.Fatalf("Timeout = %v, want 0", client.Timeout)
	}
}

func TestNewDownloadHTTPClientUsesHTTPProxyConfig(t *testing.T) {
	client := NewDownloadHTTPClient(&types.Config{
		ProxyEnabled: true,
		ProxyType:    "http",
		ProxyAddress: "127.0.0.1:7890",
	})

	if client.Transport == nil {
		t.Fatal("Transport = nil, want configured transport")
	}

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("Transport type = %T, want *http.Transport", client.Transport)
	}

	if transport.Proxy == nil {
		t.Fatal("Proxy = nil, want configured proxy function")
	}
}
