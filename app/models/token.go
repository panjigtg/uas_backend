package models

import ("github.com/golang-jwt/jwt/v5")

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Users struct {
	id int `json:"id"`
}