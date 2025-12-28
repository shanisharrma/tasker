package jobs

import (
	"context"

	"github.com/shanisharrma/tasker/internal/usecase/todo"
)

// This Job adapter satisfies the cron.Job interface in app/worker/cron
// It delegates to a usecase and remains infra-agnostic.
type AutoArchiveJob struct {
	uc *todo.AutoArchive
}

func NewAutoArchiveJob(uc *todo.AutoArchive) *AutoArchiveJob {
	return &AutoArchiveJob{uc: uc}
}

func (j *AutoArchiveJob) Name() string {
	return "auto-archive"
}

func (j *AutoArchiveJob) Description() string {
	return "Archive old completed todos"
}

func (j *AutoArchiveJob) Run(ctx context.Context, jobCtx interface{}) error {
	return j.uc.Execute(ctx)
}
