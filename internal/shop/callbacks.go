package shop

import (
	"fmt"
	"strings"
)

// Constan Callbacks
const (
	CallbackBackToList = "back_to_list"
	CallbackAddList    = "add_list"
	CallbackClearList  = "clear_list"
	CallbackNoItems    = "no_items"
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
	cbStr := fmt.Sprintf("%s_%s", c.Prefix, listName)
	if len(cbStr) >= 64 {
		return "", fmt.Errorf("cbStr is too long")
	}
	return cbStr, nil
}

func (c *ListCallbackService) Unpack(cbStr string) *ListCallbackData {
	withoutPrefix := strings.TrimPrefix(cbStr, c.Prefix + "_")
	data := strings.Split(withoutPrefix, "_")
	name := data[0]
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
	cbStr := fmt.Sprintf("%s_%s_%s", c.Prefix, itemName, listName)
	if len(cbStr) >= 64 {
		return "", fmt.Errorf("cbStr is too long")
	}
	return cbStr, nil
}

func (c *ItemCallbackService) Unpack(cbStr string) *ItemCallbackData {
	withoutPrefix := strings.TrimPrefix(cbStr, c.Prefix + "_")
	data := strings.Split(withoutPrefix, "_")
	name := data[0]
	listName := data[1]
	return &ItemCallbackData{
		ItemName: name,
		ListName: listName,
	}
}