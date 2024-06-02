package main

import (
	"github.com/Narushio/qqx5-beatmap-tool/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	e := gin.New()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.Static("/assets", "./assets")
	e.LoadHTMLGlob("view/*")

	health := handler.NewHealth()
	e.GET("/health", health.Index)

	home := handler.NewHome()
	{
		e.GET("/", home.Index)
		e.GET("/update_xml_bpm", home.UpdateXmlBPMTmpl)
		e.GET("/mcz_to_xml", home.MCZToXMLTmpl)
	}

	v1 := e.Group("api/v1")

	convert := handler.NewConvert()
	convertGroup := v1.Group("/convert")
	{
		convertGroup.POST("/osu_to_xml", convert.OSUToXML)
		convertGroup.POST("/update_xml_bpm", convert.UpdateXMLBPM)
	}
	return e
}
