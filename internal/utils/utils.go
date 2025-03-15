package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

func GenerateRandomString(n int) (string, error) {
	result := make([]byte, n)
	charsetLength := big.NewInt(int64(len(charset)))
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}
