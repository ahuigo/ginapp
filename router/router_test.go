package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahuigo/gohttptool/httpreq"
	"github.com/gin-gonic/gin"
)

// test gonic context
func TestRouterTest(t *testing.T) {
	e := gin.New()

	// middleware
	e.Use(func(c *gin.Context) {
		c.Next()
		path := c.Request.URL.Path
		header := fmt.Sprintf("%#v", c.Writer.Header())
		fmt.Printf("path:%s,response header:%s\n", path, header)
	})

	// route
	e.GET("/*any", func(c *gin.Context) {
		c.SetSameSite(http.SameSiteNoneMode)
		if c.Request.URL.Host != "x1.com" {
			panic(c.Request.Host)
		}
		if c.Request.Host != "x3.com" {
			panic(c.Request.Host)
		}
		c.SetCookie("count", "1", 1, "", "x.com", false, false)
	})

	// mock request
	// URL, _ := url.Parse("http://x1.com/https")
	// req := &http.Request{
	// 	Method: "GET",
	// 	Header: http.Header{
	// 		"Host": []string{"x2.com"},
	// 	},
	// 	Proto:      "HTTP/1.1",
	// 	ProtoMajor: 1,
	// 	ProtoMinor: 1,
	// 	URL:        URL,
	// 	Host:       "x3.com",
	// }
	req, _ := httpreq.R().SetReq("GET", "http://x1.com/https").SetHeader("Host", "x2.com").SetHost("x3.com").ToRequest()
	writer := httptest.NewRecorder()
	// writer := &http.Response{}
	e.ServeHTTP(writer, req)

}
