package router

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func StatusServer(ctx *gin.Context) {
	code, _ := strconv.Atoi(ctx.Param("code"))
	ctx.Status(code)
}

// r.GET("/status/:code", StatusServer)
func init() {
	handlers = append(handlers, Handler{
		method:  "GET",
		path:    "/status/:code",
		handler: StatusServer,
	})
}
