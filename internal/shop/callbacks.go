package shop

import (
	"encoding/json"
	"errors"
	"shopping_bot/pkg/callback"
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

// вспомогательные методы для десериализации данных в колбэк
func (l *ListCallback) Marshal() ([]byte, error) {
	return json.Marshal(l)
}

func (l *ListCallback) Unmarshal(data []byte) error {
	return json.Unmarshal(data, l)
}

// Конструктор для регистрации колбэка
func NewListCallback() callback.CallbackData {
	return &ListCallback{}
}

// ==Items Callback==
type ItemCallback struct {
	ListName string `json:"list_name"`
	Name string `json:"name"`
}

func (l *ItemCallback) Type() string {
	return "item"
}

// Validation for category
func (l *ItemCallback) Validate() error {
	if l.Name == "" {
		return errors.New("category name cannot be empty")
	}
	if len(l.Name) > 30 {
		return errors.New("category name is too long")
	}
	return nil
}

// вспомогательные методы для десериализации данных в колбэк
func (l *ItemCallback) Marshal() ([]byte, error) {
	return json.Marshal(l)
}

func (l *ItemCallback) Unmarshal(data []byte) error {
	return json.Unmarshal(data, l)
}

// NewItems callback
func NewItemsCallback() callback.CallbackData {
	return &ItemCallback{}
}
