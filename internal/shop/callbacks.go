package shop

import (
	"encoding/json"
	"errors"
	"shopping_bot/pkg/callback"
)

// CategoryCallback данные для callback категории
type CategoryCallback struct {
	Name     string `json:"name"`
}

func (c *CategoryCallback) Type() string {
	return "cat"
}

// Validation for category
func (c *CategoryCallback) Validate() error {
	if c.Name == "" {
		return errors.New("category name cannot be empty")
	}
	if len(c.Name) > 30 {
		return errors.New("category name is too long")
	}
	return nil
}

func (c *CategoryCallback) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CategoryCallback) Unmarshal(data []byte) error {
	return json.Unmarshal(data, c)
}

// NewCategoryCallback конструктор для регистрации
func NewCategoryCallback() callback.CallbackData {
	return &CategoryCallback{}
}