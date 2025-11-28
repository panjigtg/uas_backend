package repository

import (
	"uas/app/models"
)

type AuthRepository interface {
	Register (user *models.Users) error
	GetUserByEmail(identifier string) (*models.Users, error)
	GetUserWithRole(email string) (*models.UserWithRole, error)
	GetPermissionsByUserID(userID string) ([]string, error)
	GetRoleIDByName(name string) (string, error)
}