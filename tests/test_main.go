package tests

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/wangzewang/esman/config"
	"github.com/wangzewang/esman/controllers"
)

type LogSuite struct {
	suite.Suite
	config *viper.Viper
	router *gin.Engine
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.TestMode)
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
func (s *LogSuite) SetupTest() {
	config.Init("test")
	s.config = config.GetConfig()
	s.router = SetupRouter()

}

func (s *LogSuite) TestMain() {
	SetupRouter()
}

func LogTestSuite(t *testing.T) {
	suite.Run(t, new(LogSuite))
}
