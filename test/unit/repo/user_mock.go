package repo

import (
	"database/sql"
	"uas/app/models"
)

type UserMockRepo struct {
	CreateFn func(tx *sql.Tx, user *models.Users) (string, error)
}

func (m *UserMockRepo) Create(tx *sql.Tx, user *models.Users) (string, error) {
	return m.CreateFn(tx, user)
}

func (m *UserMockRepo) GetAll() ([]models.UserWithRole, error) {
	return nil, nil
}

func (m *UserMockRepo) GetByID(id string) (*models.UserWithRole, error) {
	return nil, nil
}

func (m *UserMockRepo) Update(tx *sql.Tx, userID string, req models.UserUpdateRequest) error {
	return nil
}

func (m *UserMockRepo) UpdateRole(tx *sql.Tx, userID string, roleID string) error {
	return nil
}

func (m *UserMockRepo) Delete(tx *sql.Tx, id string) error {
	return nil
}

func (m *UserMockRepo) GetIDByIndex(idx int) (string, error) {
	return "", nil
}
