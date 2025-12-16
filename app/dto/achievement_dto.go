package dto

type ReportAchievementDTO struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"achievement_type"`
	Point int    `json:"points"`
}
