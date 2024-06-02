package main

import (
	"html/template"
	"net/http"

	"github.com/Narushio/qqx5-beatmap-tool/assets"
	"github.com/Narushio/qqx5-beatmap-tool/handler"
	"github.com/Narushio/qqx5-beatmap-tool/view"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	e := gin.New()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	e.StaticFS("/assets", http.FS(assets.FS))
	e.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./assets/favicon.ico")
	})

	e.SetHTMLTemplate(template.Must(template.ParseFS(view.FS, "*.html")))

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
		convertGroup.POST("/mcz_to_xml", convert.MCZToXML)
	}
	return e
}
