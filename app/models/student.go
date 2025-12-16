package models

import "time"

type Student struct {
    ID           string     `db:"id" json:"id"`
    UserID       string     `db:"user_id" json:"user_id"`
    StudentID    string     `db:"student_id" json:"student_id"`          
    ProgramStudy *string    `db:"program_study" json:"program_study"`     
    AcademicYear *string    `db:"academic_year" json:"academic_year"`     
    AdvisorID    *string    `db:"advisor_id" json:"advisor_id"`           
    CreatedAt    time.Time  `db:"created_at" json:"created_at"`
}

type UpdateAdvisorRequest struct {
    AdvisorID *string `json:"advisor_id"`
}

type StudentStat struct {
    StudentID   string `json:"student_id"`
    StudentName string `json:"student_name"`
    Total       int    `json:"total"`
    Draft       int    `json:"draft"`
    Submitted   int    `json:"submitted"`
    Verified    int    `json:"verified"`
    Rejected    int    `json:"rejected"`
}

type Item struct {
    Reference   any `json:"reference"`
    Achievement any `json:"achievement"`
}