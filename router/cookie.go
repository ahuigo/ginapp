package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func cookieServer(c *gin.Context) {
	// get count
	count, _ := c.Cookie("count")
	i, _ := strconv.Atoi(count)
	count = fmt.Sprintf("%d", i+1)
	// set count
	c.SetCookie("count", count, 3600, "/", GetCookieDomain(c), false, true)
	fmt.Println("count:", count)
	c.JSON(http.StatusOK, count)
}
