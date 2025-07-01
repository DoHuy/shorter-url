package dto

// LoginRequest represents the login request payload
type CreateRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
}

type GetShortUrlResponse struct {
	ID          string `json:"id"`
	OriginalUrl string `json:"original_url"`
	CreatedAt   string `json:"created_at"`
}

type CreateResponse struct {
	ID       string `json:"id"`
	ShortUrl string `json:"short_url"`
}
