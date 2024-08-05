package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	mid "ginapp/middleware"
	statHandler "ginapp/router/stat"
	tpls "ginapp/router/tpl-server"

	"github.com/DeanThompson/ginpprof"
	"github.com/ahuigo/goos-tools/gonic"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	method  string // http.MethodPut ....
	path    string
	handler gin.HandlerFunc
}

var (
	handlers = []Handler{}
)

func Register(r *gin.Engine, staticFS bool, path404 string) {
	if staticFS {
		staticFsHandler(r, path404)
		return
	}
	// set tpl router // template
	tpls.TplRouter(r)
	// automatically add routers for net/http/pprof
	ginpprof.Wrap(r)
	// swagger
	RegisterSwagger(r)
	// midddleware
	r.Use(mid.LogTime, mid.CORS, mid.Error)
	// staticFS
	r.Static("/js", "./js")
	// register handlers
	for _, h := range handlers {
		r.Handle(h.method, h.path, h.handler)
	}

	// router
	r.GET("/gorm/insert", insertHandler)
	// curl m:4500/f/r/my/a.txt  -o a.txt
	r.GET("/f/r/*path", fileReadHandler)
	// curl m:4500/f/w/my/a.txt -F 'file1=@go.mod' -F 'name=alex'
	r.POST("/f/w/*path", fileWriteHandler)
	r.GET("/api/panic", panicApi)
	r.GET("/api/conf", confApi)
	r.GET("/dump/*anypath", DumpServer)
	r.POST("/dump/*anypath", DumpServer)
	r.GET("/redirect/:code", RedirectServer)
	r.POST("/redirect/form", RedirectServer)
	r.GET("/echo/:size", EchoServer)
	r.GET("/bind/*anypath", BindServer)
	r.POST("/bind/*anypath", BindServer)
	r.PUT("/bind/*anypath", BindFileServer)
	r.GET("/sleep/:second", sleepFunc)
	r.GET("/cookie", cookieServer)
	r.GET("/cpu/:second", cpuFunc)
	r.GET("/json/map", jsonMapFunc)
	r.GET("/proxy/*path", ProxyServer)
	r.POST("/proxy/*path", ProxyServer)
	r.GET("/stream", streamApi)
	r.GET("/stat/os", gonic.OsStat)
	r.GET("/stat/net", gonic.NetStat)
	r.GET("/stat/os/cosume-mem", statHandler.ConsumeMemory)
	// r.Any("/bind/*anypath", BindServer)
}

func jsonMapFunc(c *gin.Context) {
	m := map[string][]byte{
		"status": []byte("running!中国"),
	}
	c.JSON(http.StatusOK, m)
}

func EchoServer(c *gin.Context) {
	size, _ := strconv.Atoi(c.Param("size"))
	if size > 1e5 {
		size = 1e5
	}
	msg := strings.Repeat("中", size)
	c.String(http.StatusOK, msg)
}

func cpuFunc(c *gin.Context) {
	seconds, _ := strconv.Atoi(c.Param("second"))
	n := longRun(seconds)
	msg := fmt.Sprintf("sleep second: %v s, n=%d\n", seconds, n)
	c.JSON(http.StatusOK, msg)
}

func longRun(seconds int) int {
	now := time.Now()
	end_time := now.Add(time.Duration(seconds) * time.Second)
	n := 0
	for ; end_time.After(now); now = time.Now() {
		for i := 0; i < 1e8; i += 1 {
			n += 1
		}
	}
	return n
}

func sleepFunc(c *gin.Context) {
	// this is depend on: ReadTimeout, WriteTimeout, HandlerTimeout in main.go
	seconds, _ := strconv.Atoi(c.Param("second"))
	fmt.Printf("%vs sleep!\n", seconds)
	time.Sleep(time.Duration(seconds) * time.Second)
	fmt.Printf("%vs passed!\n", seconds)
	c.JSON(http.StatusOK, "sleep second: "+c.Param("second"))
}
