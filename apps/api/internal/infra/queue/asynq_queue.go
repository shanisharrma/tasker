package queue

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	domainJob "github.com/shanisharrma/tasker/internal/domain/job"
)

type AsynqQueue struct {
	client *asynq.Client
	logger *zerolog.Logger
}

func NewAsynqQueue(client *asynq.Client, logger *zerolog.Logger) *AsynqQueue {
	return &AsynqQueue{
		client: client,
		logger: logger,
	}
}

func (q *AsynqQueue) EnqueueReminderEmail(ctx context.Context, p domainJob.ReminderEmailPayload) error {
	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}

	t := asynq.NewTask("email:reminder", payload,
		asynq.MaxRetry(3),
		asynq.Queue("default"),
		asynq.Timeout(30*1e9), // 30s
	)

	if _, err := q.client.EnqueueContext(ctx, t); err != nil {
		q.logger.Error().Err(err).Msg("enqueue reminder failed")
		return err
	}
	return nil
}

func (q *AsynqQueue) EnqueueWeeklyReportEmail(ctx context.Context, p domainJob.WeeklyReportPayload) error {
	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}

	t := asynq.NewTask("email:weekly_report", payload,
		asynq.MaxRetry(3),
		asynq.Queue("default"),
		asynq.Timeout(60*1e9),
	)

	if _, err := q.client.EnqueueContext(ctx, t); err != nil {
		q.logger.Error().Err(err).Msg("enqueue weekly report failed")
		return err
	}
	return nil
}
