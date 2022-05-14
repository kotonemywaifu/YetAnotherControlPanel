package util

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789_")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
