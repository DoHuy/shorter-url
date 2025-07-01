package entity

import (
	"time"
)

// ShortURL represents the short_urls table
type ShortURL struct {
	Code        string
	OriginalURL string
	CreatedAt   time.Time
}
