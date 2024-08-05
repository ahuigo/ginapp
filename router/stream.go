package router

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/mattn/go-colorable"
)

// r.GET("/stream", func(c *gin.Context) {

func streamApi(c *gin.Context) {
	// gin.DefaultWriter = colorable.NewColorableStderr()
	chanStream := make(chan int, 10)
	go func() {
		defer close(chanStream)
		for i := 0; i < 5; i++ {
			chanStream <- i
			time.Sleep(time.Second * 1)
		}
	}()
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-chanStream; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})

}
