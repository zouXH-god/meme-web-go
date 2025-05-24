package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	// 读取环境变量
	err := godotenv.Load(".env")
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
