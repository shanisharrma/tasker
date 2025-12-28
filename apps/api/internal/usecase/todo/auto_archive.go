package todo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	domainTodo "github.com/shanisharrma/tasker/internal/domain/todo"
)

type AutoArchive struct {
	repo      domainTodo.Repository
	logger    *zerolog.Logger
	threshold int // days
	batchSize int
}

func NewAutoArchive(
	repo domainTodo.Repository,
	logger *zerolog.Logger,
	threshold, batchSize int,
) *AutoArchive {
	return &AutoArchive{
		repo:      repo,
		logger:    logger,
		threshold: threshold,
		batchSize: batchSize,
	}
}

func (u *AutoArchive) Execute(ctx context.Context) error {
	cutoff := time.Now().AddDate(0, 0, -u.threshold)

	u.logger.Info().
		Time("cutoff_date", cutoff).
		Msg("Searching for completed todos to archive")

	todos, err := u.repo.GetCompletedTodosOlderThan(ctx, cutoff, u.batchSize)
	if err != nil {
		return err
	}
	u.logger.Info().
		Int("todo_count", len(todos)).
		Msg("Found completed todos to archive")

	if len(todos) == 0 {
		u.logger.Info().Msg("No todos to archive")
		return nil
	}

	ids := make([]string, 0, len(todos))
	for _, t := range todos {
		ids = append(ids, t.ID.String())
	}

	// Convert back to uuid.UUID inside repo implementation
	// or change repo interface if desired
	return u.repo.ArchiveTodos(ctx, todosIDsToUUID(ids))
}

// helper that should be implemented/adjusted to your types
func todosIDsToUUID(ids []string) []uuid.UUID {
	// implement conversion per your codebase; placeholder here
	var out []uuid.UUID
	return out
}
