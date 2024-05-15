package utils

import (
	"math/rand"
	"time"
)


const (
	shortCodeLength = 6
	charset         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	shortCode := make([]byte, length)
	for i := range shortCode {
		shortCode[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(shortCode)
}
