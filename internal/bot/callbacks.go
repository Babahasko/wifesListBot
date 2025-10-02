package bot

import (
	"fmt"
	"strings"
)

// Constan Callbacks
const (
	CallbackBackToList = "back_to_list"
	CallbackAddList    = "add_list"
	CallbackClearList  = "clear_list"
	CallbackAddItems   = "add_items"
	CallbackDeleteList = "delete_list"
)

// Complex callbacks
type ListCallbackService struct {
	Prefix string
}

type ListCallbackData struct {
	Name string
}

func NewListCallbackService() *ListCallbackService {
	return &ListCallbackService{
		Prefix: "list",
	}
}

func (c *ListCallbackService) Pack(listName string) (string, error) {
	// Заменяем подчеркивания на безопасные символы
	safeName := strings.ReplaceAll(listName, "_", "-")
	cbStr := fmt.Sprintf("%s_%s", c.Prefix, safeName)
	if len(cbStr) >= 64 {
		return "", fmt.Errorf("cbStr is too long")
	}
	return cbStr, nil
}

func (c *ListCallbackService) Unpack(cbStr string) *ListCallbackData {
	withoutPrefix := strings.TrimPrefix(cbStr, c.Prefix+"_")
	data := strings.Split(withoutPrefix, "_")
	if len(data) < 1 || data[0] == "" {
		// Возвращаем пустые данные если callback поврежден
		return &ListCallbackData{
			Name: "",
		}
	}
	// Восстанавливаем оригинальное имя
	name := strings.ReplaceAll(data[0], "-", "_")
	return &ListCallbackData{
		Name: name,
	}
}

type ItemCallbackService struct {
	Prefix string
}

type ItemCallbackData struct {
	ItemName string
	ListName string
}

func NewItemCallbackService() *ItemCallbackService {
	return &ItemCallbackService{
		Prefix: "item",
	}
}

func (c *ItemCallbackService) Pack(itemName, listName string) (string, error) {
	// Заменяем подчеркивания на безопасные символы
	safeItemName := strings.ReplaceAll(itemName, "_", "-")
	safeListName := strings.ReplaceAll(listName, "_", "-")
	cbStr := fmt.Sprintf("%s_%s_%s", c.Prefix, safeItemName, safeListName)
	if len(cbStr) >= 64 {
		return "", fmt.Errorf("cbStr is too long")
	}
	return cbStr, nil
}

func (c *ItemCallbackService) Unpack(cbStr string) *ItemCallbackData {
	withoutPrefix := strings.TrimPrefix(cbStr, c.Prefix+"_")
	data := strings.Split(withoutPrefix, "_")
	if len(data) < 2 || data[0] == "" || data[1] == "" {
		// Возвращаем пустые данные если callback поврежден
		return &ItemCallbackData{
			ItemName: "",
			ListName: "",
		}
	}
	// Восстанавливаем оригинальные имена
	name := strings.ReplaceAll(data[0], "-", "_")
	listName := strings.ReplaceAll(data[1], "-", "_")
	return &ItemCallbackData{
		ItemName: name,
		ListName: listName,
	}
}
