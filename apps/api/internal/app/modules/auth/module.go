package auth

import (
	domain "github.com/shanisharrma/tasker/internal/domain/auth"
	uc "github.com/shanisharrma/tasker/internal/usecase/auth"
)

type Module struct {
	GetUserEmail *uc.GetUserEmail
}

func NewModule(provider domain.IdentityProvider) *Module {
	return &Module{
		GetUserEmail: uc.NewGetUserEmail(provider),
	}
}
