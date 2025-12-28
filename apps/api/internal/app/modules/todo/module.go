package todo

import (
	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/domain/category"
	"github.com/shanisharrma/tasker/internal/domain/todo"
	todoUC "github.com/shanisharrma/tasker/internal/usecase/todo"
)

type Module struct {
	CreateTodoUC  *todoUC.CreateTodo
	GetTodoByIdUC *todoUC.GetTodoByID
	ListTodosUC   *todoUC.ListTodos
	UpdateTodoUC  *todoUC.UpdateTodo
	DeleteTodoUC  *todoUC.DeleteTodo
	GetStatsUC    *todoUC.GetTodoStats
}

func NewModule(todoRepo todo.Repository, categoryRepo category.Repository, logger *zerolog.Logger) *Module {

	return &Module{
		CreateTodoUC:  todoUC.NewCreateTodo(todoRepo, categoryRepo, logger),
		GetTodoByIdUC: todoUC.NewGetTodoByID(todoRepo, logger),
		ListTodosUC:   todoUC.NewListTodos(todoRepo, logger),
		UpdateTodoUC:  todoUC.NewUpdateTodo(todoRepo, categoryRepo, logger),
		DeleteTodoUC:  todoUC.NewDeleteTodo(todoRepo, logger),
		GetStatsUC:    todoUC.NewGetTodoStats(todoRepo, logger),
	}

}
