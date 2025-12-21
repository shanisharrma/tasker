package service

import "github.com/shanisharrma/tasker/internal/server"

type AuthService struct {
	server *server.Server
}

func NewAuthService(s *server.Server) *AuthService {
	return &AuthService{
		server: s,
	}
}
