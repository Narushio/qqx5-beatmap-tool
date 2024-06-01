package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Home struct{}

func NewHome() *Home {
	return &Home{}
}

func (h *Home) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index", nil)
}

func (h *Home) UpdateXmlBPMTmpl(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "update_xml_bpm", nil)
}

func (h *Home) MCZToXMLTmpl(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "mcz_to_xml", nil)
}
