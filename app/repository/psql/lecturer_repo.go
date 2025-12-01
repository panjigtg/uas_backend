package psql

import (
	"database/sql"
	
	"uas/app/repository"
)

type LecturerRepository struct {
	DB *sql.DB
}

func NewLecturerRepo(db *sql.DB) repository.LecturerRepository {
	return &LecturerRepository{DB: db}
}


func (r *LecturerRepository) Create(tx *sql.Tx, userID string, lecturerID string) error {
	query := `
	INSERT INTO lecturers (user_id, lecturer_id)
	VALUES ($1, $2);
	`

	_, err := tx.Exec(query, userID, lecturerID)
	return err
}


func (r *LecturerRepository) DeleteByUserID(tx *sql.Tx, userID string) error {
	_, err := tx.Exec(`DELETE FROM lecturers WHERE user_id=$1`, userID)
	return err
}


func (r *LecturerRepository) GetIDByUserID(userID string) (string, error) {
	var lecID string
	err := r.DB.QueryRow(`
		SELECT id FROM lecturers WHERE user_id=$1 LIMIT 1
	`, userID).Scan(&lecID)

	return lecID, err
}
