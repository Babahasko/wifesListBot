package bot

import "errors"

var (
	ErrListExists   = errors.New("shop list already exists")
	ErrNoLists      = errors.New("no shop lists for user")
	ErrListNotFound = errors.New("shop list not found")
)
