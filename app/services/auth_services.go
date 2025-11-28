package services

import (
	"os"
	"strings"
	"uas/app/models"
	"uas/app/repository"
	"uas/helper"
	"uas/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
	refreshToken, _ := utils.GenerateRefreshToken(user.ID)
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
		RefreshToken: refreshToken,
	}

	return helper.Success(c, "Login berhasil", response)
}

func (s *AuthService) Refresh(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return helper.Unauthorized(c, "Refresh token tidak ditemukan")
	}

	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	if utils.TokenBlacklist[refreshToken] {
		return helper.Unauthorized(c, "sudah logout")
	}

	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return helper.Unauthorized(c, "Refresh token tidak valid")
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	user, err := s.repo.GetUserByID(userID)
	if err != nil || !user.IsActive {
		return helper.Forbidden(c, "Akun tidak valid atau tidak aktif")
	}

	perms, _ := s.repo.GetPermissionsByUserID(userID)

	newToken, _ := utils.GenerateToken(user.ID, user.RoleID, perms)

	return helper.Success(c, "Token diperbarui", models.RefreshResp{
	Token: newToken,
	})
}

func (s *AuthService) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return helper.Unauthorized(c, "Token tidak ditemukan")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	utils.TokenBlacklist[token] = true

	return helper.Success(c, "Logout berhasil", nil)
}

func (s *AuthService) Profile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	user, err := s.repo.GetUserProfileByID(userID)
	if err != nil {
		return helper.NotFound(c, "User tidak ditemukan")
	}

	perms, _ := s.repo.GetPermissionsByUserID(userID)

	response := models.UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Role:        user.RoleName,
		Permissions: perms,
	}

	return helper.Success(c, "Profil berhasil diambil", response)
}
