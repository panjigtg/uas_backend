package repository

import (
	"database/sql"
	"context"

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
    GetByUserID(ctx context.Context, userID string) (*models.Student, error)
	UpdateAdvisor(tx *sql.Tx, studentID string, advisorID *string) error
	GetIDByIndex(idx int) (string, error)
	FindAll(ctx context.Context) ([]models.Student, error)
    FindByID(ctx context.Context, id string) (*models.Student, error)
}

type LecturerRepository interface {
    Create(tx *sql.Tx, userID string, lecturerID string) error
    DeleteByUserID(tx *sql.Tx, userID string) error
    GetIDByUserID(userID string) (string, error)
}

type AchievementMongoRepository interface {
	Create(ctx context.Context, data *models.AchievementMongo) (string, error)
	FindByID(ctx context.Context, id string) (*models.AchievementMongo, error)
	SoftDelete(ctx context.Context, id string) error
}


type AchievementReferenceRepository interface {
	Create(ctx context.Context, ref *models.AchievementReference) error
	// FindByStudent(ctx context.Context, studentID string) ([]models.AchievementReference, error)
	// FindByStudents(ctx context.Context, studentIDs []string) ([]models.AchievementReference, error)
	// FindAll(ctx context.Context) ([]models.AchievementReference, error)
	GetByMongoID(ctx context.Context, mongoID string) (*models.AchievementReference, error)
    Update(ctx context.Context, ref *models.AchievementReference) error
}
