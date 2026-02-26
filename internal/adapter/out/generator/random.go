package generator

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type URLGenerator struct {
	rand *rand.Rand
}

func NewURLGenerator() *URLGenerator {
	return &URLGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (u URLGenerator) Generate() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[u.rand.Intn(len(charset))]
	}
	return string(b)
}
