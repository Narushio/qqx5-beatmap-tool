package server

import (
	"github.com/Narushio/qqx5-beatmap-tool/controller"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	e := gin.New()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.Static("/assets", "./assets")
	e.LoadHTMLGlob("template/*")

	health := new(controller.HealthController)
	e.GET("/health", health.Status)

	home := new(controller.HomeController)
	{
		e.GET("/", home.IndexTmpl)
		e.GET("/update_xml_bpm", home.UpdateXmlBPMTmpl)
	}

	v1 := e.Group("api/v1")
	{
		v1.POST("/convert_to_xml", home.ConvertToXml)
		v1.POST("/update_xml_bpm", home.UpdateXmlBPM)
	}
	return e
}
