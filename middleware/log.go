package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LogTime(c *gin.Context) {
	// begin time
	t := time.Now()
	// Set example variable
	c.Set("example", "12345")

	uri := c.Request.URL.Path + "?" + c.Request.URL.RawQuery
	log.Println(uri)

	// next  middleware
	c.Next()

	// after request
	latency := time.Since(t)

	// access the statusCode we are sending
	statusCode := c.Writer.Status()
	loginName := getLoginUserFromBody(c)
	log.Println(c.Request.RequestURI, statusCode, latency, loginName, "----------")
}

func getLoginUserFromBody(c *gin.Context) (name string) {
	// bind
	query := struct {
		Username string `json:"username"`
	}{}
	bindPre(c, &query)
	name = query.Username
	return
}

func bindPre(c *gin.Context, req interface{}) error {
	buf, _ := io.ReadAll(c.Request.Body)
	defer func() {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
	}()
	// revert
	c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
	return c.ShouldBind(&req)
}
