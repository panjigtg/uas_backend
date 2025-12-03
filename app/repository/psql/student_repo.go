package psql

import (
	"database/sql"
	"uas/app/models"
	"uas/app/repository"
)

type StudentRepository struct {
	DB *sql.DB
}

func NewStudentRepo(db *sql.DB) repository.StudentRepository {
	return &StudentRepository{DB: db}
}

func (r *StudentRepository) Create(tx *sql.Tx, userID string, studentID string) error {
	query := `
	INSERT INTO students (user_id, student_id)
	VALUES ($1, $2);
	`

	_, err := tx.Exec(query, userID, studentID)
	return err
}

func (r *StudentRepository) DeleteByUserID(tx *sql.Tx, userID string) error {
	_, err := tx.Exec(`DELETE FROM students WHERE user_id=$1`, userID)
	return err
}


func (r *StudentRepository) RemoveAdvisor(tx *sql.Tx, lecturerID string) error {
	_, err := tx.Exec(`
		UPDATE students
		SET advisor_id = NULL
		WHERE advisor_id = $1
	`, lecturerID)
	return err
}

func (r *StudentRepository) GetByUserID(userID string) (*models.Student, error) {
	query := `
	SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
	FROM students
	WHERE user_id = $1
	LIMIT 1;
	`

	var s models.Student
	err := r.DB.QueryRow(query, userID).Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorID,
		&s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *StudentRepository) UpdateAdvisor(tx *sql.Tx, studentID string, advisorID *string) error {
    query := `
        UPDATE students
        SET advisor_id = $1
        WHERE id = $2;
    `
    _, err := tx.Exec(query, advisorID, studentID)
    return err
}

func (r *StudentRepository) GetIDByIndex(idx int) (string, error) {
    query := `
        SELECT id
        FROM students
        ORDER BY created_at ASC
        LIMIT 1 OFFSET $1
    `
    var id string
    err := r.DB.QueryRow(query, idx).Scan(&id)

    if err != nil {
        return "", err
    }

    return id, nil
}

func (r *StudentRepository) GetByStudentID(studentID string) (*models.Student, error) {
    query := `
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
        FROM students
        WHERE id = $1
        LIMIT 1;
    `

    var s models.Student
    err := r.DB.QueryRow(query, studentID).Scan(
        &s.ID,
        &s.UserID,
        &s.StudentID,
        &s.ProgramStudy,
        &s.AcademicYear,
        &s.AdvisorID,
        &s.CreatedAt,
    )

    if err != nil {
        return nil, err
    }

    return &s, nil
}