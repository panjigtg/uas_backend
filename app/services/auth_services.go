package services

import (
	"os"
	"strings"
	"uas/app/models"
	"uas/app/repository"
	"uas/helper"
	"uas/utils"
	_ "uas/cmd/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	
)

type AuthService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Register godoc
// @Summary      Register user baru
// @Description  Registrasi user baru dengan role default Mahasiswa
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  models.RegistReq  true  "Register request"
// @Success      201   {object}  models.MetaInfo
// @Failure      400   {object}  models.MetaInfo
// @Failure      500   {object}  models.MetaInfo
// @Router       /auth/register [post]
func (s *AuthService) Register(c *fiber.Ctx) error {
	var req models.RegistReq

	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Format request tidak valid", err.Error())
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		return helper.BadRequest(c, "Semua field wajib diisi", nil)
	}

	// Ambil role mahasiswa dari database
	roleID, err := s.repo.GetRoleIDByName("Mahasiswa")
	if err != nil {
		return helper.InternalServerError(c, "Role default tidak ditemukan")
	}

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

// Login godoc
// @Summary      Login user
// @Description  Login menggunakan email dan password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  models.LoginReq  true  "Login request"
// @Success      200   {object}  models.MetaInfo
// @Failure      401   {object}  models.MetaInfo
// @Failure      403   {object}  models.MetaInfo
// @Router       /auth/login [post]
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req models.LoginReq

	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Format request tidak valid", err.Error())
	}

	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return helper.Unauthorized(c, "Email tidak ditemukan")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return helper.Unauthorized(c, "Password salah")
	}

	if !user.IsActive {
		return helper.Forbidden(c, "Akun tidak aktif")
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
		Token:         token,
		RefreshToken:  refreshToken,
	}

	return helper.Success(c, "Login berhasil", response)
}


// Refresh godoc
// @Summary      Refresh access token
// @Description  Generate access token baru menggunakan refresh token
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200   {object}  models.MetaInfo
// @Failure      401   {object}  models.MetaInfo
// @Router       /auth/refresh [post]
func (s *AuthService) Refresh(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return helper.Unauthorized(c, "Refresh token tidak ditemukan")
	}

	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	if utils.TokenBlacklist[refreshToken] {
		return helper.Unauthorized(c, "Token sudah logout")
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

	newToken, _ := utils.GenerateToken(user.ID, user.RoleName, perms)

	return helper.Success(c, "Token diperbarui", models.RefreshResp{
		AccessToken: newToken,
	})
}


// Logout godoc
// @Summary      Logout user
// @Description  Logout dan blacklist token
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200   {object}  models.MetaInfo
// @Failure      401   {object}  models.MetaInfo
// @Router       /auth/logout [post]
func (s *AuthService) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return helper.Unauthorized(c, "Token tidak ditemukan")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	utils.TokenBlacklist[token] = true

	return helper.Success(c, "Logout berhasil", nil)
}

// Profile godoc
// @Summary      Ambil profil user
// @Description  Mengambil data user yang sedang login
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200   {object}  models.MetaInfo
// @Failure      401   {object}  models.MetaInfo
// @Router       /auth/profile [get]
func (s *AuthService) Profile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	user, err := s.repo.GetUserByID(userID)
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

