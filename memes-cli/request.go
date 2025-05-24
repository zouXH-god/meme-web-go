package memes_cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// APIClient 封装 API 客户端
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewAPIClient 创建新的 API 客户端
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// Request 统一的请求方法
func (c *APIClient) Request(method, path string, body interface{}, queryParams map[string]string) ([]byte, error) {
	// 构建 URL
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, fmt.Errorf("parse URL failed: %v", err)
	}

	// 添加查询参数
	if queryParams != nil {
		q := u.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
	}

	// 准备请求体
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body failed: %v", err)
		}
		fmt.Println(string(jsonBody))
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// 创建请求
	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	// 设置请求头
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status code: %d, data: %s", resp.StatusCode, respBody)
	}

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %v", err)
	}

	return respBody, nil
}
