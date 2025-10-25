package database

import (
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func postgresConfigFromEnv() string {
	// Получение переменных окружения
	_ = godotenv.Load()
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword +
		" dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=UTC"
	return dsn
}

func NewPostgresClient() (*gorm.DB, error) {
	dsn := postgresConfigFromEnv()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
