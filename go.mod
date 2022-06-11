module joshuahwang/habits2share

go 1.18

replace internal/auth => ./internal/auth

replace internal/auth_file => ./internal/auth_file

replace internal/habit_share => ./internal/habit_share

replace internal/habit_share_file => ./internal/habit_share_file

require (
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/google/uuid v1.3.0
	internal/auth v1.0.0
	internal/auth_file v0.0.0-00010101000000-000000000000
	internal/habit_share v0.0.0-00010101000000-000000000000
	internal/habit_share_file v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang-jwt/jwt/v4 v4.3.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/lib/pq v1.10.4 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	golang.org/x/tools v0.1.9 // indirect
	google.golang.org/genproto v0.0.0-20220218161850-94dd64e39d7c // indirect
	google.golang.org/grpc v1.44.0 // indirect
)
