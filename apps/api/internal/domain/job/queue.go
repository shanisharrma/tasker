package job

import (
	"context"
	"time"
)

// Queue is the domain-level interface for enqueuing background tasks.
// Implementations live under internal/infra/queue.
type Queue interface {
	EnqueueReminderEmail(ctx context.Context, p ReminderEmailPayload) error
	EnqueueWeeklyReportEmail(ctx context.Context, p WeeklyReportPayload) error
}

// Payloads used by domain/usecase layers
type ReminderEmailPayload struct {
	UserID    string
	TodoID    string
	TodoTitle string
	DueDate   time.Time
	TaskType  string // e.g. "due_date_reminder" or "overdue_notification"
}

type WeeklyReportPayload struct {
	UserID         string
	WeekStart      time.Time
	WeekEnd        time.Time
	CompletedCount int
	ActiveCount    int
	OverdueCount   int
	CompletedTodos interface{} // keep domain types or concrete types as needed
	OverdueTodos   interface{}
}
