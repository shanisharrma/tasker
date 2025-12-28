package todo

import (
	"context"

	"github.com/rs/zerolog"
	category "github.com/shanisharrma/tasker/internal/domain/category"
	todoDomain "github.com/shanisharrma/tasker/internal/domain/todo"
	"github.com/shanisharrma/tasker/internal/shared/errs"
)

type CreateTodo struct {
	todoRepo     todoDomain.Repository
	categoryRepo category.Repository
	logger       *zerolog.Logger
}

func NewCreateTodo(todoRepo todoDomain.Repository, categoryRepo category.Repository, logger *zerolog.Logger) *CreateTodo {
	return &CreateTodo{
		todoRepo:     todoRepo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *CreateTodo) Execute(ctx context.Context, userID string, payload *todoDomain.CreateTodoPayload) (*todoDomain.Todo, error) {
	// Validate parent todo exists and belongs to user (if provided)
	if payload.ParentTodoID != nil {
		parentTodo, err := uc.todoRepo.CheckExists(ctx, userID, *payload.ParentTodoID)
		if err != nil {
			uc.logger.Error().Err(err).Msg("parent todo validation failed")
			return nil, err
		}

		if !parentTodo.CanHaveChildren() {
			err := errs.NewBadRequestError("Parent todo cannot have children (subtasks can't have subtasks)", false, nil, nil, nil)
			uc.logger.Warn().Msg("parent todo cannot have children")
			return nil, err
		}
	}

	// Validate category exists and belongs to user (if provided)
	if payload.CategoryID != nil {
		_, err := uc.categoryRepo.GetCategoryByID(ctx, userID, *payload.CategoryID)
		if err != nil {
			uc.logger.Error().Err(err).Msg("category validation failed")
			return nil, err
		}
	}

	todoItem, err := uc.todoRepo.Create(ctx, userID, payload)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to create todo")
		return nil, err
	}

	// Business event log
	uc.logger.Info().
		Str("event", "todo_created").
		Str("todo_id", todoItem.ID.String()).
		Str("title", todoItem.Title).
		Str("category_id", func() string {
			if todoItem.CategoryID != nil {
				return todoItem.CategoryID.String()
			}
			return ""
		}()).
		Str("priority", string(todoItem.Priority)).
		Msg("Todo created successfully")

	return todoItem, nil
}
