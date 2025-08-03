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

func TestShoppingListOperations(t *testing.T) {
	db := integration.SetupTestDB(t)
	repo := repository.NewPostgresShoppingRepository(db)

	t.Run("Add get and delete ShoppingList", func(t *testing.T) {
		userID := int64(489)
		listName1 := "Еда"
		listName2 := "Авто"
		listName3 := "Любимой"

		err := repo.AddShoppingList(userID, listName1)
		require.NoError(t, err)

		err = repo.AddShoppingList(userID, listName2)
		require.NoError(t, err)

		err = repo.AddShoppingList(userID, listName3)
		require.NoError(t, err)

		userLists, err := repo.GetUserLists(userID)
		require.NoError(t, err)

		// Проверяем, что список не пустой
		assert.NotEmpty(t, userLists)

		// Проверяем, что в списке есть элемент с нужным именем
		assert.Contains(t, userLists, listName1)

		// Удаляем список
		err = repo.DeleteList(userID, listName1)
		require.NoError(t, err)

		// Проверяем, что список удалён
		userLists, err = repo.GetUserLists(userID)
		require.NoError(t, err)
		assert.NotContains(t, userLists, listName1)
	})
}