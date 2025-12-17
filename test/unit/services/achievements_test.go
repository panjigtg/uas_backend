package services_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"uas/app/models"
	"uas/app/services"
	"uas/test/unit/repo"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupAchievementService(
	studentRepo *repo.StudentMockRepo,
	userRepo *repo.UserMockRepo,
	mongoRepo *repo.AchievementMongoMockRepo,
	pgRepo *repo.AchievementReferenceMockRepo,
) *services.AchievementService {
	return &services.AchievementService{
		StudentRepo: studentRepo,
		UserRepo:    userRepo,
		MongoRepo:   mongoRepo,
		PgRepo:      pgRepo,
	}
}

func createMultipartRequest(url string, fieldName, fileName string, fileContent []byte) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	if fileName != "" {
		part, err := writer.CreateFormFile(fieldName, fileName)
		if err != nil {
			return nil, err
		}
		part.Write(fileContent)
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestAchievement_Create_Success(t *testing.T) {
	app := fiber.New()
	userID := "user-uuid-123"

	mockStudentRepo := &repo.StudentMockRepo{
		GetByUserIDFn: func(ctx context.Context, uid string) (*models.Student, error) {
			return &models.Student{ID: "student-uuid", StudentID: "NIM123"}, nil
		},
	}

	mockUserRepo := &repo.UserMockRepo{
		GetByIDFn: func(id string) (*models.UserWithRole, error) {
			return &models.UserWithRole{FullName: "Panji Mahasiswa"}, nil
		},
	}

	mockMongoRepo := &repo.AchievementMongoMockRepo{
		CreateFn: func(ctx context.Context, data *models.AchievementMongo) (string, error) {
			assert.Equal(t, "NIM123", data.StudentID)
			assert.Equal(t, "Kompetisi", data.AchievementType)
			
			return "mongo-obj-id-123", nil
		},
	}

	mockPgRepo := &repo.AchievementReferenceMockRepo{
		CreateFn: func(ctx context.Context, ref *models.AchievementReference) error {
			assert.Equal(t, "mongo-obj-id-123", ref.MongoAchievementID)
			assert.Equal(t, "Panji Mahasiswa", ref.StudentName)
			return nil
		},
	}

	svc := setupAchievementService(mockStudentRepo, mockUserRepo, mockMongoRepo, mockPgRepo)

	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return svc.Create(c)
	})

	reqBody := `{
        "achievement_type": "Kompetisi",
        "title": "Juara 1 Golang",
        "description": "Lomba tingkat nasional",
        "tags": ["coding"],
        "details": {"rank": "1"}
    }`

	req := httptest.NewRequest("POST", "/achievements", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestAchievement_Create_StudentNotFound(t *testing.T) {
	app := fiber.New()

	mockStudentRepo := &repo.StudentMockRepo{
		GetByUserIDFn: func(ctx context.Context, uid string) (*models.Student, error) {
			return nil, nil
		},
	}

	svc := setupAchievementService(
		mockStudentRepo,
		&repo.UserMockRepo{},
		&repo.AchievementMongoMockRepo{},
		&repo.AchievementReferenceMockRepo{},
	)

	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user_id", "unknown-user")
		return svc.Create(c)
	})

	req := httptest.NewRequest("POST", "/achievements", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAchievement_Create_ValidationError(t *testing.T) {
	app := fiber.New()

	mockStudentRepo := &repo.StudentMockRepo{
		GetByUserIDFn: func(ctx context.Context, uid string) (*models.Student, error) {
			return &models.Student{ID: "s1"}, nil
		},
	}

	svc := setupAchievementService(
		mockStudentRepo,
		&repo.UserMockRepo{},
		&repo.AchievementMongoMockRepo{},
		&repo.AchievementReferenceMockRepo{},
	)

	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user_id", "u1")
		return svc.Create(c)
	})

	reqBody := `{
        "achievement_type": "", 
        "title": ""
    }`
	req := httptest.NewRequest("POST", "/achievements", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAchievement_Create_MongoError(t *testing.T) {
	app := fiber.New()

	mockStudentRepo := &repo.StudentMockRepo{
		GetByUserIDFn: func(ctx context.Context, uid string) (*models.Student, error) {
			return &models.Student{ID: "s1"}, nil
		},
	}
	mockUserRepo := &repo.UserMockRepo{
		GetByIDFn: func(id string) (*models.UserWithRole, error) {
			return &models.UserWithRole{FullName: "Test"}, nil
		},
	}
	mockMongoRepo := &repo.AchievementMongoMockRepo{
		CreateFn: func(ctx context.Context, data *models.AchievementMongo) (string, error) {
			return "", errors.New("mongo connection lost")
		},
	}

	svc := setupAchievementService(
		mockStudentRepo,
		mockUserRepo,
		mockMongoRepo,
		&repo.AchievementReferenceMockRepo{},
	)

	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user_id", "u1")
		return svc.Create(c)
	})

	reqBody := `{"achievement_type": "A", "title": "T"}`
	req := httptest.NewRequest("POST", "/achievements", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestAchievement_Create_Rollback_On_PgError(t *testing.T) {
	app := fiber.New()
	mockMongoID := "mongo-id-to-delete"

	mockStudentRepo := &repo.StudentMockRepo{
		GetByUserIDFn: func(ctx context.Context, uid string) (*models.Student, error) {
			return &models.Student{ID: "s1", StudentID: "123"}, nil
		},
	}
	mockUserRepo := &repo.UserMockRepo{
		GetByIDFn: func(id string) (*models.UserWithRole, error) {
			return &models.UserWithRole{FullName: "Test User"}, nil
		},
	}

	softDeleteCalled := false

	mockMongoRepo := &repo.AchievementMongoMockRepo{
		CreateFn: func(ctx context.Context, data *models.AchievementMongo) (string, error) {
			return mockMongoID, nil
		},
		SoftDeleteFn: func(ctx context.Context, id string) error {
			softDeleteCalled = true
			assert.Equal(t, mockMongoID, id)
			return nil
		},
	}

	mockPgRepo := &repo.AchievementReferenceMockRepo{
		CreateFn: func(ctx context.Context, ref *models.AchievementReference) error {
			return errors.New("postgres connection error")
		},
	}

	svc := setupAchievementService(mockStudentRepo, mockUserRepo, mockMongoRepo, mockPgRepo)

	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user_id", "u1")
		return svc.Create(c)
	})

	reqBody := `{"achievement_type": "A", "title": "T"}`
	req := httptest.NewRequest("POST", "/achievements", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	assert.True(t, softDeleteCalled, "SoftDelete harus dipanggil jika simpan ke Postgres gagal")
}

func TestAchievement_UploadAttachments_Success(t *testing.T) {
	_ = os.MkdirAll("./uploads", 0755)
	defer os.RemoveAll("./uploads")

	app := fiber.New()
	

    oid := primitive.NewObjectID() 
    targetID := oid.Hex() 

	mockMongoRepo := &repo.AchievementMongoMockRepo{
		FindByIDFn: func(ctx context.Context, id string) (*models.AchievementMongo, error) {
			assert.Equal(t, targetID, id)
			return &models.AchievementMongo{
				ID:          oid, // <--- Masukkan variable tipe primitive.ObjectID, bukan string
				Attachments: []models.AchievementFile{},
			}, nil
		},

		UpdateFn: func(ctx context.Context, a *models.AchievementMongo) error {
			assert.Equal(t, oid, a.ID) 
			assert.Len(t, a.Attachments, 1)
			assert.Equal(t, "test.pdf", a.Attachments[0].FileName)
			return nil
		},
	}

	svc := setupAchievementService(nil, nil, mockMongoRepo, nil)
	app.Post("/achievements/:id/attachments", svc.UploadAttachments)

	req, err := createMultipartRequest(
		"/achievements/"+targetID+"/attachments", 
		"files",
		"test.pdf",
		[]byte("dummy content pdf"),
	)
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
func TestAchievement_UploadAttachments_NotFound(t *testing.T) {
	app := fiber.New()
	targetID := "mongo-id-missing"

	mockMongoRepo := &repo.AchievementMongoMockRepo{
		FindByIDFn: func(ctx context.Context, id string) (*models.AchievementMongo, error) {
			return nil, nil 
		},
	}

	svc := setupAchievementService(nil, nil, mockMongoRepo, nil)
	app.Post("/achievements/:id/attachments", svc.UploadAttachments)

	req := httptest.NewRequest("POST", "/achievements/"+targetID+"/attachments", nil)

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestAchievement_UploadAttachments_NoFileUploaded(t *testing.T) {
	app := fiber.New()
    
	oid := primitive.NewObjectID()
    targetID := oid.Hex()

	mockMongoRepo := &repo.AchievementMongoMockRepo{
		FindByIDFn: func(ctx context.Context, id string) (*models.AchievementMongo, error) {
			return &models.AchievementMongo{ID: oid}, nil // <--- Pakai oid
		},
	}

	svc := setupAchievementService(nil, nil, mockMongoRepo, nil)
	app.Post("/achievements/:id/attachments", svc.UploadAttachments)

	req, _ := createMultipartRequest(
		"/achievements/"+targetID+"/attachments",
		"files",
		"",
		nil,
	)

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAchievement_UploadAttachments_UpdateError(t *testing.T) {
	_ = os.MkdirAll("./uploads", 0755)
	defer os.RemoveAll("./uploads")

	app := fiber.New()
	oid := primitive.NewObjectID()
	targetID := "mongo-id-123"

	mockMongoRepo := &repo.AchievementMongoMockRepo{
		FindByIDFn: func(ctx context.Context, id string) (*models.AchievementMongo, error) {
			return &models.AchievementMongo{ID: oid}, nil
		},
		UpdateFn: func(ctx context.Context, a *models.AchievementMongo) error {
			return errors.New("database update failed")
		},
	}

	svc := setupAchievementService(nil, nil, mockMongoRepo, nil)
	app.Post("/achievements/:id/attachments", svc.UploadAttachments)

	req, _ := createMultipartRequest(
		"/achievements/"+targetID+"/attachments",
		"files",
		"data.txt",
		[]byte("content"),
	)

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}