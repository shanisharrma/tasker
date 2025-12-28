package todo

import (
	"context"

	"github.com/rs/zerolog"
	category "github.com/shanisharrma/tasker/internal/domain/category"
	"github.com/shanisharrma/tasker/internal/domain/todo"
	todoDomain "github.com/shanisharrma/tasker/internal/domain/todo"
	"github.com/shanisharrma/tasker/internal/shared/errs"
)

type UpdateTodo struct {
	todoRepo     todoDomain.Repository
	categoryRepo category.Repository
	logger       *zerolog.Logger
}

func NewUpdateTodo(todoRepo todoDomain.Repository, categoryRepo category.Repository, logger *zerolog.Logger) *UpdateTodo {
	return &UpdateTodo{
		todoRepo:     todoRepo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *UpdateTodo) Execute(ctx context.Context, userID string, payload *todo.UpdateTodoPayload) (*todo.Todo, error) {

	// Validate parent todo exists and belongs to user (if provided)
	if payload.ParentTodoID != nil {
		parentTodo, err := uc.todoRepo.CheckExists(ctx, userID, *payload.ParentTodoID)
		if err != nil {
			uc.logger.Error().Err(err).Msg("parent todo validation failed")
			return nil, err
		}

		if parentTodo.ID == payload.ID {
			err := errs.NewBadRequestError("Todo cannot be its own parent", false, nil, nil, nil)
			uc.logger.Warn().Msg("todo cannot be its own parent")
			return nil, err
		}

		if !parentTodo.CanHaveChildren() {
			err := errs.NewBadRequestError("Parent todo cannot have children (subtasks can't have subtasks)", false, nil, nil, nil)
			uc.logger.Warn().Msg("parent todo cannot have children")
			return nil, err
		}

		uc.logger.Debug().Msg("parent todo validation passed")
	}

	// Validate category exists and belongs to user (if provided)
	if payload.CategoryID != nil {
		_, err := uc.categoryRepo.GetCategoryByID(ctx, userID, *payload.CategoryID)
		if err != nil {
			uc.logger.Error().Err(err).Msg("category validation failed")
			return nil, err
		}

		uc.logger.Debug().Msg("category validation passed")
	}

	updatedTodo, err := uc.todoRepo.Update(ctx, userID, payload)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to update todo")
		return nil, err
	}

	// Business event log
	uc.logger.Info().
		Str("event", "todo_updated").
		Str("todo_id", updatedTodo.ID.String()).
		Str("title", updatedTodo.Title).
		Str("category_id", func() string {
			if updatedTodo.CategoryID != nil {
				return updatedTodo.CategoryID.String()
			}
			return ""
		}()).
		Str("priority", string(updatedTodo.Priority)).
		Str("status", string(updatedTodo.Status)).
		Msg("Todo updated successfully")

	return updatedTodo, nil
}
