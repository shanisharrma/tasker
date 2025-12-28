package comment

import (
	"github.com/google/uuid"
	"github.com/shanisharrma/tasker/internal/shared/types"
)

type Comment struct {
	types.Base

	TodoID  uuid.UUID `json:"todoId" db:"todo_id"`
	UserID  string    `json:"userId" db:"user_id"`
	Content string    `json:"content" db:"content"`
}
