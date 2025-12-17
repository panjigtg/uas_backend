package repo

import (
	"context"
	"database/sql"
	"uas/app/models"
)

type StudentMockRepo struct {
	CreateFn          func(tx *sql.Tx, userID string, studentID string) error
	DeleteByUserIDFn  func(tx *sql.Tx, userID string) error
	RemoveAdvisorFn   func(tx *sql.Tx, lecturerID string) error
	GetByUserIDFn     func(ctx context.Context, userID string) (*models.Student, error)
	UpdateAdvisorFn   func(tx *sql.Tx, studentID string, advisorID *string) error
	GetIDByIndexFn    func(idx int) (string, error)
	FindAllFn         func(ctx context.Context) ([]models.Student, error)
	FindByIDFn        func(ctx context.Context, id string) (*models.Student, error)
	FindByAdvisorIDFn func(ctx context.Context, advisorID string) ([]models.Student, error)
	FindAdviseesIDFn  func(ctx context.Context, advisorID string) ([]models.AdviseeResponse, error)
}

func (m *StudentMockRepo) Create(tx *sql.Tx, userID string, studentID string) error {
	if m.CreateFn == nil {
		return nil
	}
	return m.CreateFn(tx, userID, studentID)
}

func (m *StudentMockRepo) DeleteByUserID(tx *sql.Tx, userID string) error {
	if m.DeleteByUserIDFn == nil {
		return nil
	}
	return m.DeleteByUserIDFn(tx, userID)
}

func (m *StudentMockRepo) RemoveAdvisor(tx *sql.Tx, lecturerID string) error {
	if m.RemoveAdvisorFn == nil {
		return nil
	}
	return m.RemoveAdvisorFn(tx, lecturerID)
}

func (m *StudentMockRepo) GetByUserID(ctx context.Context, userID string) (*models.Student, error) {
	if m.GetByUserIDFn == nil {
		return nil, nil
	}
	return m.GetByUserIDFn(ctx, userID)
}

func (m *StudentMockRepo) UpdateAdvisor(tx *sql.Tx, studentID string, advisorID *string) error {
	if m.UpdateAdvisorFn == nil {
		return nil
	}
	return m.UpdateAdvisorFn(tx, studentID, advisorID)
}

func (m *StudentMockRepo) GetIDByIndex(idx int) (string, error) {
	if m.GetIDByIndexFn == nil {
		return "", nil
	}
	return m.GetIDByIndexFn(idx)
}

func (m *StudentMockRepo) FindAll(ctx context.Context) ([]models.Student, error) {
	if m.FindAllFn == nil {
		return nil, nil
	}
	return m.FindAllFn(ctx)
}

func (m *StudentMockRepo) FindByID(ctx context.Context, id string) (*models.Student, error) {
	if m.FindByIDFn == nil {
		return nil, nil
	}
	return m.FindByIDFn(ctx, id)
}

func (m *StudentMockRepo) FindByAdvisorID(ctx context.Context, advisorID string) ([]models.Student, error) {
	if m.FindByAdvisorIDFn == nil {
		return nil, nil
	}
	return m.FindByAdvisorIDFn(ctx, advisorID)
}

func (m *StudentMockRepo) FindAdviseesID(ctx context.Context, advisorID string) ([]models.AdviseeResponse, error) {
	if m.FindAdviseesIDFn == nil {
		return nil, nil
	}
	return m.FindAdviseesIDFn(ctx, advisorID)
}