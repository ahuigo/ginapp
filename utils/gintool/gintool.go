package gintool

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var httpsDomainPattern = regexp.MustCompile(`(domain1|domain2)\.com$`)

func IsHttpsRequest(c *gin.Context) bool {
	matched := httpsDomainPattern.MatchString(c.Request.Host)
	return (c.Request.TLS != nil ||
		c.GetHeader("X-Forwarded-Proto") == "https" ||
		strings.HasSuffix(c.Request.Host, ".github.com") ||
		c.Request.URL.Scheme == "https" ||
		matched)
}
