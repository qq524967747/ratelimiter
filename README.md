ratelimiter

golang的数据下载限流工具

使用方法：

```go
package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/qq524967747/ratelimiter"
)

func main() {
	app := gin.New()
	app.GET("", copy)
	app.Run(":8080")
}

func copy(c *gin.Context) {
	f, err := os.Open("source/test.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// 适当调整buf和rate速率
	buf := make([]byte, 100*ratelimiter.KB)
	TotalLimit := ratelimiter.NewRateLimiter(ratelimiter.TransRate(100*ratelimiter.KB), 2)
	limitReader := ratelimiter.NewLimitReaderWithLimiter(TotalLimit, f, false)
	start := time.Now()
	_, copyErr := io.CopyBuffer(c.Writer, limitReader, buf)
	if copyErr != nil {
		log.Println(copyErr.Error())
	}
	println("耗时：", time.Now().Sub(start).String())
}
```