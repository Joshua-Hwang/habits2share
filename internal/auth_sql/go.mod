module internal/auth_sql

go 1.18

replace internal/auth => ../auth

replace internal/habit_share => ../habit_share

require internal/auth v0.0.0-00010101000000-000000000000

require (
	github.com/golang-jwt/jwt/v4 v4.3.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	internal/habit_share v0.0.0-00010101000000-000000000000 // indirect
)
