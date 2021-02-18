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
		initGroup := v1.Group("hello")
		{
			hello := new(controllers.HelloController)
			initGroup.GET("/hello", hello.Hello)
		}
	}
	return router

}
