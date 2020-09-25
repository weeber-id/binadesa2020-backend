package tools

import (
	"math/rand"
	"time"
)

// RandomString generator
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, length)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
