package repo

import (
	"context"
	"database/sql"
	"uas/app/models"
)

type LecturerMockRepo struct {
	CreateFn func(tx *sql.Tx, userID, lecturerID string) error
}

func (m *LecturerMockRepo) Create(tx *sql.Tx, userID, lecturerID string) error {
	return m.CreateFn(tx, userID, lecturerID)
}

func (m *LecturerMockRepo) DeleteByUserID(tx *sql.Tx, userID string) error {
	return nil
}

func (m *LecturerMockRepo) GetIDByUserID(userID string) (string, error) {
	return "", nil
}

func (m *LecturerMockRepo) GetByUserID(ctx context.Context, userID string) (*models.Lecturer, error) {
	return nil, nil
}

func (m *LecturerMockRepo) FindAll(ctx context.Context) ([]models.Lecturer, error) {
	return nil, nil
}
