package memes_cli

import (
	"encoding/json"
	"fmt"
)

// UploadImage 上传图片
func (c *APIClient) UploadImage(imageType, urlStr, headers, path, data string) (string, error) {
	requestBody := map[string]interface{}{
		"type": imageType,
	}

	switch imageType {
	case "url":
		requestBody["url"] = urlStr
		if headers != "" {
			requestBody["headers"] = headers
		}
	case "path":
		requestBody["path"] = path
	case "data":
		requestBody["data"] = data
	default:
		return "", fmt.Errorf("invalid image type: %s", imageType)
	}

	respBody, err := c.Request("POST", "/image/upload", requestBody, nil)
	if err != nil {
		return "", err
	}

	var result struct {
		ImageID string `json:"image_id"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("unmarshal response failed: %v", err)
	}

	return result.ImageID, nil
}

// GetImage 获取图片
func (c *APIClient) GetImage(imageID string) ([]byte, error) {
	return c.Request("GET", "/image/"+imageID, nil, nil)
}

// GetVersion 获取版本号
func (c *APIClient) GetVersion() (string, error) {
	respBody, err := c.Request("GET", "/meme/version", nil, nil)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}

// GetMemeKeys 获取表情名列表
func (c *APIClient) GetMemeKeys() ([]string, error) {
	respBody, err := c.Request("GET", "/meme/keys", nil, nil)
	if err != nil {
		return nil, err
	}

	var keys []string
	if err := json.Unmarshal(respBody, &keys); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}

	return keys, nil
}

// GetMemeInfos 获取表情信息列表
func (c *APIClient) GetMemeInfos() ([]MemeInfo, error) {
	respBody, err := c.Request("GET", "/meme/infos", nil, nil)
	if err != nil {
		return nil, err
	}

	var infos []MemeInfo
	if err := json.Unmarshal(respBody, &infos); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}

	return infos, nil
}

// SearchMeme 搜索表情
func (c *APIClient) SearchMeme(query string, includeTags bool) ([]string, error) {
	params := map[string]string{
		"query": query,
	}
	if includeTags {
		params["include_tags"] = "true"
	}

	respBody, err := c.Request("GET", "/meme/search", nil, params)
	if err != nil {
		return nil, err
	}

	var results []string
	if err := json.Unmarshal(respBody, &results); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}

	return results, nil
}

// GetMemeInfo 获取表情信息
func (c *APIClient) GetMemeInfo(key string) (*MemeInfo, error) {
	respBody, err := c.Request("GET", "/memes/"+key+"/info", nil, nil)
	if err != nil {
		return nil, err
	}

	var info MemeInfo
	if err := json.Unmarshal(respBody, &info); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}

	return &info, nil
}

// GetMemePreview 获取表情预览
func (c *APIClient) GetMemePreview(key string) (string, error) {
	respBody, err := c.Request("GET", "/memes/"+key+"/preview", nil, nil)
	if err != nil {
		return "", err
	}

	var result struct {
		ImageID string `json:"image_id"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("unmarshal response failed: %v", err)
	}

	return result.ImageID, nil
}

// CreateMeme 制作表情
func (c *APIClient) CreateMeme(key string, images []Image, texts []string, options map[string]interface{}) (string, error) {
	type CreateMemeBody struct {
		Images  []Image                `json:"images"`
		Texts   []string               `json:"texts"`
		Options map[string]interface{} `json:"options"`
	}
	if len(images) <= 0 {
		images = []Image{}
	}
	if len(texts) <= 0 {
		texts = []string{}
	}
	if options == nil {
		options = map[string]interface{}{}
	}
	requestBody := CreateMemeBody{
		Images:  images,
		Texts:   texts,
		Options: options,
	}

	respBody, err := c.Request("POST", "/memes/"+key, requestBody, nil)
	if err != nil {
		return "", err
	}

	var result struct {
		ImageID string `json:"image_id"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("unmarshal response failed: %v", err)
	}

	return result.ImageID, nil
}
