package view

import (
	"fmt"
	"github.com/gin-gonic/gin"
	memesCli "meme-web-go/memes-cli"
	"net/http"
	"strings"
)

// MemesApiRout 封装 API 路由
type MemesApiRout struct {
	client    *memesCli.APIClient
	MemeInfos []memesCli.MemeInfo
}

// MemesApiRouts 封装 API 路由
type MemesApiRouts struct {
	memes    []MemesApiRout
	keys     []string
	memeKeys map[int][]string
	keyIndex map[string]int
	keyInfo  map[string]*memesCli.MemeInfo
	Infos    []*memesCli.MemeInfo
}

// NewMemesApiRouts 创建新的 API 路由
func NewMemesApiRouts(clients []*memesCli.APIClient) *MemesApiRouts {
	apis := &MemesApiRouts{}
	for _, client := range clients {
		api := &MemesApiRout{client: client}
		apis.memes = append(apis.memes, *api)
	}
	apis.InitInfo()
	return apis
}

// ContainsKey 使用 for 循环判断 key 是否存在 slice 中
func ContainsKey(slice []string, key string) bool {
	for _, s := range slice {
		if s == key {
			return true
		}
	}
	return false
}

func (r *MemesApiRouts) InitInfo() {
	// 变量初始化
	var err error
	r.keys = []string{}
	r.memeKeys = map[int][]string{}
	r.keyIndex = map[string]int{}
	r.keyInfo = map[string]*memesCli.MemeInfo{}
	r.Infos = []*memesCli.MemeInfo{}
	// 获取所有信息
	for index, client := range r.memes {
		client.MemeInfos, err = client.client.GetMemeInfos()
		if err != nil {
			panic(err)
		}
		// 结构化 keys
		for _, info := range client.MemeInfos {
			if ContainsKey(r.keys, info.Key) {
				continue
			}
			r.keys = append(r.keys, info.Key)
			r.memeKeys[index] = append(r.memeKeys[index], info.Key)
			r.keyIndex[info.Key] = index
			r.keyInfo[info.Key] = &info
			r.Infos = append(r.Infos, &info)
		}
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

	var imageIDs []string
	for _, client := range r.memes {
		imageID, err := client.client.UploadImage(request.Type, request.URL, request.Headers, request.Path, request.Data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		imageIDs = append(imageIDs, imageID)
	}

	c.JSON(http.StatusOK, gin.H{"image_id": "ADD_" + strings.Join(imageIDs, ",")})
}

// HandleGetImage 获取图片
func (r *MemesApiRouts) HandleGetImage(c *gin.Context) {
	imageID := c.Param("image_id")
	if strings.HasPrefix(imageID, "ADD_") {
		imageID = strings.Split(imageID, ",")[1]
	}
	var imageData []byte
	var err error
	for _, client := range r.memes {
		imageData, err = client.client.GetImage(imageID)
		if imageData != nil {
			break
		}
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/*", imageData)
}

// HandleGetVersion 获取版本号
func (r *MemesApiRouts) HandleGetVersion(c *gin.Context) {
	versions := make(map[string]string)
	for _, client := range r.memes {
		version, err := client.client.GetVersion()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		versions[client.client.BaseURL] = version
	}

	c.JSON(http.StatusOK, versions)
}

// HandleGetMemeKeys 获取表情名列表
func (r *MemesApiRouts) HandleGetMemeKeys(c *gin.Context) {
	c.JSON(http.StatusOK, r.keys)
}

// HandleGetMemeKeywords 获取表情关键词名列表
func (r *MemesApiRouts) HandleGetMemeKeywords(c *gin.Context) {
	var keywords [][]string
	for _, key := range r.keys {
		keywords = append(keywords, r.keyInfo[key].Keywords)
	}
	c.JSON(http.StatusOK, keywords)
}

// HandleGetMemeInfos 获取表情信息列表
func (r *MemesApiRouts) HandleGetMemeInfos(c *gin.Context) {
	c.JSON(http.StatusOK, r.Infos)
}

// HandleSearchMeme 搜索表情
func (r *MemesApiRouts) HandleSearchMeme(c *gin.Context) {
	query := c.Query("query")
	includeTags := c.Query("include_tags") == "true"

	var results []*memesCli.MemeInfo
	for _, key := range r.keys {
		if strings.Contains(key, query) {
			results = append(results, r.keyInfo[key])
		}
	}
	for _, info := range r.Infos {
		if includeTags && strings.Contains(strings.Join(info.Keywords, " "), query) {
			results = append(results, info)
		}
	}
	c.JSON(http.StatusOK, results)
}

// HandleGetMemeInfo 获取表情信息
func (r *MemesApiRouts) HandleGetMemeInfo(c *gin.Context) {
	key := c.Param("key")
	c.JSON(http.StatusOK, r.keyInfo[key])
}

// HandleGetMemePreview 获取表情预览
func (r *MemesApiRouts) HandleGetMemePreview(c *gin.Context) {
	key := c.Param("key")
	cli := r.memes[r.keyIndex[key]]
	imageID, err := cli.client.GetMemePreview(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"image_id": imageID})
}

// HandleCreateMeme 制作表情
func (r *MemesApiRouts) HandleCreateMeme(c *gin.Context) {
	key := c.Param("key")
	memeIndex := r.keyIndex[key]
	cli := r.memes[memeIndex]

	var request struct {
		Images  []memesCli.Image       `json:"images"`
		Texts   []string               `json:"texts"`
		Options map[string]interface{} `json:"options"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for index, image := range request.Images {
		if strings.HasPrefix(image.ID, "ADD_") {
			imageID := strings.Replace(image.ID, "ADD_", "", 1)
			imageID = strings.Split(imageID, ",")[memeIndex]
			request.Images[index].ID = imageID
		}
	}

	imageID, err := cli.client.CreateMeme(key, request.Images, request.Texts, request.Options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"image_id": imageID})
}

// HandleMkRenderList 创建渲染列表
func (r *MemesApiRouts) HandleMkRenderList(c *gin.Context) {
	var images [][]byte
	// 获取图片列表
	for index, cli := range r.memes {
		imageID, err := cli.client.MkRenderList(r.memeKeys[index])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		image, err := cli.client.GetImage(imageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		images = append(images, image)
	}
	// 拼接图片
	imageID, err := memesCli.CombineImagesVertically(images)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"image_id": imageID})
}
