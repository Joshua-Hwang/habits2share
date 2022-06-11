module internal/auth_http

go 1.18

replace internal/habit_share => ../habit_share

require (
	github.com/golang-jwt/jwt/v4 v4.3.0
	github.com/google/uuid v1.3.0
	internal/habit_share v0.0.0-00010101000000-000000000000
)
