package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/qq524967747/ratelimiter/pkg"
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
	buf := make([]byte, 1*pkg.KB)
	TotalLimit := pkg.NewRateLimiter(pkg.TransRate(1*pkg.KB), 2)
	limitReader := pkg.NewLimitReaderWithLimiter(TotalLimit, f, false)
	start := time.Now()
	_, copyErr := io.CopyBuffer(c.Writer, limitReader, buf)
	if copyErr != nil {
		log.Println(copyErr.Error())
	}
	println("耗时：", time.Now().Sub(start).String())
}
