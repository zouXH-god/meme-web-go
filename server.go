package main

import (
	"embed"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

//go:embed static/*
var embeddedStatic embed.FS

//go:embed .env.example
var envExample string

func main() {
	err := initFiles()
	if err != nil {
		log.Fatal("Error create static files")
	}
	// 读取环境变量
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// 创建gin实例
	r := gin.Default()

	// 定义CORS配置
	CORSConfig := cors.Config{
		AllowAllOrigins:  true, // 允许所有源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(CORSConfig))

	// 注册路由
	registerRoutes(r)

	// 启动服务器
	err = r.Run(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
		return
	}
}

// initFiles 初始化必要的文件和目录
func initFiles() error {
	// 检查并创建 .env 文件
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		if err := os.WriteFile(".env", []byte(envExample), 0644); err != nil {
			return err
		}
		log.Println("Created .env file from example")
	}

	// 检查并创建 static 目录
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		// 从嵌入的文件系统中提取 static 目录
		if err := extractEmbeddedStatic(); err != nil {
			return err
		}
		log.Println("Created static directory from embedded files")
	}

	return nil
}

// extractEmbeddedStatic 从嵌入的文件系统中提取 static 目录
func extractEmbeddedStatic() error {
	return fs.WalkDir(embeddedStatic, "static", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 去掉 static/ 前缀
		relPath, err := filepath.Rel("static", path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join("static", relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := embeddedStatic.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, 0644)
	})
}
