package main

import (
	"log"
	"shopping_bot/test/integration"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)
func main() {
	dsn := "host=localhost user=test password=test dbname=test port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	if err := integration.MigrateTestDb(db); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	log.Println("All migrations applied successfully!")
	
	// Проверка существования таблиц
	var tables []string
	db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema='public'").Scan(&tables)
	log.Printf("Existing tables: %v", tables)
}