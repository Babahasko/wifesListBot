package shop

import (
	"encoding/json"
	"errors"
	"shopping_bot/pkg/callback"
)

// ListCallback данные для callback категории
type ListCallback struct {
	Name     string `json:"name"`
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

func (l *ListCallback) Marshal() ([]byte, error) {
	return json.Marshal(l)
}

func (l *ListCallback) Unmarshal(data []byte) error {
	return json.Unmarshal(data, l)
}

// NewCategoryCallback конструктор для регистрации
func NewListCallback() callback.CallbackData {
	return &ListCallback{}
}