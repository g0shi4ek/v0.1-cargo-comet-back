package main

import (
	"log"
	"os"

	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/cmd/clients"
	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/domain"
	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/handlers"
	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/repository"
	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// Получение переменных окружения
	appPort := os.Getenv("APP_PORT")
	runMigrations := os.Getenv("RUN_MIGRATIONS")

	cometRepo := repository.NewCometsRepository()

	// Выполнение миграций
	if runMigrations == "true" {
		if err := Migrate(); err != nil {
			log.Fatal("Migration failed:", err)
		}
	}

	// Получение адресов сервисов из переменных окружения
	authServiceAddr := os.Getenv("AUTH_SERVICE_ADDR")   //"localhost:50051"
	orbitServiceAddr := os.Getenv("ORBIT_SERVICE_ADDR") //"localhost:50052"

	// Инициализация клиентов
	var orbitCalcClient domain.IOrbitCalculationClient
	var authClient domain.IAuthClient

	// В зависимости от окружения используем реальные или mock клиенты
	useRealClients := os.Getenv("USE_REAL_CLIENTS")

	realAuthClient, err := clients.NewRealAuthClient(authServiceAddr)
	if err != nil {
		log.Fatal("Failed to create auth client:", err)
	}
	authClient = realAuthClient

	if useRealClients == "true" {
		log.Println("Using real gRPC clients")

		// Реальные клиенты
		_, err := clients.NewRealOrbitCalculationClient(orbitServiceAddr)
		if err != nil {
			log.Fatal("Failed to create orbit calculation client:", err)
		}
		// orbitCalcClient = realOrbitClient
		orbitCalcClient = NewMockOrbitCalculationClient()

	} else {
		orbitCalcClient = NewMockOrbitCalculationClient()
		authClient = NewMockAuthClient()
	}

	// Инициализация сервиса
	cometsService := service.NewCometsService(cometRepo, orbitCalcClient)

	// Настройка роутера
	router := gin.Default()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	// Настройка маршрутов
	handlers.SetupRoutes(router, cometsService, authClient)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "comets-service",
		})
	})

	// Запуск сервера
	log.Printf("Starting server on port %s", appPort)
	if err := router.Run(":" + appPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// CORSMiddleware настройка CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
