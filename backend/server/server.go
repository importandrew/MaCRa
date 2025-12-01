package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/models"
)

func NewServer(origins ...string) *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies(nil)

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = append(corsCfg.AllowOrigins, "https://"+models.DefaultDomain)
	corsCfg.AllowOrigins = append(corsCfg.AllowOrigins, origins...)
	corsCfg.AllowCredentials = true

	router.Use(cors.New(corsCfg))

	return router
}
