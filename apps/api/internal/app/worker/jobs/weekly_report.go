package jobs

import (
	"context"

	"github.com/shanisharrma/tasker/internal/usecase/todo"
)

// This Job adapter satisfies the cron.Job interface in app/worker/cron
// It delegates to a usecase and remains infra-agnostic.
type WeeklyReportsJob struct {
	uc *todo.WeeklyReports
}

func NewWeeklyReportsJob(uc *todo.WeeklyReports) *WeeklyReportsJob {
	return &WeeklyReportsJob{uc: uc}
}

func (j *WeeklyReportsJob) Name() string {
	return "weekly-reports"
}

func (j *WeeklyReportsJob) Description() string {
	return "Enqueue weekly productivity reports"
}

func (j *WeeklyReportsJob) Run(ctx context.Context, jobCtx interface{}) error {
	return j.uc.Execute(ctx)
}
