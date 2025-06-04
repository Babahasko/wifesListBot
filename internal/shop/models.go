package shop

import "time"

// Структура элемента списка покупок
type Purchase struct {
    Name string
	Category string
	Price int32
	Priority int16
	DateAdded time.Time
	IsCompleted bool
}