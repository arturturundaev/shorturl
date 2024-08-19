package user

import "fmt"

type UrlListItemResponse struct {
	ShortURL  string `json:"short_url"`
	OriginURL string `json:"original_url"`
}

func NewUrlResponse(baseURL, shortURL, originURL string) UrlListItemResponse {
	return UrlListItemResponse{OriginURL: originURL, ShortURL: fmt.Sprintf("%s/%s", baseURL, shortURL)}
}
