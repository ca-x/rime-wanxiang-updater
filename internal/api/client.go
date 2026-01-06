package api

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
	"rime-wanxiang-updater/internal/types"
)

// Client API 客户端
type Client struct {
	httpClient  *http.Client
	config      *types.Config
	githubToken string
}

// NewClient 创建新的 API 客户端
func NewClient(config *types.Config) *Client {
	return &Client{
		httpClient:  getHTTPClient(config),
		config:      config,
		githubToken: config.GithubToken,
	}
}

// GetHTTPClient 返回配置了代理的 HTTP 客户端
func getHTTPClient(config *types.Config) *http.Client {
	if !config.ProxyEnabled {
		return &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	var transport *http.Transport

	switch config.ProxyType {
	case "http", "https":
		proxyURL, err := url.Parse(fmt.Sprintf("http://%s", config.ProxyAddress))
		if err != nil {
			return &http.Client{Timeout: 30 * time.Second}
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

	case "socks5":
		dialer, err := proxy.SOCKS5("tcp", config.ProxyAddress, nil, proxy.Direct)
		if err != nil {
			return &http.Client{Timeout: 30 * time.Second}
		}
		transport = &http.Transport{
			Dial: dialer.Dial,
		}

	default:
		return &http.Client{Timeout: 30 * time.Second}
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

// Get 发送 GET 请求
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "RIME-Updater/1.0")
	if c.githubToken != "" && !c.config.UseMirror {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.githubToken))
	}

	// 如果使用镜像，设置特殊的 Accept 头
	if c.config.UseMirror {
		req.Header.Set("Accept", "application/vnd.cnb.web+json")
	}

	return c.httpClient.Do(req)
}

// Head 发送 HEAD 请求
func (c *Client) Head(url string) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "RIME-Updater/1.0")
	if c.githubToken != "" && !c.config.UseMirror {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.githubToken))
	}

	return c.httpClient.Do(req)
}
