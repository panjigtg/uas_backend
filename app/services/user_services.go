package services

import (
	"strconv"
	"uas/app/models"
	"uas/app/repository"
	"uas/helper"
	"uas/utils"
    "database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserService struct {
    DB           *sql.DB
    userRepo     repository.UserRepository
    studentRepo  repository.StudentRepository
    lecturerRepo repository.LecturerRepository
}


func NewUserService(
    db *sql.DB,
    userRepo repository.UserRepository,
    studentRepo repository.StudentRepository,
    lecturerRepo repository.LecturerRepository,
) *UserService {
    return &UserService{
        DB:           db,
        userRepo:     userRepo,
        studentRepo:  studentRepo,
        lecturerRepo: lecturerRepo,
    }
}


func isUUID(s string) bool {
    _, err := uuid.Parse(s)
    return err == nil
}

func (s *UserService) resolveID(id string) (string, error) {
    // Jika angka → convert ke UUID berdasarkan urutan user
    if idx, err := strconv.Atoi(id); err == nil {
        uuid, err := s.userRepo.GetIDByIndex(idx - 1)
        if err != nil {
            return "", err
        }
        return uuid, nil
    }

    // Jika UUID → langsung pakai
    if isUUID(id) {
        return id, nil
    }

    return "", fiber.NewError(fiber.StatusBadRequest, "Format ID tidak valid")
}


func (s *UserService) GetAll(c *fiber.Ctx) error {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return helper.InternalServerError(c, err.Error())
	}

	return helper.Success(c, "Daftar user berhasil diambil", users)
}

func (s *UserService) GetByID(c *fiber.Ctx) error {
    idParam := c.Params("id")

    resolvedID, err := s.resolveID(idParam)
    if err != nil {
        return helper.NotFound(c, "User tidak ditemukan")
    }

    user, err := s.userRepo.GetByID(resolvedID)
    if err != nil {
        return helper.NotFound(c, "User tidak ditemukan")
    }

    return helper.Success(c, "User ditemukan", user)
}



func (s *UserService) Create(c *fiber.Ctx) error {
	var body models.UserCreateRequest

	if err := c.BodyParser(&body); err != nil {
		return helper.BadRequest(c, "Input tidak valid", err.Error())
	}

	hashed, _ := utils.HashPassword(body.Password)

	u := models.Users{
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: hashed,
		FullName:     body.FullName,
		RoleID:       body.RoleID,
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return helper.InternalServerError(c, "Gagal memulai transaksi")
	}

	// CREATE USER
	newUserID, err := s.userRepo.Create(tx, &u)
	if err != nil {
		tx.Rollback()
		return helper.InternalServerError(c, "Gagal membuat user")
	}

	// GENERATE SHORT ID
	genShort := func(prefix string) string {
		return prefix + uuid.New().String()[:8]
	}

	// CREATE PROFILE BASED ON ROLE
	switch body.RoleID {

	case utils.ROLE_MAHASISWA:
		studentID := genShort("STD-")
		if err := s.studentRepo.Create(tx, newUserID, studentID); err != nil {
			tx.Rollback()
			return helper.InternalServerError(c, "Gagal membuat profil mahasiswa")
		}

	case utils.ROLE_DOSEN:
		lecID := genShort("DSN-")
		if err := s.lecturerRepo.Create(tx, newUserID, lecID); err != nil {
			tx.Rollback()
			return helper.InternalServerError(c, "Gagal membuat profil dosen")
		}
	}

	if err := tx.Commit(); err != nil {
		return helper.InternalServerError(c, "Gagal commit transaksi")
	}

	return helper.Created(c, "User berhasil dibuat", fiber.Map{
		"user_id": newUserID,
		"role_id": body.RoleID,
	})
}


func (s *UserService) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")

	resolvedID, err := s.resolveID(idParam)
	if err != nil {
		return helper.NotFound(c, "User tidak ditemukan")
	}

	var req models.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Format request tidak valid", err.Error())
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return helper.InternalServerError(c, "Gagal memulai transaksi")
	}

	if err := s.userRepo.Update(tx, resolvedID, req); err != nil {
		tx.Rollback()
		return helper.InternalServerError(c, "Gagal memperbarui user")
	}

	if err := tx.Commit(); err != nil {
		return helper.InternalServerError(c, "Gagal commit transaksi")
	}

	return helper.Success(c, "User berhasil diperbarui", nil)
}


func (s *UserService) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")

	resolvedID, err := s.resolveID(idParam)
	if err != nil {
		return helper.NotFound(c, "User tidak ditemukan")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return helper.InternalServerError(c, "Gagal memulai transaksi")
	}

	// DELETE PROFILE FIRST (IGNORE ERR karena mungkin tidak punya)
	_ = s.studentRepo.DeleteByUserID(tx, resolvedID)
	_ = s.lecturerRepo.DeleteByUserID(tx, resolvedID)

	// DELETE USER
	if err := s.userRepo.Delete(tx, resolvedID); err != nil {
		tx.Rollback()
		return helper.InternalServerError(c, "Gagal menghapus user")
	}

	if err := tx.Commit(); err != nil {
		return helper.InternalServerError(c, "Gagal commit transaksi")
	}

	return helper.Success(c, "User berhasil dihapus", nil)
}


func (s *UserService) UpdateRole(c *fiber.Ctx) error {
	idParam := c.Params("id")

	resolvedID, err := s.resolveID(idParam)
	if err != nil {
		return helper.NotFound(c, "User tidak ditemukan")
	}

	var req models.UserRoleUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Format request tidak valid", err.Error())
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return helper.InternalServerError(c, "Gagal memulai transaksi")
	}

	_ = s.studentRepo.DeleteByUserID(tx, resolvedID)
	_ = s.lecturerRepo.DeleteByUserID(tx, resolvedID)

	if err := s.userRepo.UpdateRole(tx, resolvedID, req.RoleID); err != nil {
		tx.Rollback()
		return helper.InternalServerError(c, "Gagal update role user")
	}

	genShort := func(prefix string) string {
		return prefix + uuid.New().String()[:8]
	}

	switch req.RoleID {

	case utils.ROLE_MAHASISWA:
		studentID := genShort("STD-")
		if err := s.studentRepo.Create(tx, resolvedID, studentID); err != nil {
			tx.Rollback()
			return helper.InternalServerError(c, "Gagal membuat profil mahasiswa baru")
		}

	case utils.ROLE_DOSEN:
		lecturerID := genShort("DSN-")
		if err := s.lecturerRepo.Create(tx, resolvedID, lecturerID); err != nil {
			tx.Rollback()
			return helper.InternalServerError(c, "Gagal membuat profil dosen baru")
		}
	}

	if err := tx.Commit(); err != nil {
		return helper.InternalServerError(c, "Gagal commit transaksi")
	}

	return helper.Success(c, "Role user berhasil diperbarui", nil)
}
