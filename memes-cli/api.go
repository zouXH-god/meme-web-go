package memes_cli

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func GetLocalImagePath(imageID string) string {
	imageDir := "./images"
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		fmt.Println("Failed to create image directory: %v", err)
		return ""
	}
	imagePath := filepath.Join("./images", imageID)
	return imagePath
}

func SaveLocalImage(data []byte) (string, error) {
	// 1. 生成唯一的 image_id（这里用 MD5 哈希，也可以用 UUID）
	hash := md5.Sum(data)
	imageID := "memeWeb" + hex.EncodeToString(hash[:])

	// 2. 保存图片到本地（示例：存储到 ./images/{imageID}）
	imagePath := GetLocalImagePath(imageID)
	if err := os.WriteFile(imagePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}

	// 3. 返回生成的 image_id
	return imageID, nil
}

func GetLocalImage(imageID string) ([]byte, error) {
	imagePath := GetLocalImagePath(imageID)
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %v", err)
	}
	return data, nil
}

// UploadImage 上传图片
// imageType: 图片类型 ("url", "path", "data")
// urlStr: 图片URL（当 type="url" 时使用）
// headers: HTTP 请求头（JSON 字符串，当 type="url" 时使用）
// path: 本地文件路径（当 type="path" 时使用）
// data: Base64 编码的图片数据（当 type="data" 时使用）
func (c *APIClient) UploadImage(imageType, urlStr, headers, path, data string) (string, error) {
	// 1. 校验参数
	if imageType != "url" && imageType != "path" && imageType != "data" {
		return "", fmt.Errorf("invalid image type: %s", imageType)
	}

	// 2. Python 模式：本地保存图片
	if c.memeType == "python" {
		switch imageType {
		case "url":
			// 下载 URL 图片并保存
			resp, err := http.Get(urlStr)
			if err != nil {
				return "", fmt.Errorf("failed to download image from URL: %v", err)
			}
			defer resp.Body.Close()

			imageData, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", fmt.Errorf("failed to read image data: %v", err)
			}

			return SaveLocalImage(imageData)
		case "path":
			// 读取本地文件
			imageData, err := os.ReadFile(path)
			if err != nil {
				return "", fmt.Errorf("failed to read local image: %v", err)
			}
			return SaveLocalImage(imageData)
		case "data":
			// 解析 Base64 数据
			imageData, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				return "", fmt.Errorf("failed to decode base64 image: %v", err)
			}
			return SaveLocalImage(imageData)
		}
	}

	// 3. 默认模式：API 上传
	requestBody := map[string]interface{}{
		"type": imageType,
	}

	switch imageType {
	case "url":
		requestBody["url"] = urlStr
		if headers != "" {
			var headersMap map[string]string
			if err := json.Unmarshal([]byte(headers), &headersMap); err != nil {
				return "", fmt.Errorf("invalid headers JSON: %v", err)
			}
			requestBody["headers"] = headersMap
		}
	case "path":
		requestBody["path"] = path
	case "data":
		requestBody["data"] = data
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
	if c.memeType == "python" {
		return GetLocalImage(imageID)
	}
	return c.Request("GET", "/image/"+imageID, nil, nil)
}

// GetMemeType 获取版本号
func (c *APIClient) GetMemeType() (string, error) {
	respBody, err := c.Request("GET", "/meme/code", nil, nil)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
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

// GetMemePreview 获取 meme 预览，返回 image_id 或错误
func (c *APIClient) GetMemePreview(key string) (string, error) {
	respBody, err := c.Request("GET", "/memes/"+key+"/preview", nil, nil)
	if err != nil {
		return "", err
	}
	fmt.Println(c.memeType == "python")
	if c.memeType == "python" {
		return SaveLocalImage(respBody)
	}
	// 尝试解析 JSON
	var result struct {
		ImageID string `json:"image_id"`
	}
	if err := json.Unmarshal(respBody, &result); err == nil && result.ImageID != "" {
		return result.ImageID, nil
	}
	return "", err
}

// CreateMeme 制作表情包
// key: 表情模板ID
// images: 图片列表（可以是本地文件或已有 Image ID）
// texts: 文本列表
// options: 额外选项
func (c *APIClient) CreateMeme(key string, images []Image, texts []string, options map[string]interface{}) (string, error) {
	// 1. 处理默认值
	if len(images) == 0 {
		images = []Image{}
	}
	if len(texts) == 0 {
		texts = []string{}
	}
	if options == nil {
		options = map[string]interface{}{}
	}

	// 2. 根据模式选择不同的请求方式
	if c.memeType == "python" {
		// Python 模式：multipart/form-data 上传本地文件
		return c.createMemePythonMode(key, images, texts, options)
	} else {
		// 默认模式：JSON 上传（直接传 Image ID）
		return c.createMemeDefaultMode(key, images, texts, options)
	}
}

// createMemeDefaultMode 默认模式（JSON 上传）
func (c *APIClient) createMemeDefaultMode(key string, images []Image, texts []string, options map[string]interface{}) (string, error) {
	requestBody := map[string]interface{}{
		"images":  images,
		"texts":   texts,
		"options": options,
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

// createMemePythonMode Python 模式（multipart/form-data 上传文件）
func (c *APIClient) createMemePythonMode(key string, images []Image, texts []string, options map[string]interface{}) (string, error) {
	// 1. 创建 multipart 请求体
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 2. 添加 texts 字段
	for _, text := range texts {
		if err := writer.WriteField("texts", text); err != nil {
			return "", fmt.Errorf("failed to write text field: %v", err)
		}
	}

	// 3. 添加 options 字段（JSON 编码）
	if len(options) > 0 {
		optionsJSON, err := json.Marshal(options)
		if err != nil {
			return "", fmt.Errorf("failed to marshal options: %v", err)
		}
		if err := writer.WriteField("options", string(optionsJSON)); err != nil {
			return "", fmt.Errorf("failed to write options field: %v", err)
		}
	}

	// 4. 添加 images 文件
	for _, image := range images {
		imagePath := GetLocalImagePath(image.ID)
		file, err := os.Open(imagePath)
		if err != nil {
			return "", fmt.Errorf("failed to open image file: %v", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile("images", filepath.Base(imagePath))
		if err != nil {
			return "", fmt.Errorf("failed to create form file: %v", err)
		}

		if _, err := io.Copy(part, file); err != nil {
			return "", fmt.Errorf("failed to copy file data: %v", err)
		}
	}

	// 5. 关闭 writer
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// 6. 创建请求
	req, err := http.NewRequest("POST", c.BaseURL+"/memes/"+key, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// 7. 设置 Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 8. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 9. 解析响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return SaveLocalImage(respBody)
}
