package repository

import (
	"errors"
	"sync"
)

type MemoryShoppingRepository struct {
	rwMux sync.RWMutex

	shoppingLists map[int64]map[string]*ShoppingList
	userStates    map[int64]*UserState
}

func NewMemoryShoppingRepository() *MemoryShoppingRepository {
	return &MemoryShoppingRepository{
		shoppingLists: make(map[int64]map[string]*ShoppingList),
		userStates:    make(map[int64]*UserState),
	}
}

func (r *MemoryShoppingRepository) GetUserState(userID int64) (*UserState, error) {
	r.rwMux.Lock()
	defer r.rwMux.Unlock()

	state, ok := r.userStates[userID]
	if !ok {
		return nil, errors.New("no state for this user")
	}
	return state, nil
}

func (r *MemoryShoppingRepository) SetUserState(userID int64, state *UserState) error {
	r.rwMux.Lock()
	defer r.rwMux.Unlock()

	if r.userStates == nil {
		return errors.New("userStates map is not initialized")
	}
	r.userStates[userID] = state
	return nil
}

func (r *MemoryShoppingRepository) AddShoppingList(userID int64, listName string) error {
	return nil
}
func (r *MemoryShoppingRepository) GetUserLists(userID int64) ([]string, error) {
	listStr := []string{}
	return listStr, nil
}
func (r *MemoryShoppingRepository) DeleteList(userID int64, listName string) error {
	return nil
}

func (r *MemoryShoppingRepository) AddItemToShoppingList(userID int64, listName, itemName string) error {
	return nil
}
func (r *MemoryShoppingRepository) GetListItems(userID int64, listName string) ([]*ShoppingItem, error) {
	return []*ShoppingItem{}, nil
}
func (r *MemoryShoppingRepository) MarkItem(userID int64, listName, itemName string) error {
	return nil
}
func (r *MemoryShoppingRepository) DeleteMarkedItems(userID int64, listName string) error {
	return nil
}
