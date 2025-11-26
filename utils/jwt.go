package utils

import (
	"os"
	"strconv"
	"time"
	"log"

	"uas/app/models"
	"github.com/golang-jwt/jwt/v5"
)

func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set in environment")
	}
	return []byte(secret)
}

func GenerateToken(userID, roleID string, permissions []string) (string, error) {
	expHour, _ := strconv.Atoi(os.Getenv("JWT_EXPIRED"))

	claims := models.JWTClaims{
		UserID:      userID,
		RoleID:      roleID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJWTSecret())
}

func ValidateToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims := token.Claims.(*models.JWTClaims)
	return claims, nil
}
