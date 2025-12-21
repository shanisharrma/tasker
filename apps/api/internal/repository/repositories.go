package repository

import "github.com/shanisharrma/tasker/internal/server"

type Repositories struct {
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}
