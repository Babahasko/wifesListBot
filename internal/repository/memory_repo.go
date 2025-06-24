package repository

import (
	"errors"
	"sync"
)

var (
	ErrListExists = errors.New("shop list already exists")
	ErrNoLists = errors.New("no shop lists for user")
	ErrListNotFound = errors.New("shop list not found")
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
	r.rwMux.Lock()
	defer r.rwMux.Unlock()
	if r.shoppingLists == nil {
		return errors.New("shopping lists is not initialized")
	}
	if _, ok := r.shoppingLists[userID]; !ok {
		r.shoppingLists[userID] = make(map[string]*ShoppingList)
	}
	if _, exists := r.shoppingLists[userID][listName]; exists{
		return ErrListExists
	}
	r.shoppingLists[userID][listName] = &ShoppingList{
		Name: listName,
		Items: make([]*ShoppingItem,0),
	}
	return nil
}

func (r *MemoryShoppingRepository) GetUserLists(userID int64) (map[string]*ShoppingList, error) {
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
	list.Items = append(list.Items, &ShoppingItem{
		ListName: listName,
		Name: itemName,
		Checked: false,
	})
	return nil
}

func (r *MemoryShoppingRepository) GetListItems(userID int64, listName string) ([]*ShoppingItem, error) {
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
	newItems := make([]*ShoppingItem, 0, len(list.Items))
	for _, item := range list.Items {
		if !item.Checked {
			newItems = append(newItems, item)
		}
	}
	list.Items = newItems
	return nil
}
