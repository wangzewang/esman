package controllers

import (
	"io"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"

	"github.com/wangzewang/esman/es"
)

type LogController struct{}

func (l LogController) All(c *gin.Context) {

	task := c.Params.ByName("task")
	if task == "" {
		c.JSON(http.StatusOK, gin.H{"task": task, "status": "no value"})
		return
	}

	res := es.NewEsQuery(task)

	c.JSON(http.StatusOK, gin.H{"task": task, "status": res})
}

func (l LogController) Stream(c *gin.Context) {

	task := c.Params.ByName("task")
	if task == "" {
		c.JSON(http.StatusOK, gin.H{"task": task, "status": "no value"})
		return
	}

	stop := int64(0)
	resp := make(chan string)
	go es.NewEsStreamQuery(task, &stop, resp)

	c.Stream(func(w io.Writer) bool {
	queryLoop:
		for {
			select {
			case msg := <-resp:
				c.SSEvent("message", msg)
				return true
			case <-c.Request.Context().Done():
				atomic.StoreInt64(&stop, 1)
				break queryLoop
			}
		}
		return false
	})
}
