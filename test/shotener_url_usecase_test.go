package test

import (
	"errors"
	"shorter-rest-api/internal/domain/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository or service dependencies
type MockShortUrlRepo struct {
	mock.Mock
}

func (m *MockShortUrlRepo) GetByCode(id string) (*dto.GetShortUrlResponse, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.GetShortUrlResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockShortUrlRepo) ExistsByOriginalUrl(url string) (bool, error) {
	args := m.Called(url)
	return args.Bool(0), args.Error(1)
}

func (m *MockShortUrlRepo) Create(req *dto.CreateRequest) (*dto.GetShortUrlResponse, error) {
	args := m.Called(req)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.GetShortUrlResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// Example ShortUrlUseCase struct for illustration
type ShortUrlUseCase struct {
	repo *MockShortUrlRepo
}

func (u *ShortUrlUseCase) GetShortUrlByCode(id string) (*dto.GetShortUrlResponse, error) {
	return u.repo.GetByCode(id)
}

func (u *ShortUrlUseCase) ValidateDuplicateShortUrl(url string) (bool, error) {
	exists, err := u.repo.ExistsByOriginalUrl(url)
	return !exists, err
}

func (u *ShortUrlUseCase) CreateShortUrl(req *dto.CreateRequest) (*dto.GetShortUrlResponse, error) {
	return u.repo.Create(req)
}

func TestGetShortUrlByCode_Success(t *testing.T) {
	mockRepo := new(MockShortUrlRepo)
	usecase := &ShortUrlUseCase{repo: mockRepo}

	expected := &dto.GetShortUrlResponse{
		ID:          "abc123",
		OriginalUrl: "https://example.com",
		CreatedAt:   "2025-07-01T00:00:00Z",
	}
	mockRepo.On("GetByCode", "abc123").Return(expected, nil)

	result, err := usecase.GetShortUrlByCode("abc123")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetShortUrlByCode_NotFound(t *testing.T) {
	mockRepo := new(MockShortUrlRepo)
	usecase := &ShortUrlUseCase{repo: mockRepo}

	mockRepo.On("GetByCode", "notfound").Return(nil, errors.New("not found"))

	result, err := usecase.GetShortUrlByCode("notfound")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestValidateDuplicateShortUrl_Duplicate(t *testing.T) {
	mockRepo := new(MockShortUrlRepo)
	usecase := &ShortUrlUseCase{repo: mockRepo}

	mockRepo.On("ExistsByOriginalUrl", "https://example.com").Return(true, nil)

	isUnique, err := usecase.ValidateDuplicateShortUrl("https://example.com")
	assert.NoError(t, err)
	assert.False(t, isUnique)
}

func TestValidateDuplicateShortUrl_Unique(t *testing.T) {
	mockRepo := new(MockShortUrlRepo)
	usecase := &ShortUrlUseCase{repo: mockRepo}

	mockRepo.On("ExistsByOriginalUrl", "https://unique.com").Return(false, nil)

	isUnique, err := usecase.ValidateDuplicateShortUrl("https://unique.com")
	assert.NoError(t, err)
	assert.True(t, isUnique)
}

func TestCreateShortUrl_Success(t *testing.T) {
	mockRepo := new(MockShortUrlRepo)
	usecase := &ShortUrlUseCase{repo: mockRepo}

	req := &dto.CreateRequest{OriginalUrl: "https://example.com"}
	expected := &dto.GetShortUrlResponse{
		ID:          "abc123",
		OriginalUrl: "https://example.com",
		CreatedAt:   "2025-07-01T00:00:00Z",
	}
	mockRepo.On("Create", req).Return(expected, nil)

	result, err := usecase.CreateShortUrl(req)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateShortUrl_Error(t *testing.T) {
	mockRepo := new(MockShortUrlRepo)
	usecase := &ShortUrlUseCase{repo: mockRepo}

	req := &dto.CreateRequest{OriginalUrl: "https://fail.com"}
	mockRepo.On("Create", req).Return(nil, errors.New("create error"))

	result, err := usecase.CreateShortUrl(req)
	assert.Error(t, err)
	assert.Nil(t, result)
}
