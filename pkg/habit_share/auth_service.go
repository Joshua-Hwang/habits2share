package habit_share

type AuthInterface interface {
	GetCurrentUser() (string, error)
}
