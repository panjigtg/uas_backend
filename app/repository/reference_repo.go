package repository

import (
	"context"
	"database/sql"
	"uas/app/models"

	"github.com/lib/pq"
)

type AchievementReferenceRepository interface {
	Create(ctx context.Context, ref *models.AchievementReference) error
    GetByMongoID(ctx context.Context, mongoID string) (*models.AchievementReference, error)
    Update(ctx context.Context, ref *models.AchievementReference) error

	FindByStudentID(ctx context.Context, studentID string) ([]models.AchievementReference, error)
	FindAll(ctx context.Context) ([]models.AchievementReference, error)
	FindByStudentIDs(ctx context.Context, ids []string) ([]models.AchievementReference, error)
    FindForAdvisor(ctx context.Context, ids []string) ([]models.AchievementReference, error)

    FindAllPaginated(ctx context.Context, limit, offset int) ([]models.AchievementReference, int, error)
	FindByStudentIDPaginated(ctx context.Context, studentID string, limit, offset int) ([]models.AchievementReference, int, error)
	FindForAdvisorPaginated(ctx context.Context, ids []string, limit, offset int) ([]models.AchievementReference, int, error)
}

type achievementReferenceRepository struct {
    db *sql.DB
}

func NewAchievementReferenceRepository(db *sql.DB) AchievementReferenceRepository {
    return &achievementReferenceRepository{db: db}
}

func scanAchievementRef(rows *sql.Rows) (models.AchievementReference, error) {
	var ref models.AchievementReference
	err := rows.Scan(
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
		&ref.StudentCode,
		&ref.StudentName,
	)
	return ref, err
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
        SELECT 
            ar.id, ar.student_id, ar.mongo_achievement_id, ar.status,
            ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note,
            ar.created_at, ar.updated_at,
            s.student_id AS student_code,
            u.full_name  AS student_name
        FROM achievement_references ar
        JOIN students s ON s.id = ar.student_id
        JOIN users u    ON u.id = s.user_id
        WHERE ar.student_id = $1
        AND ar.status != 'deleted'
    `

    rows, err := r.db.QueryContext(ctx, query, studentID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var refs []models.AchievementReference

    for rows.Next() {
        var ref models.AchievementReference

        err := rows.Scan(
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
            &ref.StudentCode,
            &ref.StudentName,
        )

        if err != nil {
            return nil, err
        }

        refs = append(refs, ref)
    }

    return refs, nil
}


func (r *achievementReferenceRepository) FindAll(ctx context.Context) ([]models.AchievementReference, error) {
    query := `
                SELECT 
            ar.id,
            ar.student_id,
            ar.mongo_achievement_id,
            ar.status,
            ar.submitted_at,
            ar.verified_at,
            ar.verified_by,
            ar.rejection_note,
            ar.created_at,
            ar.updated_at,
            s.student_id AS student_code,
            u.full_name  AS student_name
        FROM achievement_references ar
        JOIN students s ON s.id = ar.student_id
        JOIN users u    ON u.id = s.user_id
        WHERE ar.status != 'deleted'
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
            &ref.StudentCode,
            &ref.StudentName,
        ); err != nil {
            return nil, err
        }
        refs = append(refs, ref)
    }

    return refs, nil
}

func (r *achievementReferenceRepository) FindByStudentIDs(ctx context.Context, ids []string) ([]models.AchievementReference, error) {
    query := `
                SELECT 
            ar.id, ar.student_id, ar.mongo_achievement_id, ar.status,
            ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note,
            ar.created_at, ar.updated_at,
            s.student_id AS student_code,
            u.full_name  AS student_name
        FROM achievement_references ar
        JOIN students s ON s.id = ar.student_id
        JOIN users u    ON u.id = s.user_id
        WHERE ar.student_id = ANY($1)
        AND ar.status != 'deleted'
    `

    rows, err := r.db.QueryContext(ctx, query, pq.Array(ids))
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
        &ref.StudentCode,
        &ref.StudentName,
    ); err != nil {
        return nil, err
    }
        refs = append(refs, ref)
    }
    return refs, nil
}

func (r *achievementReferenceRepository) FindForAdvisor(ctx context.Context, ids []string) ([]models.AchievementReference, error) {
    query := `
        SELECT 
            ar.id, ar.student_id, ar.mongo_achievement_id, ar.status,
            ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note,
            ar.created_at, ar.updated_at,
            s.student_id AS student_code,
            u.full_name  AS student_name
        FROM achievement_references ar
        JOIN students s ON s.id = ar.student_id
        JOIN users u    ON u.id = s.user_id
        WHERE ar.student_id = ANY($1)
          AND ar.status = 'submitted'
    `

    rows, err := r.db.QueryContext(ctx, query, pq.Array(ids))
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
            &ref.StudentCode,
            &ref.StudentName,
        ); err != nil {
            return nil, err
        }
        refs = append(refs, ref)
    }
    return refs, nil
}


func (r *achievementReferenceRepository) FindAllPaginated(
	ctx context.Context,
	limit, offset int,
) ([]models.AchievementReference, int, error) {

	countQuery := `SELECT COUNT(*) FROM achievement_references WHERE status != 'deleted'`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			ar.id, ar.student_id, ar.mongo_achievement_id, ar.status,
			ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note,
			ar.created_at, ar.updated_at,
			s.student_id AS student_code,
			u.full_name  AS student_name
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		JOIN users u    ON u.id = s.user_id
		WHERE ar.status != 'deleted'
		ORDER BY ar.updated_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var refs []models.AchievementReference
	for rows.Next() {
		ref, err := scanAchievementRef(rows)
		if err != nil {
			return nil, 0, err
		}
		refs = append(refs, ref)
	}

	return refs, total, nil
}


func (r *achievementReferenceRepository) FindByStudentIDPaginated(
	ctx context.Context,
	studentID string,
	limit, offset int,
) ([]models.AchievementReference, int, error) {

	countQuery := `
		SELECT COUNT(*)
		FROM achievement_references
		WHERE student_id = $1 AND status != 'deleted'
	`

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, studentID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			ar.id, ar.student_id, ar.mongo_achievement_id, ar.status,
			ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note,
			ar.created_at, ar.updated_at,
			s.student_id AS student_code,
			u.full_name  AS student_name
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		JOIN users u    ON u.id = s.user_id
		WHERE ar.student_id = $1
		  AND ar.status != 'deleted'
		ORDER BY ar.updated_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var refs []models.AchievementReference
	for rows.Next() {
		ref, err := scanAchievementRef(rows)
		if err != nil {
			return nil, 0, err
		}
		refs = append(refs, ref)
	}

	return refs, total, nil
}

func (r *achievementReferenceRepository) FindForAdvisorPaginated(
	ctx context.Context,
	ids []string,
	limit, offset int,
) ([]models.AchievementReference, int, error) {

	countQuery := `
		SELECT COUNT(*)
		FROM achievement_references
		WHERE student_id = ANY($1)
		  AND status = 'submitted'
	`

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, pq.Array(ids)).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			ar.id, ar.student_id, ar.mongo_achievement_id, ar.status,
			ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note,
			ar.created_at, ar.updated_at,
			s.student_id AS student_code,
			u.full_name  AS student_name
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		JOIN users u    ON u.id = s.user_id
		WHERE ar.student_id = ANY($1)
		  AND ar.status = 'submitted'
		ORDER BY ar.submitted_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(ids), limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var refs []models.AchievementReference
	for rows.Next() {
		ref, err := scanAchievementRef(rows)
		if err != nil {
			return nil, 0, err
		}
		refs = append(refs, ref)
	}

	return refs, total, nil
}
