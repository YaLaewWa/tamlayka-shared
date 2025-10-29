package jwt

import (
	"fmt"

	"github.com/YaLaewWa/tamlayka-shared/apperror"
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func DecodeJWT(inputToken string, jwtSecret []byte) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(inputToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return new(JWTClaims), apperror.UnauthorizedError(err, "parse token failed")
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return new(JWTClaims), apperror.UnauthorizedError(err, "invalid token")
	}

	return claims, nil
}
