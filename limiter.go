package ratelimiter

import (
	"io"
)

func CopyBuffer(w *io.Writer, r io.Reader, speed int64) error {
	// 适当调整buf和rate速率
	buf := make([]byte, speed)
	TotalLimit := NewRateLimiter(TransRate(speed), 2)
	limitReader := NewLimitReaderWithLimiter(TotalLimit, r, false)
	_, copyErr := io.CopyBuffer(*w, limitReader, buf)
	if copyErr != nil {
		return copyErr
	}
	return nil
}
