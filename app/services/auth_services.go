package services

import (
	"uas/app/models"
	"uas/app/repository"
	"uas/helper"
	"uas/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(c *fiber.Ctx) error {
	var req models.RegistReq

	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Format request tidak valid", err.Error())
	}

	// Validasi sederhana
	if req.Email == "" || req.Password == "" || req.Username == "" {
		return helper.BadRequest(c, "Semua field wajib diisi", nil)
	}

	roleID, err := s.repo.GetRoleIDByName("Mahasiswa")
	if err != nil {
		return helper.InternalServerError(c, "Role default tidak ditemukan")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return helper.InternalServerError(c, "Gagal memproses password")
	}

	user := &models.Users{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		RoleID:       roleID,
	}

	if err := s.repo.Register(user); err != nil {
		return helper.InternalServerError(c, "Gagal mendaftarkan user")
	}

	return helper.Created(c, "Registrasi berhasil", nil)
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req models.LoginReq

	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Format request tidak valid", err.Error())
	}

	user, err := s.repo.GetUserWithRole(req.Email)
	if err != nil {
		return helper.Unauthorized(c, "Email tidak ditemukan")
	}

	if !user.IsActive {
		return helper.Forbidden(c, "Akun tidak aktif")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return helper.Unauthorized(c, "Password salah")
	}

	permissions, _ := s.repo.GetPermissionsByUserID(user.ID)

	token, err := utils.GenerateToken(user.ID, user.RoleName, permissions)
	if err != nil {
		return helper.InternalServerError(c, "Gagal membuat token")
	}

	response := models.LoginResponse{
		User: models.UserProfile{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FullName:    user.FullName,
			Role:        user.RoleName,
			Permissions: permissions,
		},
		Token: token,
	}

	return helper.Success(c, "Login berhasil", response)
}
