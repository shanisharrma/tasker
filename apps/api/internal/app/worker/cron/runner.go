package cron

import (
	"context"
	"fmt"
)

// Runner executes a single job by name using a prepared registry and container
// It does NOT create infra; everything must already be wired in bootstrap.
type Runner struct {
	registry *JobRegistry
}

func NewRunner(registry *JobRegistry) *Runner {
	return &Runner{registry: registry}
}

func (r *Runner) Run(jobName string) error {
	job, err := r.registry.Get(jobName)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := job.Run(ctx, nil); err != nil {
		return fmt.Errorf("job %s failed: %w", jobName, err)
	}

	return nil
}
