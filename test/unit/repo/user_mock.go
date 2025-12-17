package repo

import (
	"database/sql"
	"uas/app/models"
)

type UserMockRepo struct {
	GetAllFn       func() ([]models.UserWithRole, error)
	GetByIDFn      func(id string) (*models.UserWithRole, error)
	CreateFn       func(tx *sql.Tx, user *models.Users) (string, error)
	UpdateFn       func(tx *sql.Tx, userID string, req models.UserUpdateRequest) error
	UpdateRoleFn   func(tx *sql.Tx, userID string, roleID string) error
	DeleteFn       func(tx *sql.Tx, id string) error
	GetIDByIndexFn func(idx int) (string, error)
}

func (m *UserMockRepo) GetAll() ([]models.UserWithRole, error) {
	if m.GetAllFn == nil {
		return nil, nil
	}
	return m.GetAllFn()
}

func (m *UserMockRepo) GetByID(id string) (*models.UserWithRole, error) {
	if m.GetByIDFn == nil {
		return nil, nil
	}
	return m.GetByIDFn(id)
}

func (m *UserMockRepo) Create(tx *sql.Tx, user *models.Users) (string, error) {
	if m.CreateFn == nil {
		return "", nil
	}
	return m.CreateFn(tx, user)
}

func (m *UserMockRepo) Update(tx *sql.Tx, userID string, req models.UserUpdateRequest) error {
	if m.UpdateFn == nil {
		return nil
	}
	return m.UpdateFn(tx, userID, req)
}

func (m *UserMockRepo) UpdateRole(tx *sql.Tx, userID string, roleID string) error {
	if m.UpdateRoleFn == nil {
		return nil
	}
	return m.UpdateRoleFn(tx, userID, roleID)
}

func (m *UserMockRepo) Delete(tx *sql.Tx, id string) error {
	if m.DeleteFn == nil {
		return nil
	}
	return m.DeleteFn(tx, id)
}

func (m *UserMockRepo) GetIDByIndex(idx int) (string, error) {
	if m.GetIDByIndexFn == nil {
		return "", nil
	}
	return m.GetIDByIndexFn(idx)
}