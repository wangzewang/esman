package server

import (
	"github.com/gin-gonic/gin"

	"github.com/wangzewang/esman/controllers"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)

	v1 := router.Group("v1")
	{
		logGroup := v1.Group("logs")
		{
			log := new(controllers.LogController)
			logGroup.GET("/all/:task", log.All)
			logGroup.GET("/sse/:task", log.Stream)
		}
	}
	return router

}
