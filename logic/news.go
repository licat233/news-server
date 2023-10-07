package logic

import (
	"news/common/resp"
	"news/global"
	"strconv"

	"github.com/gin-gonic/gin"
)

func News(ctx *gin.Context) {
	columnIdV := ctx.DefaultQuery("columnId", "-2") //最新
	columnId, err := strconv.Atoi(columnIdV)
	if err != nil {
		resp.ExtErr(ctx, err)
		return
	}

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

	viewsCountV := ctx.DefaultQuery("viewsCount", "0")
	viewsCount, err := strconv.Atoi(viewsCountV)
	if err != nil {
		resp.ExtErr(ctx, err)
		return
	}
	a := global.DB.Article
	fields := articleListFields(global.DB)
	do := a.WithContext(ctx).Select(fields...).Limit(pageSize)
	if minId > 0 && columnId != -1 {
		do = do.Where(a.ID.Lt(int32(minId)))
	}
	if columnId > 0 {
		do = do.Where(a.ColumnID.Eq(strconv.Itoa(columnId))).Order(a.ID.Desc())
	} else if columnId == -1 { //熱門
		if viewsCount > 0 {
			do = do.Where(a.ViewsCount.Lt(int32(viewsCount)))
		}
		do = do.Order(a.ViewsCount.Desc(), a.ID.Desc())
	} else { //最新
		do = do.Order(a.ID.Desc())
	}

	articles, err := do.Find()
	resp.New(ctx, articles, err)
}
