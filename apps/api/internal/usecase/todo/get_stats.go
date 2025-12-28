package todo

import (
	"context"

	"github.com/rs/zerolog"
	todoDomain "github.com/shanisharrma/tasker/internal/domain/todo"
)

type GetTodoStats struct {
	todoRepo todoDomain.Repository
	logger   *zerolog.Logger
}

func NewGetTodoStats(todoRepo todoDomain.Repository, logger *zerolog.Logger) *GetTodoStats {
	return &GetTodoStats{
		todoRepo: todoRepo,
		logger:   logger,
	}
}

func (uc *GetTodoStats) Execute(ctx context.Context, userID string) (*todoDomain.TodoStats, error) {

	stats, err := uc.todoRepo.GetStats(ctx, userID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to fetch todo statistics")
		return nil, err
	}

	return stats, nil
}
