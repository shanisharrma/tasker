package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/app/http/handler"
	"github.com/shanisharrma/tasker/internal/app/http/router"
	"github.com/shanisharrma/tasker/internal/app/modules/attachment"
	"github.com/shanisharrma/tasker/internal/app/modules/auth"
	"github.com/shanisharrma/tasker/internal/app/modules/category"
	"github.com/shanisharrma/tasker/internal/app/modules/comment"
	"github.com/shanisharrma/tasker/internal/app/modules/todo"
	clerkInfra "github.com/shanisharrma/tasker/internal/infra/auth/clerk"
	"github.com/shanisharrma/tasker/internal/infra/database"
	"github.com/shanisharrma/tasker/internal/infra/database/postgres"
	loggerPkg "github.com/shanisharrma/tasker/internal/infra/logger"
	"github.com/shanisharrma/tasker/internal/infra/storage/s3"
	"github.com/shanisharrma/tasker/internal/server"
	"github.com/shanisharrma/tasker/internal/shared/aws"
	"github.com/shanisharrma/tasker/internal/shared/config"
)

type Application struct {
	Server *server.Server
	Logger zerolog.Logger
}

func Bootstrap(ctx context.Context) (*Application, error) {
	// 1. Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %s", err.Error())
	}

	// 2. Logger / Observability
	loggerService := loggerPkg.NewLoggerService(cfg.Observability)
	// (optional but recommended)
	defer loggerService.Shutdown()

	log := loggerPkg.NewLoggerWithService(cfg.Observability, loggerService)
	// 3. Database migration
	if cfg.Primary.Env != "local" {
		if err := database.Migrate(ctx, &log, cfg); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate database!")
			return nil, err
		}
	}

	// 4. Database
	db, err := database.New(cfg, &log, loggerService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	aws, err := aws.NewAWS(cfg)

	// 5. Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
	})

	if loggerService.GetApplication() != nil {
		redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(pingCtx).Err(); err != nil {
		log.Warn().Err(err).Msg("redis unavailable, continuing without redis")
	}

	// 7. Server (runtime container)
	srv := server.New(
		cfg,
		&log,
		loggerService,
		db,
		aws,
		redisClient,
	)

	clerkProvider := clerkInfra.NewIdentityProvider(
		cfg.Auth.SecretKey,
	)

	authModule := auth.NewModule(clerkProvider)

	// 8. Repositories
	todoRepo := postgres.NewTodoRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)
	attachmentRepo := postgres.NewAttachmentRepository(db)
	commentRepo := postgres.NewCommentRepository(db)
	s3Storage := s3.NewObjectStorageRepository(aws.S3)

	// 9. Modules
	todoModule := todo.NewModule(todoRepo, categoryRepo, &log)
	attachmentModule := attachment.NewModule(todoRepo, attachmentRepo, s3Storage, &log, cfg)
	categoryModule := category.NewModule(categoryRepo, &log)
	commentModule := comment.NewModule(todoRepo, commentRepo, &log)

	// 10. Handlers
	handlers := handler.NewHandlers(srv, todoModule, attachmentModule, categoryModule, commentModule, authModule)

	// 11. Router
	r := router.NewRouter(srv, handlers)

	// 12. HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	srv.SetupHTTPServer(httpServer)

	return &Application{
		Server: srv,
		Logger: log,
	}, nil
}
