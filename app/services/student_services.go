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
}

func NewStudentService(db *sql.DB, sRepo repository.StudentRepository, lRepo repository.LecturerRepository) *StudentService {
    return &StudentService{
        DB:           db,
        studentRepo:  sRepo,
        lecturerRepo: lRepo,
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

func (s *StudentService) GetAll(c *fiber.Ctx) error {
    list, err := s.studentRepo.FindAll(c.Context())
    if err != nil {
        return helper.InternalServerError(c, err.Error())
    }

    return helper.Success(c, "Daftar mahasiswa ditemukan", list)
}


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

func (s *StudentService) GetAchievements(c *fiber.Ctx) error {
    return helper.Success(c, "Endpoint achievements mahasiswa belum diimplementasikan", nil)
}
