package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/internal/domain"
	"github.com/g0shi4ek/v0.1-cargo-comet-back/cometsService/pkg/database"
	"gorm.io/gorm"
)

// Migrate создает или обновляет таблицы в базе данных
func Migrate() error {
	log.Println("Starting database migration...")
	db, err := database.NewPostgresClient()
	if err != nil {
		panic("failed to connect postgres")
	}

	// Автомиграция для основных сущностей
	err = db.AutoMigrate(
		&domain.Comet{},
		&domain.Observation{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}
	seedTestData := os.Getenv("SEED_TEST_DATA")
	if seedTestData == "true" {
		if err := SeedTestData(db); err != nil {
			log.Printf("Warning: Failed to seed test data: %v", err)
		}
	}

	log.Println("Database migration completed successfully")
	return nil
}

// SeedTestData заполняет базу тестовыми данными
func SeedTestData(db *gorm.DB) error {
	log.Println("Seeding test data...")
	// Тестовые кометы
	comets := []domain.Comet{
		{
			UserID:               1,
			Name:                 "Комета Галлея",
			PhotoURL:             "https://example.com/photo1.jpg",
			SemiMajorAxis:        17.8,
			Eccentricity:         0.967,
			RaanDeg:              162.3,
			AscendingNodeLong:    58.42,
			ArgumentOfPerihelion: 111.33,
			TrueAnomalyDeg:       162.3,
		},
		{
			UserID:               1,
			Name:                 "Комета NEOWISE",
			PhotoURL:             "https://example.com/photo2.jpg",
			SemiMajorAxis:        280.0,
			Eccentricity:         0.999,
			RaanDeg:              129.0,
			AscendingNodeLong:    61.0,
			ArgumentOfPerihelion: 37.0,
			TrueAnomalyDeg:       129.0,
		},
		{
			UserID:               2,
			Name:                 "Комета Энке",
			SemiMajorAxis:        2.21,
			Eccentricity:         0.847,
			RaanDeg:              11.78,
			AscendingNodeLong:    334.57,
			ArgumentOfPerihelion: 186.54,
			TrueAnomalyDeg:       11.78,
		},
	}

	// Создаем кометы
	for i := range comets {
		result := db.Create(&comets[i])
		if result.Error != nil {
			return fmt.Errorf("failed to create comet: %w", result.Error)
		}
	}

	// Тестовые наблюдения
	observations := []domain.Observation{
		{
			UserID:         1,
			CometID:        &comets[0].ID,
			RightAscension: 45.67,
			Declination:    23.45,
			ObservedAt:     parseTime("2023-12-01T20:00:00Z"),
		},
		{
			UserID:         1,
			CometID:        &comets[0].ID,
			RightAscension: 46.12,
			Declination:    23.78,
			ObservedAt:     parseTime("2023-12-02T21:00:00Z"),
		},
		{
			UserID:         1,
			CometID:        &comets[0].ID,
			RightAscension: 46.89,
			Declination:    24.12,
			ObservedAt:     parseTime("2023-12-03T22:00:00Z"),
		},
		{
			UserID:         2,
			CometID:        &comets[2].ID,
			RightAscension: 120.45,
			Declination:    -15.67,
			ObservedAt:     parseTime("2023-10-20T19:30:00Z"),
		},
		{
			UserID:         2,
			CometID:        &comets[2].ID,
			RightAscension: 121.23,
			Declination:    -15.89,
			ObservedAt:     parseTime("2023-10-21T20:15:00Z"),
		},
	}

	// Создаем наблюдения
	for i := range observations {
		result := db.Create(&observations[i])
		if result.Error != nil {
			return fmt.Errorf("failed to create observation: %w", result.Error)
		}
	}
	return nil
}

// Вспомогательная функция для парсинга времени
func parseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		log.Printf("Warning: failed to parse time %s: %v", timeStr, err)
		return time.Now()
	}
	return t
}
