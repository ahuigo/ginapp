package router

import (
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
	prefixPath := "/"
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "https", Host: "httpbin.org", Path: prefixPath})
	// proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "m:4590", Path: "/"})
	if host != "" {
		proxy = &httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.Header.Set("X-Forwarded-Host", r.Host)
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
