package tests

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wangzewang/esman/config"
	"github.com/wangzewang/esman/controllers"
)

func Test(t *testing.T) { Testing(t) }

var _ = Suite(&initSuite{})

type initSuite struct {
	config *viper.Viper
	router *gin.Engine
}

func (s *initSuite) SetUpTest(c *C) {
	config.Init("test")
	s.config = config.GetConfig()
	s.router = SetupRouter()
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.TestMode)
	health := new(controllers.HealthController)
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

func TestMain(m *testing.M) {
	SetupRouter()
}
