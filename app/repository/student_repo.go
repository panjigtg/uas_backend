package repository

import (
	"database/sql"
	"uas/app/models"

	"context"
)

type StudentRepository interface {
    Create(tx *sql.Tx, userID string, studentID string) error
    DeleteByUserID(tx *sql.Tx, userID string) error
    RemoveAdvisor(tx *sql.Tx, lecturerID string) error
    GetByUserID(ctx context.Context, userID string) (*models.Student, error)
	UpdateAdvisor(tx *sql.Tx, studentID string, advisorID *string) error
	GetIDByIndex(idx int) (string, error)
	FindAll(ctx context.Context) ([]models.Student, error)
    FindByID(ctx context.Context, id string) (*models.Student, error)
	FindByAdvisorID(ctx context.Context, advisorID string) ([]models.Student, error)
}

type studentRepository struct {
	DB *sql.DB
}

func NewStudentRepo(db *sql.DB) StudentRepository {
	return &studentRepository{DB: db}
}

func (r *studentRepository) Create(tx *sql.Tx, userID string, studentID string) error {
	query := `
	INSERT INTO students (user_id, student_id)
	VALUES ($1, $2);
	`

	_, err := tx.Exec(query, userID, studentID)
	return err
}

func (r *studentRepository) DeleteByUserID(tx *sql.Tx, userID string) error {
	_, err := tx.Exec(`DELETE FROM students WHERE user_id=$1`, userID)
	return err
}


func (r *studentRepository) RemoveAdvisor(tx *sql.Tx, lecturerID string) error {
	_, err := tx.Exec(`
		UPDATE students
		SET advisor_id = NULL
		WHERE advisor_id = $1
	`, lecturerID)
	return err
}

func (r *studentRepository) GetByUserID(ctx context.Context, userID string) (*models.Student, error) {
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

func (r *studentRepository) UpdateAdvisor(tx *sql.Tx, studentID string, advisorID *string) error {
    query := `
        UPDATE students
        SET advisor_id = $1
        WHERE id = $2;
    `
    _, err := tx.Exec(query, advisorID, studentID)
    return err
}

func (r *studentRepository) GetIDByIndex(idx int) (string, error) {
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

func (r *studentRepository) GetByStudentID(studentID string) (*models.Student, error) {
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

func (r *studentRepository) FindAll(ctx context.Context) ([]models.Student, error) {
    query := `
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
        FROM students
        ORDER BY created_at ASC
    `

    rows, err := r.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []models.Student
    for rows.Next() {
        var st models.Student
        if err := rows.Scan(
            &st.ID, &st.UserID, &st.StudentID,
            &st.ProgramStudy, &st.AcademicYear,
            &st.AdvisorID, &st.CreatedAt,
        ); err != nil {
            return nil, err
        }
        list = append(list, st)
    }

    return list, nil
}

func (r *studentRepository) FindByID(ctx context.Context, id string) (*models.Student, error) {
    query := `
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
        FROM students
        WHERE id = $1
        LIMIT 1
    `

    var st models.Student
    err := r.DB.QueryRowContext(ctx, query, id).Scan(
        &st.ID, &st.UserID, &st.StudentID,
        &st.ProgramStudy, &st.AcademicYear,
        &st.AdvisorID, &st.CreatedAt,
    )

    if err != nil {
        return nil, err
    }

    return &st, nil
}

func (r *studentRepository) FindByAdvisorID(ctx context.Context, advisorID string) ([]models.Student, error) {
    query := `
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
        FROM students
        WHERE advisor_id = $1
    `

    rows, err := r.DB.QueryContext(ctx, query, advisorID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []models.Student

    for rows.Next() {
        var s models.Student
        if err := rows.Scan(
            &s.ID, &s.UserID, &s.StudentID,
            &s.ProgramStudy, &s.AcademicYear,
            &s.AdvisorID, &s.CreatedAt,
        ); err != nil {
            return nil, err
        }
        list = append(list, s)
    }

    return list, nil
}
