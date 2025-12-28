package todo

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	todoDomain "github.com/shanisharrma/tasker/internal/domain/todo"
)

type DeleteTodo struct {
	todoRepo todoDomain.Repository
	logger   *zerolog.Logger
}

func NewDeleteTodo(todoRepo todoDomain.Repository, logger *zerolog.Logger) *DeleteTodo {
	return &DeleteTodo{
		todoRepo: todoRepo,
		logger:   logger,
	}
}

func (uc *DeleteTodo) Execute(ctx context.Context, userID string, todoID uuid.UUID) error {
	err := uc.todoRepo.Delete(ctx, userID, todoID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to delete todo")
		return err
	}

	// Business event log
	uc.logger.Info().
		Str("event", "todo_deleted").
		Str("todo_id", todoID.String()).
		Msg("Todo deleted successfully")

	return nil
}
