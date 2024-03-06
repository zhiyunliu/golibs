package xrandom

import (
	"math/rand"
	"time"
)

const (
	letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	lettersLen int = len(letters)
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Str 返回n个字符长度的字符串
func Str(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(lettersLen)]
	}
	return string(b)
}

// Range 返回[min,max)中间的任意数字
func Range(min, max int) int {
	if max <= min {
		return min
	}
	return min + rand.Intn(max-min)
}
