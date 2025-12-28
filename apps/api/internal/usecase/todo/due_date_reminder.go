package todo

import (
	"context"

	"github.com/rs/zerolog"
	domainJob "github.com/shanisharrma/tasker/internal/domain/job"
	domainTodo "github.com/shanisharrma/tasker/internal/domain/todo"
)

type DueDateReminder struct {
	repo       domainTodo.Repository
	queue      domainJob.Queue
	logger     *zerolog.Logger
	hours      int
	batchSize  int
	maxPerUser int
}

func NewDueDateReminder(
	repo domainTodo.Repository,
	queue domainJob.Queue,
	logger *zerolog.Logger,
	hours int,
	batchSize int,
	maxPerUser int,
) *DueDateReminder {
	return &DueDateReminder{
		repo:       repo,
		queue:      queue,
		logger:     logger,
		hours:      hours,
		batchSize:  batchSize,
		maxPerUser: maxPerUser,
	}
}

func (u *DueDateReminder) Execute(ctx context.Context) error {
	todos, err := u.repo.GetTodosDueInHours(ctx, u.hours, u.batchSize)
	if err != nil {
		return err
	}

	u.logger.Info().
		Int("todo_count", len(todos)).
		Int("hours", u.hours).
		Msg("Found todos due soon")

		// keep track per user to limit notifications
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
			TaskType:  "due_date_reminder",
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
