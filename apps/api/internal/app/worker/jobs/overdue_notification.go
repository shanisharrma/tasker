package jobs

import (
	"context"

	"github.com/shanisharrma/tasker/internal/usecase/todo"
)

// This Job adapter satisfies the cron.Job interface in app/worker/cron
// It delegates to a usecase and remains infra-agnostic.
type OverdueNotificationsJob struct {
	uc *todo.OverdueNotifications
}

func NewOverdueNotificationsJob(uc *todo.OverdueNotifications) *OverdueNotificationsJob {
	return &OverdueNotificationsJob{uc: uc}
}

func (j *OverdueNotificationsJob) Name() string {
	return "overdue-notifications"
}

func (j *OverdueNotificationsJob) Description() string {
	return "Enqueue notifications for overdue todos"
}

func (j *OverdueNotificationsJob) Run(ctx context.Context, jobCtx interface{}) error {
	return j.uc.Execute(ctx)
}
