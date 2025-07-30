package repository

import (
	"errors"
	"shopping_bot/internal/models"
)

var (
	ErrListExists   = errors.New("shop list already exists")
	ErrNoLists      = errors.New("no shop lists for user")
	ErrListNotFound = errors.New("shop list not found")
	ErrNoState      = errors.New("no state for this user")
)

type ShoppingRepository interface {
	GetUserState(userID int64) (*models.UserState, error)
	SetUserState(userID int64, state *models.UserState) error

	AddShoppingList(userID int64, listName string) error
	GetUserLists(userID int64) (map[string]*models.ShoppingList, error)
	DeleteList(userID int64, listName string) error

	AddItemToShoppingList(userID int64, listName, itemName string) error
	GetListItems(userID int64, listName string) ([]*models.ShoppingItem, error)
	MarkItem(userID int64, listName, itemName string) error
	DeleteMarkedItems(userID int64, listName string) error
}
