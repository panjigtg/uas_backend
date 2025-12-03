package models

import "time"

type Users struct {
	ID           string    `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
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
	PasswordHash string `db:"password_hash" json:"-"`
	RoleName     string `db:"role"`
	IsActive     bool   `db:"is_active"`
}

type UserCreateRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FullName  string `json:"full_name"`
	RoleID    string `json:"role_id"`
}

type UserUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type UserRoleUpdateRequest struct {
	RoleID string `json:"role_id"`
}