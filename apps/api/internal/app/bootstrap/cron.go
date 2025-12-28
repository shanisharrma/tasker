package bootstrap

import (
	"fmt"

	"github.com/hibiken/asynq"

	jobsAdapters "github.com/shanisharrma/tasker/internal/app/worker/jobs"
	"github.com/shanisharrma/tasker/internal/infra/database"
	"github.com/shanisharrma/tasker/internal/infra/database/postgres"
	infraLogger "github.com/shanisharrma/tasker/internal/infra/logger"
	"github.com/shanisharrma/tasker/internal/infra/queue"
	"github.com/shanisharrma/tasker/internal/shared/config"
	todoUC "github.com/shanisharrma/tasker/internal/usecase/todo"
)

// WorkerContainer holds constructed dependencies for worker processes
type WorkerContainer struct {
	Config         *config.Config
	JobClient      *asynq.Client
	DueDateJob     *jobsAdapters.DueDateRemindersJob
	OverdueJob     *jobsAdapters.OverdueNotificationsJob
	WeeklyJob      *jobsAdapters.WeeklyReportsJob
	AutoArchiveJob *jobsAdapters.AutoArchiveJob
}

func BuildWorkerContainer() (*WorkerContainer, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// logger
	loggerService := infraLogger.NewLoggerService(cfg.Observability)
	loggerInstance := infraLogger.NewLoggerWithService(cfg.Observability, loggerService)

	// db
	db, err := database.New(cfg, &loggerInstance, loggerService)
	if err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	// construct repositories (infra implementation)
	todoRepo := postgres.NewTodoRepository(db)

	// asynq client
	jobClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.Redis.Address, Password: cfg.Redis.Password})
	queueAdapter := queue.NewAsynqQueue(jobClient, &loggerInstance)

	// build usecases
	dueUC := todoUC.NewDueDateReminder(todoRepo, queueAdapter, &loggerInstance, cfg.Cron.ReminderHours, cfg.Cron.BatchSize, cfg.Cron.MaxTodosPerUserNotification)
	overdueUC := todoUC.NewOverdueNotifications(todoRepo, queueAdapter, &loggerInstance, cfg.Cron.BatchSize, cfg.Cron.MaxTodosPerUserNotification)
	weeklyUC := todoUC.NewWeeklyReports(todoRepo, queueAdapter, &loggerInstance, cfg.Cron.BatchSize)
	autoArchiveUC := todoUC.NewAutoArchive(todoRepo, &loggerInstance, cfg.Cron.ArchiveDaysThreshold, cfg.Cron.BatchSize)

	// job adapters
	dueJob := jobsAdapters.NewDueDateRemindersJob(dueUC)
	overdueJob := jobsAdapters.NewOverdueNotificationsJob(overdueUC)
	weeklyJob := jobsAdapters.NewWeeklyReportsJob(weeklyUC)
	autoArchiveJob := jobsAdapters.NewAutoArchiveJob(autoArchiveUC)

	return &WorkerContainer{
		Config:         cfg,
		JobClient:      jobClient,
		DueDateJob:     dueJob,
		OverdueJob:     overdueJob,
		WeeklyJob:      weeklyJob,
		AutoArchiveJob: autoArchiveJob,
	}, nil
}
