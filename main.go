package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "shorter-rest-api/docs"
	"shorter-rest-api/internal/application/usecase"
	"shorter-rest-api/internal/config"
	"shorter-rest-api/internal/infrastructure/cache"
	"shorter-rest-api/internal/interfaces/api"

	"github.com/gin-gonic/gin"
)

// @title          			   Shorter API Documentation
// @version         		   1.0
// @description     		   Swagger Shorter API Documentation.
// @host            		   localhost:8080
// @BasePath       			   /
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up database connection
	inMemDB, err := cache.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create use cases
	shorterUseCase := usecase.NewShortUrlUseCase(cfg, inMemDB)

	// Create Gin router
	router := gin.New()

	// Register swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register controllers
	shorterController := api.NewShortUrlController(shorterUseCase)

	// Register routes
	shorterController.RegisterRoutes(router)

	// Add health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Create server
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server is running on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shut down server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
