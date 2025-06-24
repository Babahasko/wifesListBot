package repository

type ShoppingRepository interface {
    GetUserState(userID int64) (*UserState, error)
    SetUserState(userID int64, state *UserState) error

    AddShoppingList(userID int64, listName string) error
    GetUserLists(userID int64) ([]string, error)
    DeleteList(userID int64, listName string) error

    AddItemToShoppingList(userID int64, listName, itemName string) error
    GetListItems(userID int64, listName string) ([]*ShoppingItem, error)
    MarkItem(userID int64, listName, itemName string) error
    DeleteMarkedItems(userID int64, listName string) error
}

