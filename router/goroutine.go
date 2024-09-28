package router

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var gowork_num int32 = 0

func gowork(seconds int) {
	defer atomic.AddInt32(&gowork_num, -1)
	atomic.AddInt32(&gowork_num, 1)

	duration := time.Second * time.Duration(seconds)
	time.Sleep(duration)
}
func goHandler(c *gin.Context) {
	size, _ := strconv.Atoi(c.Param("size"))
	seconds, _ := strconv.Atoi(c.Query("seconds"))
	if size > 1000 {
		c.String(http.StatusBadRequest, "size too large>1000")
		return
	}
	if seconds > 600 {
		seconds = 600
	} else if seconds <= 0 {
		seconds = 30
	}
	new_num := int32(size) - atomic.LoadInt32(&gowork_num)
	for i := 0; i < int(new_num); i++ {
		go gowork(seconds)
	}
	c.String(http.StatusOK, fmt.Sprintf("new goroutines(%d)?seconds=%d ", new_num, seconds))
}
