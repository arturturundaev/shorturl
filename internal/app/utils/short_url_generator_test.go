package utils

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkGenerateShortURL(b *testing.B) {
	b.Run("generate by length url = 10", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			GenerateShortURL(generateRandomString(10))
		}
	})

	b.Run("generate by length url = 40", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			GenerateShortURL(generateRandomString(40))
		}
	})
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}
