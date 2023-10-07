package resp

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(ctx *gin.Context, data any, err error) {
	if err != nil {
		IntErr(ctx, err)
	} else {
		Ok(ctx, data)
	}
}

func Ok(ctx *gin.Context, data any) {
	ctx.AbortWithStatusJSON(http.StatusOK, data)
}

func ExtErr(ctx *gin.Context, err error) {
	ctx.AbortWithError(http.StatusBadRequest, err)
}

func IntErr(ctx *gin.Context, err error) {
	log.Println(err.Error())
	ctx.AbortWithStatus(http.StatusInternalServerError)
}
