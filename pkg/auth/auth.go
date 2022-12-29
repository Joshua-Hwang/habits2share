package auth

import (
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
)

// AuthService doesn't do anything a couple of hooks could
// This could be extended in future and so is provided as a service it's also
// orthogonal to other pieces so the over-engineering doesn't affect anything
type AuthService struct {
	UserId string
}

var _ habit_share.AuthInterface = (*AuthService)(nil)

// GetCurrentUser implements habit_share.AuthInterface
func (a AuthService) GetCurrentUser() (string, error) {
	userId := a.UserId
	if userId == "" {
		return "", habit_share.UserNotFoundError
	}
	return userId, nil
}
