package repo

import (
	"context"
	"uas/app/models"
)

type AchievementMongoMockRepo struct {
	CreateFn     func(ctx context.Context, data *models.AchievementMongo) (string, error)
	FindByIDFn   func(ctx context.Context, id string) (*models.AchievementMongo, error)
	SoftDeleteFn func(ctx context.Context, id string) error
	UpdateFn     func(ctx context.Context, a *models.AchievementMongo) error
}

func (m *AchievementMongoMockRepo) Create(ctx context.Context, data *models.AchievementMongo) (string, error) {
	if m.CreateFn == nil {
		return "", nil
	}
	return m.CreateFn(ctx, data)
}

func (m *AchievementMongoMockRepo) FindByID(ctx context.Context, id string) (*models.AchievementMongo, error) {
	if m.FindByIDFn == nil {
		return nil, nil
	}
	return m.FindByIDFn(ctx, id)
}

func (m *AchievementMongoMockRepo) SoftDelete(ctx context.Context, id string) error {
	if m.SoftDeleteFn == nil {
		return nil
	}
	return m.SoftDeleteFn(ctx, id)
}

func (m *AchievementMongoMockRepo) Update(ctx context.Context, a *models.AchievementMongo) error {
	if m.UpdateFn == nil {
		return nil
	}
	return m.UpdateFn(ctx, a)
}