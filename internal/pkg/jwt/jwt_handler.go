package jwt

import (
	"time"
	"twitter_clone/internal/pkg/apperror"

	"github.com/golang-jwt/jwt/v5"
)

const hmacSampleSecret = "mamad-server"

type CustomClaims struct {
	UserName string `json:"username"`
	ID       int64  `json:"id"`
	jwt.RegisteredClaims
}

func BuildToken(userName string, id int64) (string, *apperror.AppError) {
	claims := CustomClaims{
		UserName: userName,
		ID:       id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", apperror.DB("failed to sign token", err)
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*CustomClaims, *apperror.AppError) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.Validation("unexpected signing method", nil, nil)
		}
		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		return nil, apperror.UnauthorizedErr("failed to parse token", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, apperror.Forbidden("invalid token", nil)
}
