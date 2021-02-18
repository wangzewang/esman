package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HelloController struct{}

func (h HelloController) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	return
}
