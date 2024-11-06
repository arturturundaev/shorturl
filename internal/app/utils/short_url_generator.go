package utils

import (
	"crypto/sha1"
	"encoding/base64"
)

// GenerateShortURL формированиерое краткого url
func GenerateShortURL(url string) string {
	hash := sha1.New()
	hash.Write([]byte(url))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))[:8]

}
