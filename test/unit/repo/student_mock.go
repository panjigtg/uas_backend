package repo

import (
	"context"
	"database/sql"
	"uas/app/models"
)

type StudentMockRepo struct {
	CreateFn func(tx *sql.Tx, userID, studentID string) error
}

func (m *StudentMockRepo) Create(tx *sql.Tx, userID, studentID string) error {
	return m.CreateFn(tx, userID, studentID)
}

func (m *StudentMockRepo) DeleteByUserID(tx *sql.Tx, userID string) error {
	return nil
}

func (m *StudentMockRepo) RemoveAdvisor(tx *sql.Tx, lecturerID string) error {
	return nil
}

func (m *StudentMockRepo) GetByUserID(ctx context.Context, userID string) (*models.Student, error) {
	return nil, nil
}

func (m *StudentMockRepo) UpdateAdvisor(tx *sql.Tx, studentID string, advisorID *string) error {
	return nil
}

func (m *StudentMockRepo) GetIDByIndex(idx int) (string, error) {
	return "", nil
}

func (m *StudentMockRepo) FindAll(ctx context.Context) ([]models.Student, error) {
	return nil, nil
}

func (m *StudentMockRepo) FindByID(ctx context.Context, id string) (*models.Student, error) {
	return nil, nil
}

func (m *StudentMockRepo) FindByAdvisorID(ctx context.Context, advisorID string) ([]models.Student, error) {
	return nil, nil
}

func (m *StudentMockRepo) FindAdviseesID(ctx context.Context, advisorID string) ([]models.AdviseeResponse, error) {
	return nil, nil
}
