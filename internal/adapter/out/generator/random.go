package generator

import (
	"sync"
	"time"
)

const (
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLen = uint64(len(charset))
	keyLength  = 8
	epoch      = int64(1704067200000)
)

type SnowflakeGenerator struct {
	mu      sync.Mutex
	counter uint64
	lastMs  int64
}

func NewURLGenerator() *SnowflakeGenerator {
	return &SnowflakeGenerator{}
}

func (s *SnowflakeGenerator) Generate() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()
	if now == s.lastMs {
		s.counter++
	} else {
		s.lastMs = now
		s.counter = 0
	}

	id := (uint64(now-epoch) << 22) | s.counter
	return encodeBase62(id, keyLength)
}

func encodeBase62(id uint64, length int) string {
	result := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		result[i] = charset[id%charsetLen]
		id /= charsetLen
	}
	return string(result)
}
