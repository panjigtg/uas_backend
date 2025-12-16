package services

import (
	"uas/app/repository"
	"uas/helper"

	"github.com/gofiber/fiber/v2"
)

type LecturerService struct {
	studentRepo  repository.StudentRepository
	lecturerRepo repository.LecturerRepository
}

func NewLecturerService(
	lecturerRepo repository.LecturerRepository,
	studentRepo repository.StudentRepository,
) *LecturerService {
	return &LecturerService{
		lecturerRepo: lecturerRepo,
		studentRepo:  studentRepo,
	}
}


// List lecturers
// @Summary      List dosen
// @Description  Mengambil daftar seluruh dosen
// @Tags         Lecturers & Students
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} models.MetaInfo
// @Failure      401 {object} models.MetaInfo
// @Failure      403 {object} models.MetaInfo
// @Router       /lecturers [get]
func (s *LecturerService) GetMyAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	students, err := s.studentRepo.FindAdviseesID(
		c.Context(),
		lecturerID,
	)
	if err != nil {
		return helper.InternalServerError(c, err.Error())
	}

	return helper.Success(c, "Daftar mahasiswa bimbingan dosen", students)
}

// Get lecturer advisees
// @Summary      Daftar mahasiswa bimbingan
// @Description  Mengambil daftar mahasiswa yang dibimbing oleh dosen
// @Tags         Lecturers & Students
// @Security     BearerAuth
// @Produce      json
// @Param        id   path string true "Lecturer ID"
// @Success      200 {object} models.MetaInfo
// @Failure      401 {object} models.MetaInfo
// @Failure      403 {object} models.MetaInfo
// @Failure      404 {object} models.MetaInfo
// @Router       /lecturers/{id}/advisees [get]
func (s *LecturerService) List(c *fiber.Ctx) error {
	lecturers, err := s.lecturerRepo.FindAll(c.Context())
	if err != nil {
		return helper.InternalServerError(c, err.Error())
	}

	return helper.Success(
		c,
		"Daftar dosen berhasil diambil",
		lecturers,
	)
}
