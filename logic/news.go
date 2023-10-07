package logic

import (
	"news/common/resp"
	"news/global"
	"strconv"

	"github.com/gin-gonic/gin"
)

const NEWEST = "newest"
const HOT = "hot"

func News(ctx *gin.Context) {
	category := ctx.DefaultQuery("category", NEWEST) //最新

	minIdV := ctx.DefaultQuery("minId", "0")
	minId, err := strconv.Atoi(minIdV)
	if err != nil {
		resp.ExtErr(ctx, err)
		return
	}

	pageSizeV := ctx.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(pageSizeV)
	if err != nil {
		resp.ExtErr(ctx, err)
		return
	}
	if pageSize > 100 {
		pageSize = 100
	}

	a := global.DB.Article
	fields := articleListFields(global.DB)
	do := a.WithContext(ctx).Select(fields...).Limit(pageSize)
	if minId > 0 {
		do = do.Where(a.ID.Lt(int32(minId)))
	}

	switch category {
	case NEWEST:
		do = do.Order(a.PublishAt.Desc(), a.ID.Desc())
	case HOT:
		do = do.Order(a.ViewsCount.Desc(), a.ID.Desc())
	default:
		do = do.Where(a.Category.Eq(category)).Order(a.ID.Desc())
	}

	articles, err := do.Find()
	resp.New(ctx, articles, err)
}
