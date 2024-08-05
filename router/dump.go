package router

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// EchoHandler _
func DumpServer(c *gin.Context) {
	sendCookie(c)
	sendHeaders(c)
	sendBody(c)
}

func GetCookieDomain(ctx *gin.Context) string {
	host := getOriginDomain(ctx)
	if host == "" {
		host = ctx.Request.Host
	}
	if i := strings.Index(host, ":"); i >= 0 {
		host = host[:i]
	}
	cookieDomain := host
	// not IP address, use parent domain
	hostSegs := strings.Split(host, ".")
	start := 0
	if len(hostSegs) > 2 {
		start = len(hostSegs) - 2
	}
	cookieDomain = strings.Join(hostSegs[start:], ".")
	return cookieDomain
}

func getOriginDomain(c *gin.Context) string {
	url := c.Request.Header.Get("Origin")
	if url == "" {
		return ""
	}
	u, err := neturl.Parse(url)
	if err != nil {
		return ""
	}
	return u.Host
}

// sendBody
func sendBody(c *gin.Context) {
	// dump header
	res := c.Request.Method + " " + //c.Request.URL.String() +" "+
		c.Request.Host +
		c.Request.URL.Path + "?" + c.Request.URL.RawQuery +
		" proto:" + strconv.Itoa(c.Request.ProtoMajor) +
		" proto:" + c.Request.Proto + " clientip:" +
		c.ClientIP() + "\n"
	res += dumpRequestHeader(c.Request) + "\n"
	res += fmt.Sprintf("c.Request.Host(include port):%s\n", c.Request.Host)
	res += fmt.Sprintf("c.Request.RequestURI:%s\n", c.Request.RequestURI)
	res += fmt.Sprintf("c.Request.URL.Scheme:%s\n", c.Request.URL.Scheme)
	res += fmt.Sprintf("c.Request.URL.Hostname():%s\n", c.Request.URL.Hostname())
	res += fmt.Sprintf("c.Request.URL.Host(invalid):%s\n", c.Request.URL.Host)
	res += fmt.Sprintf("c.Request.URL.PATH:%s\n", c.Request.URL.Path)
	res += fmt.Sprintf("c.FullPath():%s\n", c.FullPath())
	res += fmt.Sprintf("origin:%s\n", c.Request.Header.Get("Origin"))
	res += fmt.Sprintf("originDomain:%s\n", getOriginDomain(c))
	res += fmt.Sprintf("cookieDomain:%s\n", GetCookieDomain(c))
	res += fmt.Sprintf("referer:%s\n", c.Request.Header.Get("referer"))

	// dump body
	buf, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(buf)) // important!!
	res += string(buf)

	println(res)
	c.String(http.StatusOK, res)
}

func dumpRequestHeader(req *http.Request) string {
	var res strings.Builder
	headers := sortHeaders(req)
	for _, kv := range headers {
		res.WriteString(kv[0] + ": " + kv[1] + "\n")
	}
	return res.String()
}

// sortHeaders
func sortHeaders(request *http.Request) [][2]string {
	headers := [][2]string{}
	for k, vs := range request.Header {
		for _, v := range vs {
			headers = append(headers, [2]string{k, v})
		}
	}
	n := len(headers)
	for i := 0; i < n; i++ {
		for j := n - 1; j > i; j-- {
			jj := j - 1
			h1, h2 := headers[j], headers[jj]
			if h1[0] < h2[0] {
				headers[jj], headers[j] = headers[j], headers[jj]
			}
		}
	}
	return headers
}

// sendHeaders
func sendHeaders(c *gin.Context) {
	c.Writer.Header().Set("Remote", "echo-server by ahuigo")
	//c.Writer.Header().Set("Location", "http://baidu1.com")
}

// sendCookie
func sendCookie(c *gin.Context) {
	// c.Request.URL.
	c.SetSameSite(http.SameSiteLaxMode)
	hostname := strings.Split(c.Request.Host, ":")[0]
	countStr, _ := c.Cookie("count")
	if countStr == "" {
		countStr = "1"
	} else {
		count, _ := strconv.Atoi(countStr)
		countStr = strconv.Itoa(count + 1)
	}

	// Set-Cookie: count=1; Path=/; Domain=ahui.io; Max-Age=172800
	c.SetCookie("count", countStr, 86400, "", hostname, false, false)
	c.SetCookie("count", countStr, 86400, "", ".ahuigo1.io", false, false)
	c.SetCookie("count", countStr, 86400, "", "ahuigo2.io", false, false)
	// fmt.Printf("h:%#v\n", c.Header)
}
