package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Narushio/qqx5-beatmap-tool/form"
	"github.com/Narushio/qqx5-beatmap-tool/usecase"
)

type HomeController struct{}

func (h HomeController) IndexTmpl(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{"templateName": "osu_to_xml"})
}

func (h HomeController) UpdateXmlBPMTmpl(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{"templateName": "update_xml_bpm"})
}

func (h HomeController) ConvertToXml(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("uploading file: %s", err.Error()))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func (h HomeController) UpdateXmlBPM(c *gin.Context) {
	param := new(form.UpdateXMLBPMParam)
	err := c.Bind(param)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("binding param: %s", err.Error()))
		return
	}

	uc := new(usecase.FileUseCase)
	noteStr, err := uc.UpdateXMLBPM(c, param)
	if noteStr == nil || err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("updating XML: %s", err.Error()))
		return
	}
	c.JSON(http.StatusOK, noteStr)
}
