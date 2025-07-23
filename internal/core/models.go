package core

type UserState struct {
    CurrentList string
}

type ShoppingList struct {
    Name  string
    Items []*ShoppingItem
}

type ShoppingItem struct {
    ListName string
    Name     string
    Checked  bool
}