package render

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ConvertResponse struct {
	Success bool   `json:"success"`
	Key     string `json:"key"`
	URL     string `json:"url"`
	Error   string `json:"error,omitempty"`
}

type BatchItem struct {
	Name string `json:"name,omitempty"`
	SVG  string `json:"svg"`
}

type BatchConvertResponse struct {
	Success bool              `json:"success"`
	Total   int               `json:"total"`
	Results []ConvertResponse `json:"results"`
	Message string            `json:"message,omitempty"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type Option func(*Client)

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

func NewClient(baseURL string, opts ...Option) *Client {
	c := &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Convert 单张 SVG 转换
func (c *Client) Convert(ctx context.Context, svgContent string) (*ConvertResponse, error) {
	if strings.TrimSpace(svgContent) == "" {
		return nil, fmt.Errorf("svgContent 不能为空")
	}

	reqURL := c.baseURL + "/convert"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBufferString(svgContent))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "image/svg+xml")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("接口返回异常状态码 [%d]: %s", resp.StatusCode, string(body))
	}

	var res ConvertResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w, 响应内容: %s", err, string(body))
	}

	if !res.Success {
		return nil, fmt.Errorf("转换业务失败: %s", res.Error)
	}

	return &res, nil
}

// BatchConvert 批量 SVG 转换
func (c *Client) BatchConvert(ctx context.Context, items []BatchItem) (*BatchConvertResponse, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("批量处理列表不能为空")
	}

	payload := map[string]interface{}{
		"svgs": items,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("序列化请求 JSON 失败: %w", err)
	}

	reqURL := c.baseURL + "/batch-convert"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("创建批量请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送批量请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("接口返回异常状态码 [%d]: %s", resp.StatusCode, string(body))
	}

	var res BatchConvertResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w, 响应内容: %s", err, string(body))
	}

	if !res.Success {
		return nil, fmt.Errorf("批量转换业务失败: %s", res.Message)
	}

	return &res, nil
}
