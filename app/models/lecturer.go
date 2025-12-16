package models

import "time"

type Lecturer struct {
    ID         string    `db:"id" json:"id"`
    UserID     string    `db:"user_id" json:"user_id"`
    LecturerID string    `db:"lecturer_id" json:"lecturer_id"`  
    Department *string    `db:"department" json:"department"`      
    CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type LecturerReq struct {
    ID     string `db:"id"`
    UserID string `db:"user_id"`
}

type AdviseeResponse struct {
	ID            string     `json:"id"`
	StudentID     string     `json:"student_id"`
	Username      string     `json:"username"`
	ProgramStudy  *string    `json:"program_study"`
	AcademicYear  *string    `json:"academic_year"`
	CreatedAt     time.Time  `json:"created_at"`
}