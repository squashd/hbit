package auth

import (
	"time"

	"github.com/SQUASHD/hbit"
	"github.com/golang-jwt/jwt/v5"
)

func MakeJWT(userId, tokenSecret, issuer string, durationSeconds int) (string, error) {
	signingKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(durationSeconds) * time.Second)),
		Subject:   userId,
	})
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret, issuer string) (userId string, err error) {
	claimsStruct := jwt.RegisteredClaims{Issuer: issuer}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
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
