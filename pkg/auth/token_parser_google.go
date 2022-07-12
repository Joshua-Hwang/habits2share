package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenParserGoogle struct {
	WebClientId    string
	MobileClientId string
}

var _ TokenParser = (*TokenParserGoogle)(nil)

func (a *TokenParserGoogle) ParseToken(ctx context.Context, token string) (*OpenIdClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&OpenIdClaims{},
		func(token *jwt.Token) (interface{}, error) {
			pem, err := getPublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				return nil, err
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err
			}
			return key, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*OpenIdClaims)
	if !ok {
		return claims, errors.New("Token did not contain necessary fields")
	}
	if (a.WebClientId != "" && claims.Audience != a.WebClientId) && (a.MobileClientId != "" && claims.Audience != a.MobileClientId) {
		return claims, errors.New("Unknown client id")
	}
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return claims, errors.New("Token expired")
	}

	return claims, nil
}

func getPublicKey(keyId string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyId]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}
