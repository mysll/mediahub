package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"mediahub/internal/conf"
)

func Init(e *gin.Engine) {
	Cors(e)
	g := e.Group(conf.GetConfig().App.SiteURL)

	g.Any("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}

func Cors(e *gin.Engine) {
	config := cors.DefaultConfig()
	//config.AllowAllOrigins = true
	config.AllowOrigins = conf.GetConfig().Cors.AllowOrigins
	config.AllowHeaders = conf.GetConfig().Cors.AllowHeaders
	config.AllowMethods = conf.GetConfig().Cors.AllowMethods
	e.Use(cors.New(config))
}
