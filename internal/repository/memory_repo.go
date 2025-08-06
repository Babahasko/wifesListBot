package repository

import (
	"errors"
	"shopping_bot/internal/models"
	"sync"
)

type MemoryShoppingRepository struct {
	rwMux sync.RWMutex

	shoppingLists map[int64]map[string]*models.ShoppingList
	userStates    map[int64]*models.UserState
}

func NewMemoryShoppingRepository() *MemoryShoppingRepository {
	return &MemoryShoppingRepository{
		shoppingLists: make(map[int64]map[string]*models.ShoppingList),
		userStates:    make(map[int64]*models.UserState),
	}
}

func (r *MemoryShoppingRepository) GetUserState(userID int64) (*models.UserState, error) {
	r.rwMux.Lock()
	defer r.rwMux.Unlock()

	state, ok := r.userStates[userID]
	if !ok {
		return nil, errors.New("no state for this user")
	}
	return state, nil
}

func (r *MemoryShoppingRepository) SetUserState(userID int64, state *models.UserState) error {
	r.rwMux.Lock()
	defer r.rwMux.Unlock()

	if r.userStates == nil {
		return errors.New("userStates map is not initialized")
	}
	r.userStates[userID] = state
	return nil
}

func (r *MemoryShoppingRepository) UpdateUserState(userID int64, state *models.UserState) error {
	r.rwMux.Lock()
	defer r.rwMux.Unlock()
	r.userStates[userID] = state
	return nil
}

func (r *MemoryShoppingRepository) AddShoppingList(userID int64, listName string) error {
	r.rwMux.Lock()
	defer r.rwMux.Unlock()
	if r.shoppingLists == nil {
		return errors.New("shopping lists is not initialized")
	}
	if _, ok := r.shoppingLists[userID]; !ok {
		r.shoppingLists[userID] = make(map[string]*models.ShoppingList)
	}
	if _, exists := r.shoppingLists[userID][listName]; exists{
		return ErrListExists
	}
	r.shoppingLists[userID][listName] = &models.ShoppingList{
		Name: listName,
		Items: make([]*models.ShoppingItem,0),
	}
	return nil
}

func (r *MemoryShoppingRepository) GetUserLists(userID int64) (map[string]*models.ShoppingList, error) {
	r.rwMux.Lock()
	defer r.rwMux.Unlock()
	lists, ok := r.shoppingLists[userID]
	if !ok {
		return nil, ErrNoLists
	}
	return lists, nil
}

func (r *MemoryShoppingRepository) DeleteList(userID int64, listName string) error {
	r.rwMux.Lock()
	defer r.rwMux.Unlock()
	userLists, ok := r.shoppingLists[userID]
	if !ok {
		return ErrNoLists
	}
	if _, exists := userLists[listName]; !exists {
		return ErrListNotFound
	}
	delete(userLists, listName)
	return nil
}

func (r *MemoryShoppingRepository) AddItemToShoppingList(userID int64, listName, itemName string) error {
	r.rwMux.Lock()
	defer r.rwMux.Unlock()
	userLists, ok := r.shoppingLists[userID]
	if !ok {
		return ErrNoLists
	}
	list, ok := userLists[listName]
	if !ok {
		return ErrListNotFound
	}
	list.Items = append(list.Items, &models.ShoppingItem{
		ListName: listName,
		Name: itemName,
		Checked: false,
	})
	return nil
}

func (r *MemoryShoppingRepository) GetListItems(userID int64, listName string) ([]*models.ShoppingItem, error) {
	r.rwMux.Lock()
	defer r.rwMux.Unlock()
	userLists, ok := r.shoppingLists[userID]
	if !ok {
		return nil, ErrNoLists
	}
	list, ok := userLists[listName]
	if !ok {
		return nil, ErrListNotFound
	}
	return list.Items, nil
}

func (r *MemoryShoppingRepository) MarkItem(userID int64, listName, itemName string) error {
	r.rwMux.Lock()
	defer r.rwMux.Unlock()
	userLists, ok := r.shoppingLists[userID]
	if !ok {
		return ErrNoLists
	}
	list, ok := userLists[listName]
	if !ok {
		return ErrListNotFound
	}
	for _, item := range list.Items {
		if item.Name == itemName {
			item.Checked = !item.Checked
		}
	}
	return nil
}
func (r *MemoryShoppingRepository) DeleteMarkedItems(userID int64, listName string) error {
	r.rwMux.Lock()
	defer r.rwMux.Unlock()
	userLists, ok := r.shoppingLists[userID]
	if !ok {
		return ErrNoLists
	}
	list, ok := userLists[listName]
	if !ok {
		return ErrListNotFound
	}
	newItems := make([]*models.ShoppingItem, 0, len(list.Items))
	for _, item := range list.Items {
		if !item.Checked {
			newItems = append(newItems, item)
		}
	}
	list.Items = newItems
	return nil
}
