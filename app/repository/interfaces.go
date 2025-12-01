package repository

import (
	"database/sql"

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
	Create(tx *sql.Tx, user *models.Users) (string, error)
    Update(tx *sql.Tx, userID string, req models.UserUpdateRequest) error
    UpdateRole(tx *sql.Tx, userID string, roleID string) error
	Delete(tx *sql.Tx, id string) error
	GetIDByIndex(idx int) (string, error)
}

type StudentRepository interface {
    Create(tx *sql.Tx, userID string, studentID string) error
    DeleteByUserID(tx *sql.Tx, userID string) error
    RemoveAdvisor(tx *sql.Tx, lecturerID string) error
    GetByUserID(userID string) (*models.Student, error)
}

type LecturerRepository interface {
    Create(tx *sql.Tx, userID string, lecturerID string) error
    DeleteByUserID(tx *sql.Tx, userID string) error
    GetIDByUserID(userID string) (string, error)
}
