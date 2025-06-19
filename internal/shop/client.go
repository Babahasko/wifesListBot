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

	// структура для долгосрочного хранения списков покупок
	shoppingLists map[int64]map[string]*ShoppingList

	// структура для отслеживания состояния пользователя
	userStates map[int64]*UserState
}

type UserState struct {
	CurrentList string `json:"state"`//текущий лист покупок редактируемый пользователем
}

type ShoppingList struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
}

func (c *ShopClient) GetUserState(ctx *ext.Context) *UserState {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()

	if c.userStates == nil {
		return nil
	}
	return c.userStates[ctx.EffectiveUser.Id]
}

func (c *ShopClient) SetUserState(ctx *ext.Context, state *UserState) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()
	if c.userStates == nil {
		c.userStates = make(map[int64]*UserState)
	}
	c.userStates[ctx.EffectiveUser.Id] = state
}

func (c *ShopClient) GetCurrentList(ctx *ext.Context) string {
	state := c.GetUserState(ctx)
	if state == nil {
		return ""
	}
	return state.CurrentList
}

func (c *ShopClient) SetCurrentListName(ctx *ext.Context, listName string) {
	state := c.GetUserState(ctx)
	if state == nil {
		state = &UserState{}
	}
	state.CurrentList = listName
	c.SetUserState(ctx, state)
}

// TODO: добавить валидацию названия шоппинг листа?
func (c *ShopClient) addShoppingList(ctx *ext.Context, listName string) error {
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

func (c *ShopClient) GetListItems(ctx *ext.Context, listName string) ([]string, error) {
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
