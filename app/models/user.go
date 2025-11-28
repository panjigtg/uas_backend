package models

import "time"

type Users struct {
	ID           string    `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FullName     string    `db:"full_name"`
	RoleID       string    `db:"role_id"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type UserWithRole struct {
	ID           string `db:"id"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	FullName     string `db:"full_name"`
	PasswordHash string `db:"password_hash"`
	RoleName     string `db:"role"`
	IsActive     bool   `db:"is_active"`
}