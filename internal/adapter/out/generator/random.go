package generator

import (
	"math/rand"
	"sync"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type URLGenerator struct {
	rand *rand.Rand
	mu   *sync.Mutex
}

func NewURLGenerator() *URLGenerator {
	return &URLGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		mu:   &sync.Mutex{},
	}
}

func (u URLGenerator) Generate() string {
	b := make([]byte, 8)
	u.mu.Lock()
	for i := range b {
		b[i] = charset[rand.Intn(len(charset)-1)]
	}
	u.mu.Unlock()
	return string(b)
}
