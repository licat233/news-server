package logic

import (
	"errors"
	"news/common/resp"
	"news/global"

	"github.com/gin-gonic/gin"
)

func More(ctx *gin.Context) {
	articleUUIdV := ctx.DefaultQuery("uuid", "")
	if articleUUIdV == "" {
		resp.ExtErr(ctx, errors.New("uuid is empty"))
		return
	}
	a := global.DB.Article
	do := a.WithContext(ctx)
	article, err := do.Select().Where(a.UUID.Eq(articleUUIdV)).First()
	if errors.Is(err, ErrRecordNotFound) {
		resp.ExtErr(ctx, errNotFoundMsg)
		return
	} else if err != nil {
		resp.IntErr(ctx, err)
		return
	}

	category := article.Category
	fields := articleListFields(global.DB)
	// TODO: 随机挑选， 未实现
	articles, err := do.Select(fields...).Where(a.Category.Eq(category)).Order().Order(a.ViewsCount.Desc()).Limit(10).Find()
	resp.New(ctx, articles, err)
	return
}
