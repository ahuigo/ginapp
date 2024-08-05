package stat

import (
	"time"

	"github.com/gin-gonic/gin"
)

func ConsumeMemory(c *gin.Context) {
	data := make([]byte, 1024*1024*1024)
	for i := range data {
		data[i] = byte(i)
	}
	time.Sleep(1 * time.Second)
	data = append(data, 123)
	// 释放data 内存占用
	// 1. debug/pprof/heap?debug=1 会看到 heapAlloc/heapInuse 减少(会触发一次gc)
	// 2. go tool pprof debug/pprof/heap 则不会触发gc, 即使触发了gc 只会看到`heap_space` 减少,  `alloc_space` 未减少(它其实是从程序启动开始累计分配的所有内存TotalAlloc)
	data = data[:100]
	_ = data
	// case "large-memory":
	c.String(200, "cost large-memory")
}
