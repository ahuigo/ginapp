package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

// Error returns an error handler
func Error(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			errstr := fmt.Sprintf("error:%+v", r)
			stackstr := errstr + "\n" + string(debug.Stack())
			fmt.Println(stackstr)
			c.String(http.StatusBadRequest, stackstr)
			c.Abort()
		}
	}()
	c.Next()

	if messages := c.Errors.Errors(); len(messages) > 0 && c.Writer.Size() == 0 {
		bodyError := struct {
			Message string `json:"message"`
		}{Message: strings.Join(messages, "\n")}
		if body, err := json.Marshal(bodyError); err != nil {
			fmt.Println("middleware error.go:", err.Error())
		} else {
			c.Writer.Write(body) //nolint: errcheck
		}
	}
}
