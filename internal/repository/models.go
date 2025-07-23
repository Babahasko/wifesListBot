package repository

import "gorm.io/gorm"

type UserState struct {
	gorm.Model
	CurrentList string `json:"state"` //текущий лист покупок редактируемый пользователем
}

type ShoppingList struct {
	gorm.Model
	Name  string          `json:"name"`
	Items []*ShoppingItem `json:"items" gorm:"foreignKey:ShoppingListID"`
}

type ShoppingItem struct {
	gorm.Model
	ShoppingListID uint `json:"shopping_list_id" gorm:"index"`
	ListName string `json:"list_name"`
	Name     string `json:"item_name"`
	Checked  bool   `json:"checked"`
}