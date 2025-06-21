package shop

import (
	"errors"
	"shopping_bot/pkg/callback"
)
const (
	CallbackBackToList = "back_to_list"
)
// ==ListCallback==
type ListCallback struct {
	Name string `json:"name"`
}

func (l *ListCallback) Type() string {
	return "list"
}

// Validation for category
func (l *ListCallback) Validate() error {
	if l.Name == "" {
		return errors.New("category name cannot be empty")
	}
	if len(l.Name) > 30 {
		return errors.New("category name is too long")
	}
	return nil
}

// Конструктор для регистрации колбэка
func NewListCallback() callback.CallbackService {
	return &ListCallback{}
}

// ==Items Callback==
type ItemCallback struct {
	ListName string `json:"list_name,omitempty"`
	ItemName string `json:"item_name,omitempty"`
}

func (l *ItemCallback) Type() string {
	return "item"
}

// Validation for category
func (l *ItemCallback) Validate() error {
	if l.ItemName == "" {
		return errors.New("item name cannot be empty")
	}
	if len(l.ItemName) > 30 {
		return errors.New("item name is too long")
	}
	return nil
}

// NewItems callback
func NewItemsCallback() callback.CallbackService {
	return &ItemCallback{}
}
