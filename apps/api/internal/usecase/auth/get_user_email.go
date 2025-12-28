package auth

import (
	"context"

	domain "github.com/shanisharrma/tasker/internal/domain/auth"
)

type GetUserEmail struct {
	provider domain.IdentityProvider
}

func NewGetUserEmail(provider domain.IdentityProvider) *GetUserEmail {
	return &GetUserEmail{provider: provider}
}

func (uc *GetUserEmail) Execute(
	ctx context.Context,
	userID string,
) (string, error) {
	return uc.provider.GetUserEmail(ctx, userID)
}
