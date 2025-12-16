package services_test

import (
	"net/http/httptest"
	"strings"
	"testing"
	"database/sql"

	"uas/app/models"

	"uas/app/services"
	"uas/test/unit/repo"
	"uas/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_Create_Mahasiswa_Success(t *testing.T) {
    // ENV (kalau perlu hashing / jwt)
    t.Setenv("JWT_SECRET", "test")

    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    mock.ExpectBegin()
    mock.ExpectCommit()

    userRepo := &repo.UserMockRepo{
        CreateFn: func(tx *sql.Tx, user *models.Users) (string, error) {
            return "user-123", nil
        },
    }

    studentRepo := &repo.StudentMockRepo{
        CreateFn: func(tx *sql.Tx, userID, studentID string) error {
            return nil
        },
    }

    lecturerRepo := &repo.LecturerMockRepo{}

    service := services.NewUserService(
        db,
        userRepo,
        studentRepo,
        lecturerRepo,
    )

    app := fiber.New()
    app.Post("/users", service.Create)

    body := `{
        "username": "panji",
        "email": "panji@test.com",
        "password": "password123",
        "full_name": "Panji",
        "role_id": "` + utils.ROLE_MAHASISWA + `"
    }`

    req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)
    require.NoError(t, err)

    assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}
