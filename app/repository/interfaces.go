package repository

import (
	"uas/app/models"
)

type AuthRepository interface {
	Register(user *models.Users) error
	GetUserByEmail(email string) (*models.UserWithRole, error)
	GetUserByID(userID string) (*models.UserWithRole, error)
	GetPermissionsByUserID(userID string) ([]string, error)
	GetRoleIDByName(name string) (string, error) 
}

type UserRepository interface {
	GetAll() ([]models.UserWithRole, error)
	GetByID(id string) (*models.UserWithRole, error)
	Create(user *models.Users) error
	Update(id string, data models.UserUpdateRequest) error 
	Delete(id string) error
	UpdateRole(id string, roleID string) error
	GetIDByIndex(idx int) (string, error)
}
