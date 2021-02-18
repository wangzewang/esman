package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelloController struct{}

func (h HelloController) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	return
}
