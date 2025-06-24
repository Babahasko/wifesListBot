package service

import (
	"shopping_bot/internal/repository"
)

type ShoppingService struct {
	repo repository.ShoppingRepository
}

func NewShopService(repo repository.ShoppingRepository) *ShoppingService {
	return &ShoppingService{
		repo: repo,
	}
}

// TODO: Add validation for item and list name`s`

func (s *ShoppingService) SetCurrentList(userID int64, listName string) error {
	state, _ := s.repo.GetUserState(userID)
	if state == nil {
		state = &repository.UserState{}
	}
	state.CurrentList = listName
	err :=s.repo.SetUserState(userID, state)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShoppingService) GetCurrentList(userID int64) (string, error) {
	state, err := s.repo.GetUserState(userID)
	if err != nil {
		return "", err
	}
	return state.CurrentList, nil
}

func (s *ShoppingService) AddShoppingList(userID int64, listName string) error {
	return s.repo.AddShoppingList(userID, listName)
}

func (s *ShoppingService) GetUserLists(userID int64) ([]string, error) {
	userLists, err := s.repo.GetUserLists(userID)
	if err != nil {
		return nil, err
	}

	strListNames := make([]string, 0, len(userLists))
	for _, list := range userLists {
		strListNames = append(strListNames, list.Name)
	}
	return strListNames, nil
}

func (s *ShoppingService) DeleteList(userID int64, listName string) error {
	return s.repo.DeleteList(userID, listName)
}

func (s *ShoppingService) AddItemToShoppingList(userID int64, listName, itemName string) error {
	return s.repo.AddItemToShoppingList(userID, listName, itemName)
}

func (s *ShoppingService) GetListItems(userID int64, listName string) ([]*repository.ShoppingItem, error) {
	listItems, err := s.repo.GetListItems(userID, listName)
	if err != nil {
		return nil, err
	}
	return listItems, nil
}

func (s *ShoppingService) MarkItem(userID int64, listName, itemName string) error {
	return s.repo.MarkItem(userID, listName, itemName)
}

func (s *ShoppingService) DeleteMarkedItems(userID int64, listName string) error {
	return s.repo.DeleteMarkedItems(userID, listName)
}