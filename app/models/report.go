package models

type ReportReference struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	StudentID   string `json:"student_id"`
	StudentCode string `json:"student_code"`
	StudentName string `json:"student_name"`
	SubmittedAt any    `json:"submitted_at,omitempty"`
	VerifiedAt  any    `json:"verified_at,omitempty"`
}

type ReportAchievement struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"achievement_type"`
	Point int    `json:"points"`
}
