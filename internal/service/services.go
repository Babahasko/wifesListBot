package service

import (
	"shopping_bot/internal/repository"
)

type ShoppingService struct {
	Repo repository.ShoppingRepository
}

func NewShopService(repo repository.ShoppingRepository) *ShoppingService {
	return &ShoppingService{
		Repo: repo,
	}
}

func (s *ShoppingService) SetCurrentList(userID int64, listName string) error {
	state, _ := s.Repo.GetUserState(userID)
	if state == nil {
		state = &repository.UserState{}
	}
	state.CurrentList = listName
	err :=s.Repo.SetUserState(userID, state)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShoppingService) GetCurrentList(userID int64) (string, error) {
	state, err := s.Repo.GetUserState(userID)
	if err != nil {
		return "", err
	}
	return state.CurrentList, nil
}

func (s *ShoppingService) AddShoppingList(userID int64, listName string) error {
	return s.Repo.AddShoppingList(userID, listName)
}

func (s *ShoppingService) GetUserLists(userID int64) ([]string, error) {
	userLists, err := s.Repo.GetUserLists(userID)
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
	return s.Repo.DeleteList(userID, listName)
}

func (s *ShoppingService) AddItemToShoppingList(userID int64, listName, itemName string) error {
	return s.Repo.AddItemToShoppingList(userID, listName, itemName)
}

func (s *ShoppingService) GetListItems(userID int64, listName string) ([]string, error) {
	listItems, err := s.Repo.GetListItems(userID, listName)
	if err != nil {
		return nil, err
	}
	strListItems := make([]string, 0, len(listItems))
	for _, item := range listItems {
		strListItems = append(strListItems, item.Name)
	}
	return strListItems, nil
}

func (s *ShoppingService) MarkItem(userID int64, listName, itemName string) error {
	return s.Repo.MarkItem(userID, listName, itemName)
}

func (s *ShoppingService) DeleteMarkedItems(userID int64, listName string) error {
	return s.Repo.DeleteMarkedItems(userID, listName)
}