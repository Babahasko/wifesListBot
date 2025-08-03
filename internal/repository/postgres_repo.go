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
	UserID int64               `gorm:"index"`
	Name   string              `gorm:"not null"`
	Items  []*GormShoppingItem `gorm:"foreignKey:ShoppingListID"`
}

type GormShoppingItem struct {
	gorm.Model
	ShoppingListID uint   `gorm:"not null"`
	ListName       string `gorm:"not null"`
	Name           string `gorm:"not null"`
	Checked        bool   `gorm:"default:false"`
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
	// Проверяем, что список существует
	var existingList GormShoppingList
	result := r.db.Where("user_id = ? AND name = ?", userID, listName).First(&existingList)

	if result.Error == nil {
		return ErrListExists
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	newList := GormShoppingList{
		UserID: userID,
		Name:   listName,
		Items:  make([]*GormShoppingItem, 0),
	}

	result = r.db.Create(&newList)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostgresShoppingRepository) GetUserLists(userID int64) (map[string]*models.ShoppingList, error) {
	var gormLists []GormShoppingList
	result := r.db.Where("user_id = ?", userID).Preload("Items").Find(&gormLists)
	if result.Error != nil {
		return nil, result.Error
	}
	// Преобразуем в map[string]*models.ShoppingList
	userLists := make(map[string]*models.ShoppingList)

	for _, gormList := range gormLists {
		shoppingList := r.convertGormListToModel(&gormList)
		userLists[gormList.Name] = shoppingList
	}
	return userLists, nil
}
func (r *PostgresShoppingRepository) DeleteList(userID int64, listName string) error {
	result := r.db.Where("user_id = ? AND name = ?", userID, listName).Delete(&GormShoppingList{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrListExists
	}
	
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

// function for converting from gorm model to standart model
func (r *PostgresShoppingRepository) convertUserStateToModel(gormUS *GormUserState) *models.UserState {
	return &models.UserState{
		CurrentList: gormUS.CurrentList,
	}
}
func (r *PostgresShoppingRepository) convertUserStateToGormModel(userid int64, userState *models.UserState) *GormUserState {
	return &GormUserState{
		UserID:      userid,
		CurrentList: userState.CurrentList,
	}
}

func (r *PostgresShoppingRepository) convertGormListToModel(gormList *GormShoppingList) *models.ShoppingList {
	items := make([]*models.ShoppingItem, len(gormList.Items))
	for i, gormItem := range gormList.Items {
		items[i] = &models.ShoppingItem{
			ShoppingListID: gormItem.ShoppingListID,
			ListName: gormItem.ListName,
			Name: gormItem.Name,
			Checked: gormItem.Checked,
		}
	}
	return &models.ShoppingList{
		Name: gormList.Name,
		Items: items,
	}
}
