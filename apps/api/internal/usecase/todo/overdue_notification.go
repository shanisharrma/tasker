package todo

import (
	"context"

	"github.com/rs/zerolog"
	domainJob "github.com/shanisharrma/tasker/internal/domain/job"
	domainTodo "github.com/shanisharrma/tasker/internal/domain/todo"
)

type OverdueNotifications struct {
	repo       domainTodo.Repository
	queue      domainJob.Queue
	logger     *zerolog.Logger
	batchSize  int
	maxPerUser int
}

func NewOverdueNotifications(
	repo domainTodo.Repository,
	queue domainJob.Queue,
	logger *zerolog.Logger,
	batchSize, maxPerUser int) *OverdueNotifications {
	return &OverdueNotifications{
		repo:       repo,
		queue:      queue,
		logger:     logger,
		batchSize:  batchSize,
		maxPerUser: maxPerUser,
	}
}

func (u *OverdueNotifications) Execute(ctx context.Context) error {
	todos, err := u.repo.GetOverdueTodos(ctx, u.batchSize)
	if err != nil {
		return err
	}

	u.logger.Info().
		Int("todo_count", len(todos)).
		Msg("Found overdue todos")

	userCount := map[string]int{}

	for _, t := range todos {
		if userCount[t.UserID] >= u.maxPerUser {
			continue
		}

		p := domainJob.ReminderEmailPayload{
			UserID:    t.UserID,
			TodoID:    t.ID.String(),
			TodoTitle: t.Title,
			DueDate:   *t.DueDate,
			TaskType:  "overdue_notification",
		}

		err := u.queue.EnqueueReminderEmail(ctx, p)
		if err != nil {
			u.logger.Error().
				Err(err).
				Str("todo_id", t.ID.String()).
				Str("user_id", t.UserID).
				Msg("Failed to enqueue reminder email")
			continue
		}

		userCount[t.UserID]++
		u.logger.Info().
			Str("todo_id", t.ID.String()).
			Str("todo_title", t.Title).
			Str("user_id", t.UserID).
			Msg("Enqueued reminder for todo")
	}
	return nil
}
