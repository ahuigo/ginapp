package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RespBody(c *gin.Context) {
	w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
	c.Writer = w

	// rewrite header
	w.Header().Set("Server", "ahui server")

	// rewrite code
	w.ResponseWriter.WriteHeader(205)
	w.WriteHeaderNow() // call before c.Next

	c.Next()
	body := w.body.String()

	header := dumpHeader(w.Header())
	fmt.Printf("response header:\n%s\n", header)
	fmt.Printf("response body len:%d\n", len(body))
}

func dumpHeader(header http.Header) string {
	var res strings.Builder
	headers := sortHeader(header)
	for _, kv := range headers {
		res.WriteString(kv[0] + ": " + kv[1] + "\n")
	}
	return res.String()
}

// sortHeaders
func sortHeader(header http.Header) [][2]string {
	headers := [][2]string{}
	for k, vs := range header {
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

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer //cache
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r responseBodyWriter) WriteString(s string) (n int, err error) {
	r.body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}
