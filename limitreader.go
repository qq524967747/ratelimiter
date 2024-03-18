package ratelimiter

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
)

// NewLimitReader creates a LimitReader.
// src: reader
// rate: bytes/second
func NewLimitReader(src io.Reader, rate int64, calculateMd5 bool) *LimitReader {
	return NewLimitReaderWithLimiter(newRateLimiterWithDefaultWindow(rate), src, calculateMd5)
}

// NewLimitReaderWithLimiter creates LimitReader with a
// src: reader
// rate: bytes/second
func NewLimitReaderWithLimiter(rl *RateLimiter, src io.Reader, calculateMd5 bool) *LimitReader {
	var md5sum hash.Hash
	if calculateMd5 {
		md5sum = md5.New()
	}
	return &LimitReader{
		Src:     src,
		Limiter: rl,
		md5sum:  md5sum,
	}
}

// NewLimitReaderWithMD5Sum creates LimitReader with a md5 sum.
// src: reader
// rate: bytes/second
func NewLimitReaderWithMD5Sum(src io.Reader, rate int64, md5sum hash.Hash) *LimitReader {
	return NewLimitReaderWithLimiterAndMD5Sum(src, newRateLimiterWithDefaultWindow(rate), md5sum)
}

// NewLimitReaderWithLimiterAndMD5Sum creates LimitReader with rateLimiter and md5 sum.
// src: reader
// rate: bytes/second
func NewLimitReaderWithLimiterAndMD5Sum(src io.Reader, rl *RateLimiter, md5sum hash.Hash) *LimitReader {
	return &LimitReader{
		Src:     src,
		Limiter: rl,
		md5sum:  md5sum,
	}
}

func newRateLimiterWithDefaultWindow(rate int64) *RateLimiter {
	return NewRateLimiter(TransRate(rate), 2)
}

// LimitReader reads stream with
type LimitReader struct {
	Src     io.Reader
	Limiter *RateLimiter
	md5sum  hash.Hash
}

func (lr *LimitReader) Read(p []byte) (n int, err error) {
	n, e := lr.Src.Read(p)
	if e != nil && e != io.EOF {
		return n, e
	}
	if n > 0 {
		if lr.md5sum != nil {
			lr.md5sum.Write(p[:n])
		}
		lr.Limiter.AcquireBlocking(int64(n))
	}
	return n, e
}

// Md5 calculates the md5 of all contents read.
func (lr *LimitReader) Md5() string {
	if lr.md5sum != nil {
		return GetMd5Sum(lr.md5sum, nil)
	}
	return ""
}

// GetMd5Sum gets md5 sum as a string and appends the current hash to b.
func GetMd5Sum(md5 hash.Hash, b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}
