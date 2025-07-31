package integration

import (
	"log"
	"shopping_bot/internal/repository"
	"testing"

	"github.com/stretchr/testify/require"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=test password=test dbname=test port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to test database")

	err = MigrateTestDb(db)

	require.NoError(t, err, "Failed to migrate to test database")

	t.Cleanup(func() {
		CleanDatabase(db)
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
	return db
}

func CleanDatabase(db *gorm.DB) {
	tables := []string{"gorm_user_states", "gorm_shopping_lists", "gorm_shopping_items"}
	for _, table := range tables {
		db.Exec("TRUNCATE TABLE " + table + " CASCADE")
	}
}

func MigrateTestDb(db *gorm.DB) error {
	models := []interface{}{
		&repository.GormUserState{},
		&repository.GormShoppingList{},
		&repository.GormShoppingItem{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			log.Fatalf("Migration failed with err: %v", err)
		}
	}
	return nil
}