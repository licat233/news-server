package logic

import (
	"news/common/resp"
	"news/global"

	"github.com/gin-gonic/gin"
)

func Categories(ctx *gin.Context) {
	categories, err := global.DB.Category.WithContext(ctx).Select(global.DB.Category.ID, global.DB.Category.Name).Find()
	resp.New(ctx, categories, err)
}
