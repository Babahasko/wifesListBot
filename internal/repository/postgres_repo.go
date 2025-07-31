package repository

import (
	"errors"
	"shopping_bot/internal/models"

	"gorm.io/gorm"
)

type PostgresShoppingRepository struct {
	db *gorm.DB
}

// GORM-модели для записи в БД
type GormUserState struct {
	gorm.Model
	UserID      int64 `gorm:"uniqueIndex"`
	CurrentList string
}

type GormShoppingList struct {
	gorm.Model
	UserID int64 `gorm:"index"`
	Name   string
	Items  []*GormShoppingItem `gorm:"foreignKey:ShoppingListID"`
}

type GormShoppingItem struct {
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
	var gormUserState GormUserState
	result := r.db.First(&gormUserState, "user_id = ?", userID)
	if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, ErrNoState
        }
        return nil, result.Error
	}
	return r.convertUserStateToModel(&gormUserState), nil
}

func (r *PostgresShoppingRepository) SetUserState(userID int64, state *models.UserState) error {
    gormUserState := r.convertUserStateToGormModel(userID, state)
	result := r.db.Create(&gormUserState)
	if result.Error != nil {
		return result.Error
	}
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


//function for converting from gorm model to standart model
func (r *PostgresShoppingRepository) convertUserStateToModel(gormUS *GormUserState) *models.UserState {
    return &models.UserState{
        CurrentList: gormUS.CurrentList,
    }
}
func (r *PostgresShoppingRepository) convertUserStateToGormModel(userid int64, userState *models.UserState) *GormUserState {
    return &GormUserState{
        UserID: userid,
		CurrentList: userState.CurrentList,
    }
}