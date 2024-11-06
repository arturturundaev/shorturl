package utils

import "fmt"

func ExampleGenerateShortURL() {
	// Строка для получения краткой ссылки
	incomingURL := "http://google.com"

	// краткая ссылка
	shortURL := GenerateShortURL(incomingURL)

	fmt.Print(shortURL)
}
