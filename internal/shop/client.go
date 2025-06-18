package shop

import (
	"fmt"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// A basic handler client to share state across executions.
// Note: This is a very simple layout which uses a shared mutex.
// It is all in-memory, and so will not persist data across restarts.
type ShopClient struct {
	// Use a mutex to avoid concurrency issues.
	// If you use multiple maps, you may want to use a new mutex for each one.
	rwMux sync.RWMutex

	//структура для хранения списков покупок
	shoppingLists map[int64]map[string]*ShoppingList

	// We use a double map to:
	// - map once for the user id
	// - map a second time for the keys a user can have
	// The second map has values of type "any" so anything can be stored in them, for the purpose of this example.
	// This could be improved by using a struct with typed fields, though this would need some additional handling to
	// ensure concurrent safety.
	userData map[int64]map[string]any

	// This struct could also contain:
	// - pointers to database connections
	// - pointers cache connections
	// - localised strings
	// - helper methods for retrieving/caching chat settings
}

type ShoppingList struct {
	Name string `json:"name"`
	Items []string `json:"items"`
}

//TODO: добавить валидацию названия шоппинг листа?
func (c *ShopClient)addShoppingList(ctx *ext.Context, listName string) error {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()
	if c.shoppingLists == nil {
        c.shoppingLists = make(map[int64]map[string]*ShoppingList)
    }

    if _, ok := c.shoppingLists[ctx.EffectiveUser.Id]; !ok {
        c.shoppingLists[ctx.EffectiveUser.Id] = make(map[string]*ShoppingList)
    }

    if _, exists := c.shoppingLists[ctx.EffectiveUser.Id][listName]; exists {
        return fmt.Errorf("shopping list '%s' already exists", listName)
    }

    c.shoppingLists[ctx.EffectiveUser.Id][listName] = &ShoppingList{
        Name:  listName,
        Items: []string{},
    }
	return nil
}

func (c *ShopClient) AddItemToShoppingList(ctx *ext.Context, listName, itemName string) error {
    c.rwMux.Lock()
    defer c.rwMux.Unlock()

    lists, ok := c.shoppingLists[ctx.EffectiveUser.Id]
    if !ok {
        return fmt.Errorf("no shopping lists found for user %d", ctx.EffectiveUser.Id)
    }

    list, ok := lists[listName]
    if !ok {
        return fmt.Errorf("shopping list '%s' not found", listName)
    }

    list.Items = append(list.Items, itemName)
    return nil
}

func (c *ShopClient) GetShoppingList(ctx *ext.Context, listName string) ([]string, error) {
	c.rwMux.RLock()
    defer c.rwMux.RUnlock()

    userID := ctx.EffectiveUser.Id

    // Проверяем, есть ли вообще списки у пользователя
    lists, ok := c.shoppingLists[userID]
    if !ok {
        return nil, fmt.Errorf("no shopping lists found for user %d", userID)
    }

    // Ищем нужный список по имени
    shoppingList, exists := lists[listName]
    if !exists {
        return nil, fmt.Errorf("shopping list '%s' not found for user %d", listName, userID)
    }

    // Возвращаем копию списка товаров, чтобы избежать внешних изменений
    items := make([]string, len(shoppingList.Items))
    copy(items, shoppingList.Items)

    return items, nil
}

func (c *ShopClient) getUserData(ctx *ext.Context, key string) (any, bool) {
	c.rwMux.RLock()
	defer c.rwMux.RUnlock()

	if c.userData == nil {
		return nil, false
	}

	userData, ok := c.userData[ctx.EffectiveUser.Id]
	if !ok {
		return nil, false
	}

	v, ok := userData[key]
	return v, ok
}

func (c *ShopClient) setUserData(ctx *ext.Context, key string, val any) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()

	if c.userData == nil {
		c.userData = map[int64]map[string]any{}
	}

	_, ok := c.userData[ctx.EffectiveUser.Id]
	if !ok {
		c.userData[ctx.EffectiveUser.Id] = map[string]any{}
	}
	c.userData[ctx.EffectiveUser.Id][key] = val
}