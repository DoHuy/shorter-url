package api

import (
	"net/http"
	"shorter-rest-api/internal/application/usecase"
	"shorter-rest-api/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

// UserController handles HTTP requests for users
type ShortUrlController struct {
	shortUrlUseCase usecase.ShortUrlUseCase
}

// NewUserController creates a new user controller
func NewShortUrlController(shortUrlUseCase usecase.ShortUrlUseCase) *ShortUrlController {
	return &ShortUrlController{
		shortUrlUseCase: shortUrlUseCase,
	}
}

// RegisterRoutes registers the routes for the user controller
func (c *ShortUrlController) RegisterRoutes(router *gin.Engine) {

	// Register protected  routes

	router.GET("/api/shortlinks/:id", c.GetShortByCode)
	router.POST("/api/shortlinks", c.CreateShortUrl)
	router.GET("/shortlinks/:id", c.Redirect)
}

// GetShortByCode gets a shorturl by ID
// @Summary      Get shorturl by ID
// @Description  Retrieves a specific shorturl by its ID
// @Tags         shorturl
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "short id"
// @Success      200  {object}  dto.GetShortUrlResponse
// @Failure      400  "id is required"
// @Failure 	 500 "Internal Server Error"
// @Router       /api/shortlinks/{id} [get]
func (c *ShortUrlController) GetShortByCode(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	result, err := c.shortUrlUseCase.GetShortUrlByCode(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// Redirect shorturl by ID
// @Summary      Redirect to original URL
// @Description  Redirects to the original URL for the given short code
// @Tags         shorturl
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "short id"
// @Success      200  {object}  dto.GetShortUrlResponse
// @Failure      400  "Bad Request - Invalid input"
// @Failure      302 "Found - Redirects to original URL"
// @Failure 	 500 "Internal Server Error"
// @Router       /shortlinks/{id} [get]
func (c *ShortUrlController) Redirect(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	result, err := c.shortUrlUseCase.GetShortUrlByCode(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, result.OriginalUrl)
}

// CreateShortUrl creates a new shorturl
// @Summary      Create shorturl
// @Description  Creates a new shorturl
// @Tags         shorturl
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateRequest  true  "URL object"
// @Success      201  {object}  dto.CreateResponse
// @Failure      400  "Bad Request - Invalid input"
// @Failure      409  "Conflict - Short URL already exists"
// @Failure      500  "Internal Server Error"
// @Router       /api/shortlinks [post]
func (c *ShortUrlController) CreateShortUrl(ctx *gin.Context) {
	var shortUrl dto.CreateRequest
	if err := ctx.ShouldBindJSON(&shortUrl); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate duplicate short URL
	isDuplicate, err := c.shortUrlUseCase.ValidateDuplicateShortUrl(shortUrl.OriginalUrl)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isDuplicate {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Short URL already exists"})
		return
	}

	result, err := c.shortUrlUseCase.CreateShortUrl(ctx, &shortUrl)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, result)
}
