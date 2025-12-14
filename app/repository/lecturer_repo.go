package repository

import (
	"database/sql"
	"uas/app/models"
	"context"
	
)

type LecturerRepository interface {
    Create(tx *sql.Tx, userID string, lecturerID string) error
    DeleteByUserID(tx *sql.Tx, userID string) error
    GetIDByUserID(userID string) (string, error)
	GetByUserID(ctx context.Context, userID string) (*models.Lecturer, error)
}


type lecturerRepository struct {
	DB *sql.DB
}

func NewLecturerRepo(db *sql.DB) LecturerRepository {
	return &lecturerRepository{DB: db}
}


func (r *lecturerRepository) Create(tx *sql.Tx, userID string, lecturerID string) error {
	query := `
	INSERT INTO lecturers (user_id, lecturer_id)
	VALUES ($1, $2);
	`

	_, err := tx.Exec(query, userID, lecturerID)
	return err
}


func (r *lecturerRepository) DeleteByUserID(tx *sql.Tx, userID string) error {
	_, err := tx.Exec(`DELETE FROM lecturers WHERE user_id=$1`, userID)
	return err
}


func (r *lecturerRepository) GetIDByUserID(userID string) (string, error) {
	var lecID string
	err := r.DB.QueryRow(`
		SELECT id FROM lecturers WHERE user_id=$1 LIMIT 1
	`, userID).Scan(&lecID)

	return lecID, err
}

func (r *lecturerRepository) GetByUserID(ctx context.Context, userID string) (*models.Lecturer, error) {
    const query = `
        SELECT id, user_id, lecturer_id, department, created_at
        FROM lecturers
        WHERE user_id = $1
        LIMIT 1
    `

    lec := new(models.Lecturer)

    err := r.DB.QueryRowContext(ctx, query, userID).Scan(
        &lec.ID,
        &lec.UserID,
        &lec.LecturerID,
        &lec.Department,
        &lec.CreatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }

    return lec, err
}