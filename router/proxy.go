package router

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// curl m:4500/proxy/abc?host=m:4500
func ProxyServer(c *gin.Context) {
	host := c.Query("host")
	path := c.Param("path")
	c.Request.URL.Path = path
	fmt.Println("nihao:", path, "ahuigo!!!", host)
	prefixPath := "/"
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "https", Host: "httpbin.org", Path: prefixPath})
	// proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "m:4590", Path: "/"})
	if host != "" {
		fmt.Println("nihao2:", path, "ahuigo!!!", host)
		proxy = &httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
				r.URL.Scheme = "http"
				r.URL.Path = "/dump" + r.URL.Path
				r.URL.Host = host // with:port
				r.Host = host
			},
			ModifyResponse: func(r *http.Response) error {
				r.Header.Del("Access-Control-Allow-Origin")
				r.Header.Del("Access-Control-Allow-Credentials")
				r.Header.Del("Access-Control-Allow-Headers")
				r.Header.Del("Access-Control-Allow-Methods")
				return nil
			},
		}

	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
