package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/handlers"
	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/repository"
	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Инициализация логгера
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Получение переменных окружения
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "mydbpass")
	dbName := getEnv("DB_NAME", "cometdb")
	appPort := getEnv("APP_PORT", "8080")
	runMigrations := getEnv("RUN_MIGRATIONS", "true")
	seedTestData := getEnv("SEED_TEST_DATA", "true")

	// Подключение к базе данных
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + 
		" dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=UTC"
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to database successfully")

	// Выполнение миграций
	if runMigrations == "true" {
		if err := Migrate(db); err != nil {
			log.Fatal("Migration failed:", err)
		}
	}

	// Заполнение тестовыми данными
	if seedTestData == "true" {
		if err := SeedTestData(db); err != nil {
			log.Printf("Warning: Failed to seed test data: %v", err)
		}
	}

	// Инициализация зависимостей
	cometRepo := repository.NewCometsRepository(db)
	
	// Инициализация клиентов (заглушки для демонстрации)
	orbitCalcClient := NewMockOrbitCalculationClient()
	fileStorageClient := NewMockFileStorageClient()
	authClient := NewMockAuthClient()

	// Инициализация сервиса
	cometsService := service.NewCometsService(cometRepo, orbitCalcClient, fileStorageClient)

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
			"status": "ok",
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

// getEnv получение переменной окружения с значением по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}