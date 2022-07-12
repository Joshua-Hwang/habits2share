package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

// Apparently Email is guaranteed in OpenID
// https://openid.net/specs/openid-connect-core-1_0.html
type OpenIdClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	jwt.StandardClaims
}

type TokenParser interface {
	ParseToken(ctx context.Context, token string) (*OpenIdClaims, error)
}
