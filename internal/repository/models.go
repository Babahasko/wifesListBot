package repository

type UserState struct {
	CurrentList string `json:"state"` //текущий лист покупок редактируемый пользователем
}

type ShoppingList struct {
	Name  string          `json:"name"`
	Items []*ShoppingItem `json:"items"`
}

type ShoppingItem struct {
	ListName string `json:"list_name"`
	Name     string `json:"item_name"`
	Checked  bool   `json:"checked"`
}