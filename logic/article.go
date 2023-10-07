package logic

import (
	"errors"
	"news/common/resp"
	"news/global"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Article(ctx *gin.Context) {
	idV := ctx.Param("id")
	id, err := strconv.Atoi(idV)
	if err != nil {
		resp.ExtErr(ctx, err)
		return
	}
	a := global.DB.Article
	do := a.WithContext(ctx)
	data, err := do.Select().Where(a.ID.Eq(int32(id))).First()
	if errors.Is(err, ErrRecordNotFound) {
		resp.ExtErr(ctx, errNotFound)
		return
	} else if err != nil {
		resp.IntErr(ctx, err)
		return
	}
	//浏览数+1
	data.ViewsCount++
	_, err = do.Where(a.ID.Eq(int32(id))).UpdateColumn(a.ViewsCount, data.ViewsCount)
	if err != nil {
		resp.IntErr(ctx, err)
		return
	}
	resp.Ok(ctx, data)
	return
}
