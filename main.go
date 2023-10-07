package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"news/logic"
	"news/router"
	"os"
	"time"

	_ "news/global"

	"github.com/gin-gonic/gin"
)

func main() {
	handler()
}

func handler() {
	start := time.Now()
	defer func() {
		duration := time.Now().Sub(start)
		fmt.Println("耗時:", duration)
	}()

	err := logic.ArticleHandler(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func api() {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// 添加CORS中间件
	r.Use(corsMiddleware)
	// 设置日志输出到文件
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	router.Load(r)
	r.Run("0.0.0.0:82") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// CORS中间件
func corsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()
}
