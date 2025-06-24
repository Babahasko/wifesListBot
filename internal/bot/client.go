package bot

import (
	"errors"
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
	CurrentList string `json:"state"` //текущий лист покупок редактируемый пользователем
}

type ShoppingList struct {
	Name  string          `json:"name"`
	Items []*ShoppingItem `json:"items"`
}

type ShoppingItem struct {
	ListName string `json:"list_name"`
	Name     string `json:"item_name"`
	Checked  bool   `json:"checked"`
}

func (c *ShopClient) getUserState(ctx *ext.Context) *UserState {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()

	if c.userStates == nil {
		return nil
	}
	return c.userStates[ctx.EffectiveUser.Id]
}

func (c *ShopClient) setUserState(ctx *ext.Context, state *UserState) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()
	if c.userStates == nil {
		c.userStates = make(map[int64]*UserState)
	}
	c.userStates[ctx.EffectiveUser.Id] = state
}

func (c *ShopClient) getCurrentList(ctx *ext.Context) string {
	state := c.getUserState(ctx)
	if state == nil {
		return ""
	}
	return state.CurrentList
}

func (c *ShopClient) setCurrentListName(ctx *ext.Context, listName string) {
	state := c.getUserState(ctx)
	if state == nil {
		state = &UserState{}
	}
	state.CurrentList = listName
	c.setUserState(ctx, state)
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
		Items: make([]*ShoppingItem, 0),
	}
	return nil
}

func (c *ShopClient) getUserLists(ctx *ext.Context) ([]string, error) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()
	userID := ctx.EffectiveUser.Id
	if c.shoppingLists == nil {
		return nil, errors.New(ErrorNoLists)
	}
	userLists, ok := c.shoppingLists[userID]
	if !ok {
		return nil, fmt.Errorf("no shopping lists found for user: %d", userID)
	}
	listNames := make([]string, 0, len(userLists))
	for name := range userLists {
		listNames = append(listNames, name)
	}
	return listNames, nil
}

func (c *ShopClient) deleteList(ctx *ext.Context, listName string) error {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()

	userID := ctx.EffectiveUser.Id

	userLists := c.shoppingLists[userID]
	_, exists := userLists[listName]
	if !exists {
		return fmt.Errorf("shopping list with name %q not found for userId %d", listName, userID)
	}
	delete(userLists, listName)
	return nil
}

func (c *ShopClient) addItemToShoppingList(ctx *ext.Context, listName, itemName string) error {
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

	list.Items = append(list.Items, &ShoppingItem{
		ListName: listName,
		Name:     itemName,
		Checked:  false,
	})
	return nil
}

func (c *ShopClient) getListItems(ctx *ext.Context, listName string) ([]*ShoppingItem, error) {
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
	items := make([]*ShoppingItem, len(shoppingList.Items))
	copy(items, shoppingList.Items)

	return items, nil
}

func (c *ShopClient) markItem(ctx *ext.Context, listName, itemName string) error {
	userID := ctx.EffectiveUser.Id
	shopList := c.shoppingLists[userID][listName]
	for _, item := range shopList.Items {
		if item.Name == itemName {
			item.Checked = !item.Checked
		}
	}
	return nil
}

func (c *ShopClient) deleteMarkItems(ctx *ext.Context, listName string) error {
	userID := ctx.EffectiveUser.Id
	shopList := c.shoppingLists[userID][listName]
	newShopList := make([]*ShoppingItem, 0, len(shopList.Items))
	for _, item := range shopList.Items {
		if !item.Checked {
			newShopList = append(newShopList, item)
		}
	}
	shopList.Items = newShopList
	return nil
}
