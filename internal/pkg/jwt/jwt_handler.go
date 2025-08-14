package jwt

import (
	"os"
	"time"
	"twitter_clone/internal/pkg/apperror"

	"github.com/golang-jwt/jwt/v5"
)

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
	hmacSecret := os.Getenv("HMACSECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(hmacSecret))
	if err != nil {
		return "", apperror.DB("failed to sign token", err)
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*CustomClaims, *apperror.AppError) {
	hmacSecret := os.Getenv("HMACSECRET")
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.Validation("unexpected signing method", nil, nil)
		}
		return []byte(hmacSecret), nil
	})

	if err != nil {
		return nil, apperror.UnauthorizedErr("failed to parse token", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, apperror.Forbidden("invalid token", nil)
}
