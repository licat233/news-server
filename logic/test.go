package logic

import (
	"news/common/resp"
	"news/global"

	"github.com/gin-gonic/gin"
)

func Test(ctx *gin.Context) {
	a := global.DB.Article
	do := global.DB.Article.WithContext(ctx).Select(a.ID, a.Title, a.AuthorName, a.Cover, a.PublishedAt, a.ColumnID, a.ColumnName, a.ViewsCount)
	// do := global.DB.Article.WithContext(ctx).Select()
	do = do.Where(a.ID.Eq(21000))
	data, err := do.Find()
	resp.New(ctx, data, err)
}
