package jobs

import (
	"context"

	"github.com/shanisharrma/tasker/internal/usecase/todo"
)

// This Job adapter satisfies the cron.Job interface in app/worker/cron
// It delegates to a usecase and remains infra-agnostic.
type DueDateRemindersJob struct {
	uc *todo.DueDateReminder
}

func NewDueDateRemindersJob(uc *todo.DueDateReminder) *DueDateRemindersJob {
	return &DueDateRemindersJob{uc: uc}
}

func (j *DueDateRemindersJob) Name() string {
	return "due-date-reminders"
}

func (j *DueDateRemindersJob) Description() string {
	return "Enqueue email reminders for todos due soon"
}

func (j *DueDateRemindersJob) Run(ctx context.Context, jobCtx interface{}) error {
	return j.uc.Execute(ctx)
}
