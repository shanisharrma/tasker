package todo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shanisharrma/tasker/internal/shared/types"
)

type Repository interface {
	Create(ctx context.Context, userID string, payload *CreateTodoPayload) (*Todo, error)
	GetByID(ctx context.Context, userID string, todoID uuid.UUID) (*PopulatedTodo, error)
	CheckExists(ctx context.Context, userID string, todoID uuid.UUID) (*Todo, error)
	List(ctx context.Context, userID string, query *GetTodosQuery) (*types.PaginatedResponse[PopulatedTodo], error)
	Update(ctx context.Context, userID string, payload *UpdateTodoPayload) (*Todo, error)
	Delete(ctx context.Context, userID string, todoID uuid.UUID) error
	GetStats(ctx context.Context, userID string) (*TodoStats, error)

	// cron
	GetTodosDueInHours(ctx context.Context, hours, limit int) ([]Todo, error)
	GetOverdueTodos(ctx context.Context, limit int) ([]Todo, error)
	GetCompletedTodosOlderThan(ctx context.Context, cutoffDate time.Time, limit int) ([]Todo, error)
	ArchiveTodos(ctx context.Context, todoIDs []uuid.UUID) error
	GetWeeklyStatsForUsers(ctx context.Context, startDate, endDate time.Time) ([]UserWeeklyStats, error)
	GetCompletedTodosForUser(ctx context.Context, userID string, startDate, endDate time.Time) ([]PopulatedTodo, error)
	GetOverdueTodosForUser(ctx context.Context, userID string) ([]PopulatedTodo, error)
}
