package psql

import (
    "context"
    "database/sql"
    "uas/app/models"
    "uas/app/repository"
)

type achievementReferenceRepository struct {
    db *sql.DB
}

func NewAchievementReferenceRepository(db *sql.DB) repository.AchievementReferenceRepository {
    return &achievementReferenceRepository{db: db}
}

func (r *achievementReferenceRepository) Create(ctx context.Context, ref *models.AchievementReference) error {
    query := `
        INSERT INTO achievement_references
            (id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at)
        VALUES ($1,$2,$3,$4,NULL,NULL,NULL,NULL,NOW(),NOW())
    `
    _, err := r.db.ExecContext(ctx, query,
        ref.ID,
        ref.StudentID,
        ref.MongoAchievementID,
        ref.Status,
    )
    return err
}

func (r *achievementReferenceRepository) GetByMongoID(ctx context.Context, mongoID string) (*models.AchievementReference, error) {
    query := `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references
        WHERE mongo_achievement_id = $1
        LIMIT 1
    `
    var ref models.AchievementReference

    err := r.db.QueryRowContext(ctx, query, mongoID).Scan(
        &ref.ID,
        &ref.StudentID,
        &ref.MongoAchievementID,
        &ref.Status,
        &ref.SubmittedAt,
        &ref.VerifiedAt,
        &ref.VerifiedBy,
        &ref.RejectionNote,
        &ref.CreatedAt,
        &ref.UpdatedAt,
    )

    if err != nil {
        return nil, err
    }

    return &ref, nil
}

func (r *achievementReferenceRepository) Update(ctx context.Context, ref *models.AchievementReference) error {
    query := `
        UPDATE achievement_references
        SET status = $1,
            submitted_at = $2,
            verified_at = $3,
            verified_by = $4,
            rejection_note = $5,
            updated_at = NOW()
        WHERE mongo_achievement_id = $6
    `
    _, err := r.db.ExecContext(ctx, query,
        ref.Status,
        ref.SubmittedAt,
        ref.VerifiedAt,
        ref.VerifiedBy,
        ref.RejectionNote,
        ref.MongoAchievementID,
    )
    return err
}

func (r *achievementReferenceRepository) FindByStudentID(ctx context.Context, studentID string) ([]models.AchievementReference, error) {
    query := `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at,
               verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references
        WHERE student_id = $1
          AND status != 'deleted'
    `

    rows, err := r.db.QueryContext(ctx, query, studentID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    refs := []models.AchievementReference{}

    for rows.Next() {
        var ref models.AchievementReference
        if err := rows.Scan(
            &ref.ID,
            &ref.StudentID,
            &ref.MongoAchievementID,
            &ref.Status,
            &ref.SubmittedAt,
            &ref.VerifiedAt,
            &ref.VerifiedBy,
            &ref.RejectionNote,
            &ref.CreatedAt,
            &ref.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        refs = append(refs, ref)
    }

    return refs, nil
}

func (r *achievementReferenceRepository) FindAll(ctx context.Context) ([]models.AchievementReference, error) {
    query := `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at,
               verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references
        WHERE status != 'deleted'
    `

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var refs []models.AchievementReference

    for rows.Next() {
        var ref models.AchievementReference
        if err := rows.Scan(
            &ref.ID,
            &ref.StudentID,
            &ref.MongoAchievementID,
            &ref.Status,
            &ref.SubmittedAt,
            &ref.VerifiedAt,
            &ref.VerifiedBy,
            &ref.RejectionNote,
            &ref.CreatedAt,
            &ref.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        refs = append(refs, ref)
    }

    return refs, nil
}
