package services

import (
	"context"
	"uas/app/models"
	"uas/app/repository"
	"uas/helper"
	"uas/app/dto"

	"github.com/gofiber/fiber/v2"
)

type ReportService struct {
	AchievementRefRepo   repository.AchievementReferenceRepository
	AchievementMongoRepo repository.AchievementMongoRepository
}

func NewReportService(
	refRepo repository.AchievementReferenceRepository,
	mongoRepo repository.AchievementMongoRepository,
) *ReportService {
	return &ReportService{
		AchievementRefRepo:   refRepo,
		AchievementMongoRepo: mongoRepo,
	}
}


// Statistics
// @Summary      Statistik prestasi
// @Description  Menampilkan statistik prestasi global dan per mahasiswa
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} models.MetaInfo
// @Failure      401 {object} models.MetaInfo
// @Failure      403 {object} models.MetaInfo
// @Failure      500 {object} models.MetaInfo
// @Router       /reports/statistics [get]
func (s *ReportService) Statistics(c *fiber.Ctx) error {
	ctx := c.Context()

	refs, err := s.AchievementRefRepo.FindAll(ctx)
	if err != nil {
		return helper.InternalServerError(c, "failed to load report statistics")
	}

	global := map[string]int{
		"total":     0,
		"draft":     0,
		"submitted": 0,
		"verified":  0,
		"rejected":  0,
	}

	studentMap := map[string]*models.StudentStat{}

	for _, ref := range refs {
		global["total"]++
		global[ref.Status]++

		if _, ok := studentMap[ref.StudentID]; !ok {
			studentMap[ref.StudentID] = &models.StudentStat{
				StudentID:   ref.StudentID,
				StudentName: ref.StudentName,
			}
		}

		st := studentMap[ref.StudentID]
		st.Total++

		switch ref.Status {
		case "draft":
			st.Draft++
		case "submitted":
			st.Submitted++
		case "verified":
			st.Verified++
		case "rejected":
			st.Rejected++
		}
	}

	var perStudent []models.StudentStat
	for _, v := range studentMap {
		perStudent = append(perStudent, *v)
	}

	return helper.Success(c, "report statistics retrieved", fiber.Map{
		"global":      global,
		"per_student": perStudent,
	})
}


// StudentReport
// @Summary      Laporan prestasi mahasiswa
// @Description  Menampilkan daftar prestasi dan status milik satu mahasiswa
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Param        id   path string true "Student ID"
// @Success      200 {object} models.MetaInfo
// @Failure      401 {object} models.MetaInfo
// @Failure      403 {object} models.MetaInfo
// @Failure      404 {object} models.MetaInfo
// @Failure      500 {object} models.MetaInfo
// @Router       /reports/student/{id} [get]
func (s *ReportService) StudentReport(c *fiber.Ctx) error {
	ctx := c.Context()
	studentID := c.Params("id")

	refs, err := s.AchievementRefRepo.FindByStudentID(ctx, studentID)
	if err != nil {
		return helper.InternalServerError(c, "failed to load student report")
	}

	// local wrapper → ONLY for this endpoint
	type Item struct {
		Reference   dto.ReportReferenceDTO   `json:"reference"`
		Achievement dto.ReportAchievementDTO `json:"achievement"`
	}

	var items []Item

	for _, ref := range refs {
		ach, err := s.AchievementMongoRepo.FindByID(
			context.Background(),
			ref.MongoAchievementID,
		)
		if err != nil {
			continue
		}

		items = append(items, Item{
			Reference: dto.ReportReferenceDTO{
				ID:          ref.ID,
				Status:      ref.Status,
				StudentID:   ref.StudentID,
				StudentCode: ref.StudentCode,
				StudentName: ref.StudentName,
				SubmittedAt: ref.SubmittedAt,
				VerifiedAt:  ref.VerifiedAt,
			},
			Achievement: dto.ReportAchievementDTO{
				ID:    ach.ID.Hex(),
				Title: ach.Title,
				Type:  ach.AchievementType,
				Point: ach.Points,
			},
		})
	}

	return helper.Success(c, "student report retrieved", fiber.Map{
		"student_id": studentID,
		"total":      len(items),
		"items":      items,
	})
}



