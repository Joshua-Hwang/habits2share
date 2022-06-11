package auth

import (
	"internal/habit_share"
)

// AuthService doesn't do anything a couple of hooks could
// This could be extended in future and so is provided as a service it's also
// orthogonal to other pieces so the over-engineering doesn't affect anything
type AuthService struct {
	GetUserId func() (string, error)
}

var _ habit_share.AuthInterface = (*AuthService)(nil)

// GetCurrentUser implements habit_share.AuthInterface
func (a AuthService) GetCurrentUser() (string, error) {
	userId, err := a.GetUserId()
	if err != nil {
		return "", err
	}
	if userId == "" {
		return "", habit_share.UserNotFoundError
	}
	return userId, nil
}
