package utils

import (
	"math/rand"
	"time"
)

const MAXIMUM_SHORT_URL_COUNT = 1000000 // Maximum number of short URLs
var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rnd     = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func GenerateShortCode() string {
	length := 6 + rnd.Intn(3) // 6, 7, or 8
	code := make([]rune, length)
	for i := range code {
		code[i] = letters[rnd.Intn(len(letters))]
	}
	return string(code)
}
