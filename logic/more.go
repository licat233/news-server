package logic

import (
	"errors"
	"news/common/resp"
	"news/global"
	"strconv"

	"github.com/gin-gonic/gin"
)

func More(ctx *gin.Context) {
	articleIdV := ctx.DefaultQuery("id", "0")
	articleId, err := strconv.Atoi(articleIdV)
	if err != nil {
		resp.ExtErr(ctx, err)
		return
	}

	a := global.DB.Article
	do := a.WithContext(ctx)
	article, err := do.Select().Where(a.ID.Eq(int32(articleId))).First()
	if errors.Is(err, ErrRecordNotFound) {
		resp.ExtErr(ctx, errNotFound)
		return
	} else if err != nil {
		resp.IntErr(ctx, err)
		return
	}
	categoryId := article.ColumnID
	fields := articleListFields(global.DB)
	articles, err := do.Select(fields...).Where(a.ColumnID.Eq(categoryId)).Order().Order(a.ViewsCount.Desc()).Limit(10).Find()
	resp.New(ctx, articles, err)
	return
}
