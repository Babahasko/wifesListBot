package repository

import (
	"shopping_bot/internal/models"

	"gorm.io/gorm"
)

type PostgresShoppingRepository struct {
	db *gorm.DB
}

// GORM-модели
type gormUserState struct {
    gorm.Model
    UserID      int64 `gorm:"uniqueIndex"`
    CurrentList string
}

type gormShoppingList struct {
    gorm.Model
    UserID int64 `gorm:"index"`
    Name   string
    Items  []*gormShoppingItem `gorm:"foreignKey:ShoppingListID"`
}

type gormShoppingItem struct {
    gorm.Model
    ShoppingListID uint
    ListName       string
    Name           string
    Checked        bool
}

func NewPostgresShoppingRepository(db *gorm.DB) *PostgresShoppingRepository {
    return &PostgresShoppingRepository{db: db}
}


func (r *PostgresShoppingRepository) GetUserState(userID int64) (*models.UserState, error) {

 return nil, nil
}

func (r *PostgresShoppingRepository) SetUserState(userID int64, state *models.UserState) error {
	return nil
}

func (r *PostgresShoppingRepository) AddShoppingList(userID int64, listName string) error {
	return nil
}
func (r *PostgresShoppingRepository) GetUserLists(userID int64) (map[string]*models.ShoppingList, error) {
	return nil, nil
}
func (r *PostgresShoppingRepository) DeleteList(userID int64, listName string) error {
	return nil
}

func (r *PostgresShoppingRepository) AddItemToShoppingList(userID int64, listName, itemName string) error {
	return nil
}
func (r *PostgresShoppingRepository) GetListItems(userID int64, listName string) ([]*models.ShoppingItem, error) {
	return nil, nil
}
func (r *PostgresShoppingRepository) MarkItem(userID int64, listName, itemName string) error {
	return nil
}
func (r *PostgresShoppingRepository) DeleteMarkedItems(userID int64, listName string) error {
	return nil
}