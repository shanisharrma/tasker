package todo

import (
	"context"

	"github.com/rs/zerolog"
	todoDomain "github.com/shanisharrma/tasker/internal/domain/todo"
	"github.com/shanisharrma/tasker/internal/shared/types"
)

type ListTodos struct {
	todoRepo todoDomain.Repository
	logger   *zerolog.Logger
}

func NewListTodos(todoRepo todoDomain.Repository, logger *zerolog.Logger) *ListTodos {
	return &ListTodos{
		todoRepo: todoRepo,
		logger:   logger,
	}
}

func (uc *ListTodos) Execute(ctx context.Context, userID string, query *todoDomain.GetTodosQuery) (*types.PaginatedResponse[todoDomain.PopulatedTodo], error) {

	result, err := uc.todoRepo.List(ctx, userID, query)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to fetch todos")
		return nil, err
	}

	return result, nil
}
