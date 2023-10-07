package router

import (
	"net/http"
	"news/logic"

	"github.com/gin-gonic/gin"
)

func Load(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/test", logic.Test)
	api := r.Group("/api")
	{
		api.GET("/news", logic.News)
		api.GET("/news/more", logic.More)
		api.GET("/news/:id", logic.Article)
		api.GET("/categories", logic.Categories)
		api.GET("/updateImgage", logic.UpdateImage)
		// api.GET("/update", gin.BasicAuth(gin.Accounts{
		// 	"licat": "licat233",
		// }), logic.UpdateArticle)
	}
}
