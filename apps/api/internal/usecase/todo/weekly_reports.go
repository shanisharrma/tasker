package todo

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	domainJob "github.com/shanisharrma/tasker/internal/domain/job"
	domainTodo "github.com/shanisharrma/tasker/internal/domain/todo"
)

type WeeklyReports struct {
	repo      domainTodo.Repository
	queue     domainJob.Queue
	logger    *zerolog.Logger
	batchSize int
}

func NewWeeklyReports(
	repo domainTodo.Repository,
	queue domainJob.Queue,
	logger *zerolog.Logger,
	batchSize int) *WeeklyReports {
	return &WeeklyReports{
		repo:      repo,
		queue:     queue,
		logger:    logger,
		batchSize: batchSize,
	}
}

func (u *WeeklyReports) Execute(ctx context.Context) error {
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)

	stats, err := u.repo.GetWeeklyStatsForUsers(ctx, weekAgo, now)
	if err != nil {
		return err
	}

	u.logger.Info().
		Int("user_count", len(stats)).
		Msg("Generating weekly reports")

	for _, s := range stats {
		completed, err := u.repo.GetCompletedTodosForUser(ctx, s.UserID, weekAgo, now)
		if err != nil {
			completed = []domainTodo.PopulatedTodo{}
		}

		overdue, err := u.repo.GetOverdueTodosForUser(ctx, s.UserID)
		if err != nil {
			u.logger.Error().
				Err(err).
				Str("user_id", s.UserID).
				Msg("Failed to fetch overdue todos")
			overdue = []domainTodo.PopulatedTodo{}
		}

		p := domainJob.WeeklyReportPayload{
			UserID:         s.UserID,
			WeekStart:      weekAgo,
			WeekEnd:        now,
			CompletedCount: s.CompletedCount,
			ActiveCount:    s.ActiveCount,
			OverdueCount:   s.OverdueCount,
			CompletedTodos: completed,
			OverdueTodos:   overdue,
		}

		if err := u.queue.EnqueueWeeklyReportEmail(ctx, p); err != nil {
			// log and continue
			u.logger.Error().
				Err(err).
				Str("user_id", s.UserID).
				Msg("Failed to enqueue weekly report")
			continue
		}

		u.logger.Info().
			Str("user_id", s.UserID).
			Int("created", s.CreatedCount).
			Int("completed", s.CompletedCount).
			Int("active", s.ActiveCount).
			Int("overdue", s.OverdueCount).
			Msg("Enqueued weekly report")

	}
	return nil
}
