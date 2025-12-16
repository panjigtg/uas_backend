package services_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"uas/app/models"
	"uas/app/services"
	"uas/test/unit/repo"
	"uas/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestAuth_Login_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret")

	app := fiber.New()

	hashed, err := utils.HashPassword("password123")
	require.NoError(t, err)

	mockRepo := &repo.AuthMockRepo{
		GetUserByEmailFn: func(email string) (*models.UserWithRole, error) {
			return &models.UserWithRole{
				ID:           "user-1",
				Email:        email,
				PasswordHash: hashed,
				RoleName:     "Mahasiswa",
				IsActive:     true,
			}, nil
		},
		GetPermissionsByUserIDFn: func(userID string) ([]string, error) {
			return []string{"achievement:read"}, nil
		},
	}

	authService := services.NewAuthService(mockRepo)
	app.Post("/auth/login", authService.Login)

	reqBody := `{
		"email": "panji@test.com",
		"password": "password123"
	}`

	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}



func TestAuth_Login_WrongPassword(t *testing.T) {
	app := fiber.New()

	hashed, err := utils.HashPassword("password123")
	require.NoError(t, err)

	mockRepo := &repo.AuthMockRepo{
		GetUserByEmailFn: func(email string) (*models.UserWithRole, error) {
			return &models.UserWithRole{
				Email:        email,
				PasswordHash: hashed,
				IsActive:     true,
			}, nil
		},
		GetPermissionsByUserIDFn: func(string) ([]string, error) {
			return []string{}, nil
		},
	}

	authService := services.NewAuthService(mockRepo)
	app.Post("/auth/login", authService.Login)

	reqBody := `{
		"email": "panji@test.com",
		"password": "salah"
	}`

	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_Login_InactiveUser(t *testing.T) {
	app := fiber.New()

	hashed, err := utils.HashPassword("password123")
	require.NoError(t, err)

	mockRepo := &repo.AuthMockRepo{
		GetUserByEmailFn: func(email string) (*models.UserWithRole, error) {
			return &models.UserWithRole{
				Email:        email,
				PasswordHash: hashed,
				IsActive:     false,
			}, nil
		},
		GetPermissionsByUserIDFn: func(string) ([]string, error) {
			return []string{}, nil
		},
	}

	authService := services.NewAuthService(mockRepo)
	app.Post("/auth/login", authService.Login)

	reqBody := `{
		"email": "panji@test.com",
		"password": "password123"
	}`

	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}
