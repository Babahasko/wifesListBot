package shop

import "time"

// Структура элемента списка покупок
type Purchase struct {
    Name string
	Category string
	Price float64
	Priority int
	DateAdded time.Time
	IsCompleted bool
}