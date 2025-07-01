package usecase

import (
	"context"
	"fmt"
	"shorter-rest-api/internal/config"
	"shorter-rest-api/internal/domain/dto"
	"shorter-rest-api/internal/domain/entity"
	"shorter-rest-api/internal/infrastructure/cache"
	"shorter-rest-api/internal/infrastructure/utils"
	"strings"
	"time"
)

// ShortUrlUseCase defines the interface for shortUrl use cases
type ShortUrlUseCase interface {
	GetShortUrlByCode(ctx context.Context, code string) (*dto.GetShortUrlResponse, error)
	CreateShortUrl(ctx context.Context, url *dto.CreateRequest) (*dto.CreateResponse, error)
	ValidateDuplicateShortUrl(originalUrl string) (bool, error)
}

type shortUrlUseCase struct {
	cacheService cache.IRedisCache
	cfg          *config.Config
}

// NewShortUrlUseCase creates a new shortUrl use case
func NewShortUrlUseCase(config *config.Config, cacheService cache.IRedisCache) ShortUrlUseCase {
	return &shortUrlUseCase{
		cacheService: cacheService,
		cfg:          config,
	}
}

func (uc *shortUrlUseCase) ValidateDuplicateShortUrl(originalUrl string) (bool, error) {

	// Get short URL by code
	isExist, err := uc.cacheService.Exists(originalUrl)
	if err != nil && !strings.Contains(err.Error(), "redigo: nil returned") {
		return false, fmt.Errorf("failed to find short url: %w", err)
	}
	if !isExist {
		return true, nil // No duplicate found
	}
	return isExist, nil
}

func (uc *shortUrlUseCase) GetShortUrlByCode(ctx context.Context, code string) (*dto.GetShortUrlResponse, error) {

	// Get short URL by code
	shortUrl, err := uc.cacheService.Get(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find short url: %w", err)
	}

	// Map to response DTO
	response := &dto.GetShortUrlResponse{
		ID:          shortUrl.Code,
		OriginalUrl: shortUrl.OriginalURL,
		CreatedAt:   shortUrl.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

// CreateShortUrl creates a new shortUrl
func (uc *shortUrlUseCase) CreateShortUrl(ctx context.Context, shortUrl *dto.CreateRequest) (*dto.CreateResponse, error) {

	// check maximum short URL count follow configure from
	// initialization simplest will hardcode is 1 million saved keys
	count, err := uc.cacheService.CountKeysByPattern("short_urls:*")
	if err != nil {
		return nil, fmt.Errorf("failed to count short URLs: %w", err)
	}
	if count >= uc.cfg.MaximumShortUrlCount {
		return nil, fmt.Errorf("maximum short URL count reached: %d", uc.cfg.MaximumShortUrlCount)
	}
	// Create a new short URL entity
	newShortUrl := &entity.ShortURL{
		OriginalURL: shortUrl.OriginalUrl,
		CreatedAt:   time.Now(), // Set the current time as CreatedAt
	}

	// generate code for the short URL
	for {
		shortCode := utils.GenerateShortCode() // Assume this function generates a random short code
		newShortUrl.Code = shortCode
		exists, _ := uc.cacheService.Exists(fmt.Sprintf("short_urls:%s", newShortUrl.Code))
		if !exists {
			break
		}
	}

	// Store the new short URL in the cache
	if err := uc.cacheService.Set(fmt.Sprintf("short_urls:%s", newShortUrl.Code), *newShortUrl, uc.cfg.Expiration); err != nil {
		return nil, fmt.Errorf("failed to create short URL: %w", err)
	}
	// Store the original URL in the cache with the short code as the key
	if err := uc.cacheService.Set(newShortUrl.OriginalURL, *newShortUrl, uc.cfg.Expiration); err != nil {
		return nil, fmt.Errorf("failed to create short URL: %w", err)
	}

	return &dto.CreateResponse{
		ID:       newShortUrl.Code,
		ShortUrl: fmt.Sprintf("http://localhost:%s/shortlinks/%s", uc.cfg.Server.Port, newShortUrl.Code),
	}, nil

}
