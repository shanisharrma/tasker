package todo

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	todoDomain "github.com/shanisharrma/tasker/internal/domain/todo"
)

type GetTodoByID struct {
	todoRepo todoDomain.Repository
	logger   *zerolog.Logger
}

func NewGetTodoByID(todoRepo todoDomain.Repository, logger *zerolog.Logger) *GetTodoByID {
	return &GetTodoByID{
		todoRepo: todoRepo,
		logger:   logger,
	}
}

func (uc *GetTodoByID) Execute(ctx context.Context, userID string, todoID uuid.UUID) (*todoDomain.PopulatedTodo, error) {
	todoItem, err := uc.todoRepo.GetByID(ctx, userID, todoID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to fetch todo by ID")
		return nil, err
	}

	return todoItem, nil
}
