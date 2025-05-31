package main

import (
	"github.com/gin-gonic/gin"
	memesCli2 "meme-web-go/memes-cli"
	"meme-web-go/view"
	"os"
	"strings"
)

func registerRoutes(r *gin.Engine) {
	// 注册 memes 端点
	var memesClis []*memesCli2.APIClient
	points := strings.Split(os.Getenv("MEME_POINTS"), ",")
	for _, point := range points {
		memesClis = append(memesClis, memesCli2.NewAPIClient(point))
	}
	memesRouter := view.NewMemesApiRouts(memesClis)
	// 先注册 API 路由
	apiRoute := r.Group("/api")
	{
		apiRoute.GET("/", memesRouter.HandleAPIRoot)
		apiRoute.POST("/image/upload", memesRouter.HandleUploadImage)
		apiRoute.GET("/image/:image_id", memesRouter.HandleGetImage)
		apiRoute.GET("/meme/version", memesRouter.HandleGetVersion)
		apiRoute.GET("/meme/keys", memesRouter.HandleGetMemeKeys)
		apiRoute.GET("/meme/keywords", memesRouter.HandleGetMemeKeywords)
		apiRoute.GET("/meme/infos", memesRouter.HandleGetMemeInfos)
		apiRoute.GET("/meme/search", memesRouter.HandleSearchMeme)
		apiRoute.GET("/memes/:key/info", memesRouter.HandleGetMemeInfo)
		apiRoute.GET("/memes/:key/preview", memesRouter.HandleGetMemePreview)
		apiRoute.POST("/memes/:key", memesRouter.HandleCreateMeme)
	}

	// 然后手动设置静态文件路由
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// 其他静态文件
	r.Static("/static", "./static")

	// 对于 SPA 应用，可以添加以下处理
	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") &&
			!strings.HasPrefix(c.Request.URL.Path, "/static") {
			c.File("./static/index.html")
		}
	})
}
