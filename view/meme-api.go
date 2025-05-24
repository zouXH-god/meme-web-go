package view

import (
	"fmt"
	"github.com/gin-gonic/gin"
	memesCli "meme-web-go/memes-cli"
	"net/http"
)

// MemesApiRouts 封装 API 路由
type MemesApiRouts struct {
	client    *memesCli.APIClient
	MemeInfos []memesCli.MemeInfo
}

// NewMemesApiRouts 创建新的 API 路由
func NewMemesApiRouts(client *memesCli.APIClient) *MemesApiRouts {
	api := &MemesApiRouts{client: client}
	api.InitInfo()
	return api
}

func (r *MemesApiRouts) InitInfo() {
	var err error
	r.MemeInfos, err = r.client.GetMemeInfos()
	if err != nil {
		panic(err)
	}
}

// HandleAPIRoot API 根路径
func (r *MemesApiRouts) HandleAPIRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Meme Generator API",
		"routes": []string{
			"POST   /api/image/upload",
			"GET    /api/image/:image_id",
			"GET    /api/meme/version",
			"GET    /api/meme/keys",
			"GET    /api/meme/infos",
			"GET    /api/meme/search",
			"GET    /api/memes/:key/info",
			"GET    /api/memes/:key/preview",
			"POST   /api/memes/:key",
		},
	})
}

// HandleUploadImage 上传图片
func (r *MemesApiRouts) HandleUploadImage(c *gin.Context) {
	var request struct {
		Type    string `json:"type" binding:"required"`
		URL     string `json:"url"`
		Headers string `json:"headers"`
		Path    string `json:"path"`
		Data    string `json:"data"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(request)

	imageID, err := r.client.UploadImage(request.Type, request.URL, request.Headers, request.Path, request.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"image_id": imageID})
}

// HandleGetImage 获取图片
func (r *MemesApiRouts) HandleGetImage(c *gin.Context) {
	imageID := c.Param("image_id")
	imageData, err := r.client.GetImage(imageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/*", imageData)
}

// HandleGetVersion 获取版本号
func (r *MemesApiRouts) HandleGetVersion(c *gin.Context) {
	version, err := r.client.GetVersion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, version)
}

// HandleGetMemeKeys 获取表情名列表
func (r *MemesApiRouts) HandleGetMemeKeys(c *gin.Context) {
	keys, err := r.client.GetMemeKeys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, keys)
}

// HandleGetMemeKeywords 获取表情关键词名列表
func (r *MemesApiRouts) HandleGetMemeKeywords(c *gin.Context) {
	var keywords [][]string
	for _, info := range r.MemeInfos {
		keywords = append(keywords, info.Keywords)
	}
	c.JSON(http.StatusOK, keywords)
}

// HandleGetMemeInfos 获取表情信息列表
func (r *MemesApiRouts) HandleGetMemeInfos(c *gin.Context) {
	infos, err := r.client.GetMemeInfos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, infos)
}

// HandleSearchMeme 搜索表情
func (r *MemesApiRouts) HandleSearchMeme(c *gin.Context) {
	query := c.Query("query")
	includeTags := c.Query("include_tags") == "true"

	results, err := r.client.SearchMeme(query, includeTags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// HandleGetMemeInfo 获取表情信息
func (r *MemesApiRouts) HandleGetMemeInfo(c *gin.Context) {
	key := c.Param("key")
	info, err := r.client.GetMemeInfo(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

// HandleGetMemePreview 获取表情预览
func (r *MemesApiRouts) HandleGetMemePreview(c *gin.Context) {
	key := c.Param("key")
	imageID, err := r.client.GetMemePreview(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"image_id": imageID})
}

// HandleCreateMeme 制作表情
func (r *MemesApiRouts) HandleCreateMeme(c *gin.Context) {
	key := c.Param("key")

	var request struct {
		Images  []memesCli.Image       `json:"images"`
		Texts   []string               `json:"texts"`
		Options map[string]interface{} `json:"options"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imageID, err := r.client.CreateMeme(key, request.Images, request.Texts, request.Options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"image_id": imageID})
}
