package render

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OK(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, data)
}

func Error(ctx *gin.Context, code int, err error) {
	type e struct {
		Message string `json:"message"`
	}

	ctx.JSON(code, e{Message: err.Error()})
}
