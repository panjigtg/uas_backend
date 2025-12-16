package repo

import "uas/app/models"

type AuthMockRepo struct {
	GetUserByEmailFn         func(email string) (*models.UserWithRole, error)
	GetPermissionsByUserIDFn func(userID string) ([]string, error)
}

func (m *AuthMockRepo) Register(user *models.Users) error {
	return nil
}

func (m *AuthMockRepo) GetUserByEmail(email string) (*models.UserWithRole, error) {
	return m.GetUserByEmailFn(email)
}

func (m *AuthMockRepo) GetUserByID(userID string) (*models.UserWithRole, error) {
	return nil, nil
}

func (m *AuthMockRepo) GetPermissionsByUserID(userID string) ([]string, error) {
	return m.GetPermissionsByUserIDFn(userID)
}

func (m *AuthMockRepo) GetRoleIDByName(name string) (string, error) {
	return "", nil
}