package services

import (
    "database/sql"
    "strconv"
    "uas/app/models"
    "uas/app/repository"
    "uas/helper"
	"uas/utils"

    "github.com/gofiber/fiber/v2"
)

type StudentService struct {
    DB           *sql.DB
    studentRepo  repository.StudentRepository
    lecturerRepo repository.LecturerRepository
    AchRefRepo      repository.AchievementReferenceRepository
	MongoAchRepo    repository.AchievementMongoRepository
}

func NewStudentService(db *sql.DB, sRepo repository.StudentRepository, lRepo repository.LecturerRepository, achRefRepo repository.AchievementReferenceRepository,
	mongoRepo repository.AchievementMongoRepository) *StudentService {
    return &StudentService{
        DB:           db,
        studentRepo:  sRepo,
        lecturerRepo: lRepo,
        AchRefRepo:  achRefRepo,
		MongoAchRepo: mongoRepo,
    }
}

func (s *StudentService) resolveStudentID(id string) (string, error) {
    // angka → convert index → UUID
    if idx, err := strconv.Atoi(id); err == nil {
        uuid, err := s.studentRepo.GetIDByIndex(idx - 1)
        if err != nil {
            return "", err
        }
        return uuid, nil
    }

    // uuid
    if utils.IsUUID(id) {
    	return id, nil
	}

    return "", fiber.NewError(fiber.StatusBadRequest, "Format ID mahasiswa tidak valid")
}

// GetAll
// @Summary      Ambil semua mahasiswa
// @Description  Menampilkan daftar seluruh mahasiswa
// @Tags         Lecturers & Students
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} models.MetaInfo
// @Failure      401 {object} models.MetaInfo
// @Failure      403 {object} models.MetaInfo
// @Failure      500 {object} models.MetaInfo
// @Router       /students [get]
func (s *StudentService) GetAll(c *fiber.Ctx) error {
    list, err := s.studentRepo.FindAll(c.Context())
    if err != nil {
        return helper.InternalServerError(c, err.Error())
    }

    return helper.Success(c, "Daftar mahasiswa ditemukan", list)
}

// GetByID
// @Summary      Detail mahasiswa
// @Description  Menampilkan detail mahasiswa berdasarkan ID atau index
// @Tags         Lecturers & Students
// @Security     BearerAuth
// @Produce      json
// @Param        id   path string true "Student ID atau index"
// @Success      200 {object} models.MetaInfo
// @Failure      401 {object} models.MetaInfo
// @Failure      403 {object} models.MetaInfo
// @Failure      404 {object} models.MetaInfo
// @Failure      500 {object} models.MetaInfo
// @Router       /students/{id} [get]
func (s *StudentService) GetByID(c *fiber.Ctx) error {
    idParam := c.Params("id")

    resolvedID, err := s.resolveStudentID(idParam)
    if err != nil {
        return helper.NotFound(c, "Mahasiswa tidak ditemukan")
    }

    student, err := s.studentRepo.FindByID(c.Context(), resolvedID)
    if err != nil || student == nil {
        return helper.NotFound(c, "Mahasiswa tidak ditemukan")
    }

    return helper.Success(c, "Data mahasiswa ditemukan", student)
}

// UpdateAdvisor
// @Summary      Update dosen wali mahasiswa
// @Description  Menetapkan atau menghapus dosen wali mahasiswa
// @Tags         Lecturers & Students
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path string true "Student ID atau index"
// @Param        body  body models.UpdateAdvisorRequest true "Data dosen wali"
// @Success      200 {object} models.MetaInfo
// @Failure      400 {object} models.MetaInfo
// @Failure      401 {object} models.MetaInfo
// @Failure      403 {object} models.MetaInfo
// @Failure      404 {object} models.MetaInfo
// @Failure      500 {object} models.MetaInfo
// @Router       /students/{id}/advisor [put]
func (s *StudentService) UpdateAdvisor(c *fiber.Ctx) error {
    idParam := c.Params("id")

    resolvedID, err := s.resolveStudentID(idParam)
    if err != nil {
        return helper.NotFound(c, "Mahasiswa tidak ditemukan")
    }

    var req models.UpdateAdvisorRequest
    if err := c.BodyParser(&req); err != nil {
        return helper.BadRequest(c, "Format request tidak valid", err.Error())
    }

    // mulai transaksi
    tx, err := s.DB.Begin()
    if err != nil {
        return helper.InternalServerError(c, "Gagal memulai transaksi")
    }

    // cek student ada
    student, err := s.studentRepo.FindByID(c.Context(), resolvedID)
    if err != nil {
        tx.Rollback()
        return helper.NotFound(c, "Mahasiswa tidak ditemukan")
    }

    // jika set advisor
    var lecID *string = nil

    if req.AdvisorID != nil {
        id, err := s.lecturerRepo.GetIDByUserID(*req.AdvisorID)
        if err != nil {
            tx.Rollback()
            return helper.NotFound(c, "Dosen wali tidak ditemukan")
        }
        lecID = &id
    }

    // update advisor
    err = s.studentRepo.UpdateAdvisor(tx, resolvedID, lecID)
    if err != nil {
        tx.Rollback()
        return helper.InternalServerError(c, "Gagal update advisor mahasiswa")
    }

    if err := tx.Commit(); err != nil {
        return helper.InternalServerError(c, "Gagal commit transaksi")
    }

    return helper.Success(c, "Advisor mahasiswa berhasil diperbarui", fiber.Map{
        "student_id": student.ID,
        "advisor_id": req.AdvisorID,
    })
}

// GetAchievements
// @Summary      Daftar prestasi mahasiswa
// @Description  Menampilkan daftar prestasi milik mahasiswa
// @Tags         Lecturers & Students
// @Security     BearerAuth
// @Produce      json
// @Param        id   path string true "Student ID atau index"
// @Success      200 {object} models.MetaInfo
// @Failure      401 {object} models.MetaInfo
// @Failure      403 {object} models.MetaInfo
// @Failure      404 {object} models.MetaInfo
// @Failure      500 {object} models.MetaInfo
// @Router       /students/{id}/achievements [get]
func (s *StudentService) GetAchievements(c *fiber.Ctx) error {
	idParam := c.Params("id")

	resolvedID, err := s.resolveStudentID(idParam)
	if err != nil {
		return helper.NotFound(c, "Mahasiswa tidak ditemukan")
	}

	student, err := s.studentRepo.FindByID(c.Context(), resolvedID)
	if err != nil || student == nil {
		return helper.NotFound(c, "Mahasiswa tidak ditemukan")
	}

	refs, err := s.AchRefRepo.FindByStudentID(c.Context(), resolvedID)
	if err != nil {
		return helper.InternalServerError(c, "Gagal mengambil prestasi")
	}

	if len(refs) == 0 {
		return helper.Success(c, "Belum ada prestasi", []fiber.Map{})
	}

	result := []fiber.Map{}

	for _, ref := range refs {
		ach, err := s.MongoAchRepo.FindByID(c.Context(), ref.MongoAchievementID)
		if err != nil || ach == nil {
			continue
		}

		result = append(result, fiber.Map{
			"id":      ach.ID.Hex(),
			"title":   ach.Title,
			"type":    ach.AchievementType,
			"status":  ref.Status,
			"updated": ach.UpdatedAt,
		})
	}

	return helper.Success(c, "Daftar prestasi mahasiswa ditemukan", result)
}


