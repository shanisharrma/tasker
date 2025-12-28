package auth

import "context"

type IdentityProvider interface {
	GetUserEmail(ctx context.Context, userID string) (string, error)
}
