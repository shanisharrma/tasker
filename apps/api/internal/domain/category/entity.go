package category

import "github.com/shanisharrma/tasker/internal/shared/types"

type Category struct {
	types.Base

	UserID      string  `json:"userId" db:"user_id"`
	Name        string  `json:"name" db:"name"`
	Color       string  `json:"color" db:"color"`
	Description *string `json:"description"`
}
