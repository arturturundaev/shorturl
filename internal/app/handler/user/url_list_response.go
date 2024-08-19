package user

import "fmt"

type URLListItemResponse struct {
	ShortURL  string `json:"short_url"`
	OriginURL string `json:"original_url"`
}

func NewURLResponse(baseURL, shortURL, originURL string) URLListItemResponse {
	return URLListItemResponse{OriginURL: originURL, ShortURL: fmt.Sprintf("%s/%s", baseURL, shortURL)}
}
