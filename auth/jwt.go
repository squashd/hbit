package auth

import (
	"time"

	"github.com/SQUASHD/hbit"
	"github.com/golang-jwt/jwt/v5"
)

func MakeJWT(userId, tokenSecret, issuer string, duration time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(duration))),
		Subject:   userId,
	})
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (any, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "invalid token"}
	}
	id, err := token.Claims.GetSubject()
	if err != nil {
		return "", &hbit.Error{Code: hbit.EUNAUTHORIZED, Message: "invalid token"}
	}

	return id, nil
}
