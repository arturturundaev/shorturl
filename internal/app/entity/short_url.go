package entity

type ShortURLEntity struct {
	ShortURL      string `json:"short_url" db:"url_short"`
	URL           string `json:"original_url" db:"original_url"`
	CorrelationId string `json:"correlation_id" db:"correlation_id"`
}
