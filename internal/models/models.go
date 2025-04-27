package models

// ShortenRequest содержит запрос на сокращение ссылки
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse содержит ответ на запрос сокращения ссылки
type ShortenResponse struct {
	Result string `json:"result"`
}

// BatchShortenRequest содержит запрос на сокращение множества ссылок
type BatchShortenRequest []BatchShortenRequestEntry

// BatchShortenRequestEntry содержит информацию о сокращаемой ссылке
type BatchShortenRequestEntry struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchShortenResponse содержит ответ на запрос сокращение множества ссылок
type BatchShortenResponse []BatchShortenResponseEntry

// BatchShortenResponseEntry содержит информацию о сокращенной ссылке
type BatchShortenResponseEntry struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserShortenedURL содержит информацию о сокращенной ссылке
type UserShortenedURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// UserShortenedURLs содержит набор сокращенных пользователем ссылок
type UserShortenedURLs []UserShortenedURL

// DeleteURLsRequest содержит запрос на удаление сокращенных ссылок
type DeleteURLsRequest []string

// UserDeleteRequest содержит запрос на удаление сокращенной ссылки
type UserDeleteRequest struct {
	UserID          int
	ShortIDToDelete string
}
