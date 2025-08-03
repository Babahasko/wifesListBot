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

func TestShoppingItemOperations(t *testing.T) {
    db := integration.SetupTestDB(t)
    repo := repository.NewPostgresShoppingRepository(db)
    
    userID := int64(489)
    listName := "Продукты"

    t.Run("Add item to existing shopping list", func(t *testing.T) {
        // Сначала создаем список покупок
        err := repo.AddShoppingList(userID, listName)
        require.NoError(t, err)

        // Добавляем item в список
        itemName := "Молоко"
        err = repo.AddItemToShoppingList(userID, listName, itemName)
        require.NoError(t, err)

        // Проверяем, что item добавлен
        items, err := repo.GetListItems(userID, listName)
        require.NoError(t, err)
        assert.Len(t, items, 1)
        assert.Equal(t, itemName, items[0].Name)
        assert.False(t, items[0].Checked)
    })

    t.Run("Add item to non-existent shopping list", func(t *testing.T) {
        nonExistentList := "Несуществующий список"
        err := repo.AddItemToShoppingList(userID, nonExistentList, "Хлеб")
        assert.ErrorIs(t, err, repository.ErrListNotFound)
    })

    t.Run("Add multiple items with same name", func(t *testing.T) {
        itemName := "Яблоки"
        
        // Добавляем первый item
        err := repo.AddItemToShoppingList(userID, listName, itemName)
        require.NoError(t, err)

        // Добавляем второй item с тем же именем
        err = repo.AddItemToShoppingList(userID, listName, itemName)
        require.NoError(t, err)

        // Проверяем, что оба item'а добавлены
        items, err := repo.GetListItems(userID, listName)
        require.NoError(t, err)
        assert.GreaterOrEqual(t, len(items), 2)
        
        // Подсчитываем items с нужным именем
        count := 0
        for _, item := range items {
            if item.Name == itemName {
                count++
            }
        }
        assert.GreaterOrEqual(t, count, 2)
    })

    t.Run("Get items from non-existent list", func(t *testing.T) {
        nonExistentList := "Несуществующий список"
        _, err := repo.GetListItems(userID, nonExistentList)
        assert.ErrorIs(t, err, repository.ErrListNotFound)
    })

    t.Run("Mark item as checked", func(t *testing.T) {
        itemName := "Хлеб"
        
        // Добавляем item
        err := repo.AddItemToShoppingList(userID, listName, itemName)
        require.NoError(t, err)

        // Отмечаем item
        err = repo.MarkItem(userID, listName, itemName)
        require.NoError(t, err)

        // Проверяем, что item отмечен
        items, err := repo.GetListItems(userID, listName)
        require.NoError(t, err)
        
        var foundItem *models.ShoppingItem
        for _, item := range items {
            if item.Name == itemName {
                foundItem = item
                break
            }
        }
        
        require.NotNil(t, foundItem)
        assert.True(t, foundItem.Checked)
    })

    t.Run("Mark non-existent item", func(t *testing.T) {
        nonExistentItem := "Несуществующий товар"
        err := repo.MarkItem(userID, listName, nonExistentItem)
        assert.Error(t, err)
        // Проверяем, что это ошибка, связанная с не найденным item
        assert.Contains(t, err.Error(), "not found")
    })

    t.Run("Mark item in non-existent list", func(t *testing.T) {
        nonExistentList := "Несуществующий список"
        err := repo.MarkItem(userID, nonExistentList, "Хлеб")
        assert.ErrorIs(t, err, repository.ErrListNotFound)
    })

    t.Run("Delete marked items", func(t *testing.T) {
        // Добавляем несколько items
        itemsToAdd := []string{"Молоко", "Хлеб", "Яйца", "Сахар"}
        for _, itemName := range itemsToAdd {
            err := repo.AddItemToShoppingList(userID, listName, itemName)
            require.NoError(t, err)
        }

        // Отмечаем некоторые items
        itemsToMark := []string{"Яйца", "Сахар"}
        for _, itemName := range itemsToMark {
            err := repo.MarkItem(userID, listName, itemName)
            require.NoError(t, err)
        }

        // Удаляем отмеченные items
        err := repo.DeleteMarkedItems(userID, listName)
        require.NoError(t, err)

        // Проверяем, что остались только неотмеченные items
        items, err := repo.GetListItems(userID, listName)
        require.NoError(t, err)
        
        // Создаем мапу оставшихся items для удобства проверки
        remainingItems := make(map[string]bool)
        for _, item := range items {
            remainingItems[item.Name] = true
        }

        // Проверяем, что отмеченные items удалены
        assert.NotContains(t, remainingItems, "Яйца")
        assert.NotContains(t, remainingItems, "Сахар")
        
        // Проверяем, что неотмеченные items остались
        assert.Contains(t, remainingItems, "Молоко")
        assert.Contains(t, remainingItems, "Хлеб")
    })

    t.Run("Delete marked items from non-existent list", func(t *testing.T) {
        nonExistentList := "Несуществующий список"
        err := repo.DeleteMarkedItems(userID, nonExistentList)
        assert.ErrorIs(t, err, repository.ErrListNotFound)
    })

    t.Run("Delete marked items when no items are marked", func(t *testing.T) {
        newListName := "Пустой список"
        
        // Создаем новый список
        err := repo.AddShoppingList(userID, newListName)
        require.NoError(t, err)

        // Добавляем items, но не отмечаем их
        err = repo.AddItemToShoppingList(userID, newListName, "Товар1")
        require.NoError(t, err)

        err = repo.AddItemToShoppingList(userID, newListName, "Товар2")
        require.NoError(t, err)

        // Удаляем отмеченные items (их нет)
        err = repo.DeleteMarkedItems(userID, newListName)
        require.NoError(t, err)

        // Проверяем, что все items остались
        items, err := repo.GetListItems(userID, newListName)
        require.NoError(t, err)
        assert.Len(t, items, 2)
    })
}