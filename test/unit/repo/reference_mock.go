package repo

import (
	"context"
	"uas/app/models"
)

type AchievementReferenceMockRepo struct {
	CreateFn                   func(ctx context.Context, ref *models.AchievementReference) error
	GetByMongoIDFn             func(ctx context.Context, mongoID string) (*models.AchievementReference, error)
	UpdateFn                   func(ctx context.Context, ref *models.AchievementReference) error
	FindByStudentIDFn          func(ctx context.Context, studentID string) ([]models.AchievementReference, error)
	FindAllFn                  func(ctx context.Context) ([]models.AchievementReference, error)
	FindByStudentIDsFn         func(ctx context.Context, ids []string) ([]models.AchievementReference, error)
	FindForAdvisorFn           func(ctx context.Context, ids []string) ([]models.AchievementReference, error)
	FindAllPaginatedFn         func(ctx context.Context, limit, offset int) ([]models.AchievementReference, int, error)
	FindByStudentIDPaginatedFn func(ctx context.Context, studentID string, limit, offset int) ([]models.AchievementReference, int, error)
	FindForAdvisorPaginatedFn  func(ctx context.Context, ids []string, limit, offset int) ([]models.AchievementReference, int, error)
}

func (m *AchievementReferenceMockRepo) Create(ctx context.Context, ref *models.AchievementReference) error {
	if m.CreateFn == nil {
		return nil
	}
	return m.CreateFn(ctx, ref)
}

func (m *AchievementReferenceMockRepo) GetByMongoID(ctx context.Context, mongoID string) (*models.AchievementReference, error) {
	if m.GetByMongoIDFn == nil {
		return nil, nil
	}
	return m.GetByMongoIDFn(ctx, mongoID)
}

func (m *AchievementReferenceMockRepo) Update(ctx context.Context, ref *models.AchievementReference) error {
	if m.UpdateFn == nil {
		return nil
	}
	return m.UpdateFn(ctx, ref)
}

func (m *AchievementReferenceMockRepo) FindByStudentID(ctx context.Context, studentID string) ([]models.AchievementReference, error) {
	if m.FindByStudentIDFn == nil {
		return nil, nil
	}
	return m.FindByStudentIDFn(ctx, studentID)
}

func (m *AchievementReferenceMockRepo) FindAll(ctx context.Context) ([]models.AchievementReference, error) {
	if m.FindAllFn == nil {
		return nil, nil
	}
	return m.FindAllFn(ctx)
}

func (m *AchievementReferenceMockRepo) FindByStudentIDs(ctx context.Context, ids []string) ([]models.AchievementReference, error) {
	if m.FindByStudentIDsFn == nil {
		return nil, nil
	}
	return m.FindByStudentIDsFn(ctx, ids)
}

func (m *AchievementReferenceMockRepo) FindForAdvisor(ctx context.Context, ids []string) ([]models.AchievementReference, error) {
	if m.FindForAdvisorFn == nil {
		return nil, nil
	}
	return m.FindForAdvisorFn(ctx, ids)
}

func (m *AchievementReferenceMockRepo) FindAllPaginated(ctx context.Context, limit, offset int) ([]models.AchievementReference, int, error) {
	if m.FindAllPaginatedFn == nil {
		return nil, 0, nil
	}
	return m.FindAllPaginatedFn(ctx, limit, offset)
}

func (m *AchievementReferenceMockRepo) FindByStudentIDPaginated(ctx context.Context, studentID string, limit, offset int) ([]models.AchievementReference, int, error) {
	if m.FindByStudentIDPaginatedFn == nil {
		return nil, 0, nil
	}
	return m.FindByStudentIDPaginatedFn(ctx, studentID, limit, offset)
}

func (m *AchievementReferenceMockRepo) FindForAdvisorPaginated(ctx context.Context, ids []string, limit, offset int) ([]models.AchievementReference, int, error) {
	if m.FindForAdvisorPaginatedFn == nil {
		return nil, 0, nil
	}
	return m.FindForAdvisorPaginatedFn(ctx, ids, limit, offset)
}