package entity

// ShortURLEntity dto сущности
type ShortURLEntity struct {
	ShortURL      string `json:"short_url" db:"url_short"`
	URL           string `json:"original_url" db:"original_url"`
	CorrelationID string `json:"correlation_id" db:"correlation_id"`
	AddedUserID   string `json:"added_user_id" db:"added_user_id"`
	IsDeleted     bool   `json:"-" db:"is_deleted"`
}
