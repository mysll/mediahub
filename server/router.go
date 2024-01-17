package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"mediahub/internal/conf"
)

func initRouter(r *gin.Engine) {
	Cors(r)
}

func Cors(r *gin.Engine) {
	config := cors.DefaultConfig()
	//config.AllowAllOrigins = true
	config.AllowOrigins = conf.GetConfig().Cors.AllowOrigins
	config.AllowHeaders = conf.GetConfig().Cors.AllowHeaders
	config.AllowMethods = conf.GetConfig().Cors.AllowMethods
	r.Use(cors.New(config))
}
