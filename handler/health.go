package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Health struct{}

func NewHealth() *Health {
	return &Health{}
}

func (h *Health) Index(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{"msg": "ok"})
}
