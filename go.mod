module joshuahwang/habits2share

go 1.18

replace internal/auth => ./internal/auth

replace internal/auth_file => ./internal/auth_file

replace internal/habit_share => ./internal/habit_share

replace internal/habit_share_file => ./internal/habit_share_file

require (
	github.com/google/uuid v1.3.0
	internal/auth v1.0.0
	internal/auth_file v0.0.0-00010101000000-000000000000
	internal/habit_share v0.0.0-00010101000000-000000000000
	internal/habit_share_file v0.0.0-00010101000000-000000000000
)

require github.com/golang-jwt/jwt/v4 v4.3.0 // indirect
