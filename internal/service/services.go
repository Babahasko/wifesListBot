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