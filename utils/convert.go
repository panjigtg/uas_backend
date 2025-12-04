package utils

import "uas/app/models"

func MapToDetails(in map[string]interface{}) models.AchievementDetails {
	var d models.AchievementDetails

	if v, ok := in["competition_name"].(string); ok {
		d.CompetitionName = &v
	}

	if v, ok := in["score"].(float64); ok {
		s := int(v)
		d.Score = &s
	}
	
	return d
}
