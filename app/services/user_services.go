package services

import (
	"strconv"
	"uas/app/models"
	"uas/app/repository"
	"uas/helper"
	"uas/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func isUUID(s string) bool {
    _, err := uuid.Parse(s)
    return err == nil
}

func (s *UserService) resolveID(id string) (string, error) {
    // Jika angka → convert ke UUID berdasarkan urutan user
    if idx, err := strconv.Atoi(id); err == nil {
        uuid, err := s.repo.GetIDByIndex(idx - 1)
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
	users, err := s.repo.GetAll()
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

    user, err := s.repo.GetByID(resolvedID)
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

	if err := s.repo.Create(&u); err != nil {
		return helper.InternalServerError(c, err.Error())
	}

	return helper.Created(c, "User berhasil dibuat", fiber.Map{
		"username": u.Username,
		"email":    u.Email,
		"role_id":  u.RoleID,
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

    // jalankan update
    err = s.repo.Update(resolvedID, req)
    if err != nil {
        return helper.InternalServerError(c, "Gagal memperbarui user")
    }

    return helper.Success(c, "User berhasil diperbarui", nil)
}



func (s *UserService) Delete(c *fiber.Ctx) error {
    idParam := c.Params("id")

    resolvedID, err := s.resolveID(idParam)
    if err != nil {
        return helper.NotFound(c, "User tidak ditemukan")
    }

    err = s.repo.Delete(resolvedID)
    if err != nil {
        return helper.InternalServerError(c, "Gagal menghapus user")
    }

    return helper.Success(c, "User berhasil dihapus", nil)
}


func (s *UserService) UpdateRole(c *fiber.Ctx) error {
    idParam := c.Params("id")

    resolvedID, err := s.resolveID(idParam)
    if err != nil {
        return helper.NotFound(c, "User tidak ditemukan")
    }

    var req models.UserUpdateRequest
    if err := c.BodyParser(&req); err != nil {
        return helper.BadRequest(c, "Format request tidak valid", err.Error())
    }

    err = s.repo.UpdateRole(resolvedID, req.RoleID)
    if err != nil {
        return helper.InternalServerError(c, "Gagal memperbarui role user")
    }

    return helper.Success(c, "Role user berhasil diperbarui", nil)
}

