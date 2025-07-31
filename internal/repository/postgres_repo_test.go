//go:build integration

package repository_test

import (
	"shopping_bot/internal/models"
	"shopping_bot/internal/repository"
	"shopping_bot/test/integration"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
func TestUserStateOperations(t *testing.T) {
	db := integration.SetupTestDB(t)
	repo := repository.NewPostgresShoppingRepository(db)

	t.Run("Set and get user state", func(t *testing.T) {
		userID := int64(1)
		testState := &models.UserState{CurrentList: "test_list"}
		// Проверяем запись состояния
		err := repo.SetUserState(userID, testState)
		require.NoError(t, err)
		
		// Проверяем чтение состояния
		retrievedState, err := repo.GetUserState(userID)
		require.NoError(t, err)
		assert.Equal(t, testState.CurrentList, retrievedState.CurrentList)
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		_, err := repo.GetUserState(999)
		assert.ErrorIs(t, err, repository.ErrNoState)
	})
}