package clerk

import (
	"context"
	"fmt"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
)

type IdentityProvider struct{}

func NewIdentityProvider(secretKey string) *IdentityProvider {
	clerk.SetKey(secretKey)
	return &IdentityProvider{}
}

func (p *IdentityProvider) GetUserEmail(
	ctx context.Context,
	userID string,
) (string, error) {

	user, err := clerkUser.Get(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user from Clerk: %w", err)
	}

	if len(user.EmailAddresses) == 0 {
		return "", fmt.Errorf("user %s has no email addresses", userID)
	}

	if user.PrimaryEmailAddressID != nil {
		for _, email := range user.EmailAddresses {
			if email.ID == *user.PrimaryEmailAddressID {
				return email.EmailAddress, nil
			}
		}
	}

	return user.EmailAddresses[0].EmailAddress, nil
}
