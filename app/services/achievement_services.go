package services

import (
	"time"
	"uas/app/models"
	"uas/app/repository"
	"uas/helper"
	"uas/utils"

	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService struct {
	StudentRepo   	repository.StudentRepository
	MongoRepo 		repository.AchievementMongoRepository
	PgRepo    		repository.AchievementReferenceRepository
}

func NewAchievementService(
	stdRepo repository.StudentRepository,
	mongoRepo repository.AchievementMongoRepository,
	pgRepo repository.AchievementReferenceRepository,
) *AchievementService {
	return &AchievementService{
		StudentRepo: stdRepo,
		MongoRepo: mongoRepo,
		PgRepo:    pgRepo,
	}
}

func (s *AchievementService) List(c *fiber.Ctx) error {
    role := c.Locals("role_id").(string)
    userID := c.Locals("user_id").(string)

    var refs []models.AchievementReference
    var err error

    switch role {
    case "Mahasiswa":
        student, _ := s.StudentRepo.GetByUserID(c.Context(), userID)
        refs, err = s.PgRepo.FindByStudentID(c.Context(), student.ID)

    case "Admin":
        refs, err = s.PgRepo.FindAll(c.Context())
    }

    if err != nil {
        return helper.InternalServerError(c, "Gagal mengambil data prestasi")
    }

    var list []fiber.Map

    for _, ref := range refs {
        ach, err := s.MongoRepo.FindByID(c.Context(), ref.MongoAchievementID)
        if err != nil {
            continue
        }

        // ambil eventDate jika ada
		eventDate := ""
		if ach.Details != nil && ach.Details["eventDate"] != nil {
			eventDate = utils.FormatDate(ach.Details["eventDate"])
		}


        // ambil thumbnail → file pertama
        var thumbnail string
        if len(ach.Attachments) > 0 {
            thumbnail = ach.Attachments[0].FileURL
        }

        list = append(list, fiber.Map{
            "id":         ach.ID.Hex(),
            "title":      ach.Title,
            "type":       ach.AchievementType,
            "status":     ref.Status,
            "event_date": eventDate,
            "thumbnail":  thumbnail,
            "updated_at": ach.UpdatedAt,
        })
    }

    return helper.Success(c, "Daftar prestasi ditemukan", list)
}


func (s *AchievementService) Detail(c *fiber.Ctx) error {
    mongoID := c.Params("id")

    ref, err := s.PgRepo.GetByMongoID(c.Context(), mongoID)
    if err != nil || ref == nil {
        return helper.NotFound(c, "Prestasi tidak ditemukan")
    }

    ach, err := s.MongoRepo.FindByID(c.Context(), mongoID)
    if err != nil {
        return helper.InternalServerError(c, "Gagal mengambil data")
    }

    details := ach.Details
    if details == nil {
        details = map[string]interface{}{}
    }

    tags := ach.Tags
    if tags == nil {
        tags = []string{}
    }

    eventDate := utils.FormatDate(details["eventDate"])

    location := ""
    if v, ok := details["location"].(string); ok {
        location = v
    }

    organizer := ""
    if v, ok := details["organizer"].(string); ok {
        organizer = v
    }

    delete(details, "eventDate")
    delete(details, "location")
    delete(details, "organizer")

    return helper.Success(c, "Detail prestasi ditemukan", fiber.Map{
        "id":           ach.ID.Hex(),
        "title":        ach.Title,
        "type":         ach.AchievementType,
        "description":  ach.Description,
        "status":       ref.Status,

        "event_date":   eventDate,
        "location":     location,
        "organizer":    organizer,

        "attachments":  ach.Attachments,
        "tags":         tags,
        "details":      details,

        "created_at":   ach.CreatedAt,
        "updated_at":   ach.UpdatedAt,
    })
}


func (s *AchievementService) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return helper.Unauthorized(c, "User tidak terautentikasi")
	}

	
	student, err := s.StudentRepo.GetByUserID(c.Context(), userID)
	if err != nil {
		return helper.InternalServerError(c, "Gagal mengambil data mahasiswa")
	}
	if student == nil {
		return helper.BadRequest(c, "Mahasiswa tidak ditemukan", nil)
	}

	achievementType := c.FormValue("achievement_type")
	title := c.FormValue("title")
	description := c.FormValue("description")
	tagsStr := c.FormValue("tags")

	
	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
	}

	var details map[string]interface{}
	if raw := c.FormValue("details"); raw != "" {
		_ = json.Unmarshal([]byte(raw), &details)
	}
	
	var attachments []models.AchievementFile

	form, err := c.MultipartForm()
	if err != nil && err != http.ErrNotMultipart {
		return helper.BadRequest(c, "Format multipart tidak valid", err.Error())
	}

	if form != nil && form.File["files[]"] != nil {
		files := form.File["files[]"]

		uploadFolder := "uploads/achievements/" + student.StudentID
		_ = os.MkdirAll(uploadFolder, os.ModePerm)

		for _, file := range files {
			dst := uploadFolder + "/" + strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename

			if saveErr := c.SaveFile(file, dst); saveErr != nil {
				return helper.InternalServerError(c, "Gagal upload file")
			}

			attachments = append(attachments, models.AchievementFile{
				FileName:   file.Filename,
				FileURL:    dst,
				FileType:   file.Header.Get("Content-Type"),
				UploadedAt: time.Now(),
			})
		}
	}

	now := time.Now()


	achievement := models.AchievementMongo{
		ID:              primitive.NewObjectID(),
		StudentID:       student.StudentID,
		AchievementType: achievementType,
		Title:           title,
		Description:     description,
		Details:         details,
		Tags:            tags,
		Attachments:     attachments,
		Points:          0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	mongoID, err := s.MongoRepo.Create(c.Context(), &achievement)
	if err != nil {
		return helper.InternalServerError(c, "Gagal menyimpan data ke MongoDB")
	}

	ref := models.AchievementReference{
		ID:                 uuid.NewString(),
		StudentID:          student.ID,
		MongoAchievementID: mongoID,
		Status:             utils.AchievementStatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.PgRepo.Create(c.Context(), &ref); err != nil {
		return helper.InternalServerError(c, "Gagal menyimpan reference ke PostgreSQL")
	}

	
	return helper.Created(c, "Prestasi berhasil dibuat", fiber.Map{
		"achievement": achievement,
		"status":   ref.Status,
	})
}

func (s *AchievementService) Submit(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	mongoID := c.Params("id")

	ref, err := s.PgRepo.GetByMongoID(c.Context(), mongoID)
	if err != nil || ref == nil {
		return helper.NotFound(c, "Prestasi tidak ditemukan")
	}

	student, err := s.StudentRepo.GetByUserID(c.Context(), userID)
	if err != nil || student == nil {
		return helper.BadRequest(c, "Mahasiswa tidak ditemukan", nil)
	}
	if ref.StudentID != student.ID {
		return helper.Forbidden(c, "Tidak dapat submit prestasi milik pengguna lain")
	}

	if ref.Status != utils.AchievementStatusDraft {
		return helper.BadRequest(c, "Prestasi hanya dapat disubmit dari status draft", nil)
	}

	now := time.Now()
	ref.Status = utils.AchievementStatusSubmitted
	ref.SubmittedAt = &now
	ref.UpdatedAt = now

	if err := s.PgRepo.Update(c.Context(), ref); err != nil {
		return helper.InternalServerError(c, "Gagal update status prestasi")
	}

	return helper.Success(c, "Prestasi berhasil dikirim untuk verifikasi", fiber.Map{
		"status":       ref.Status,
		"submitted_at": ref.SubmittedAt,
	})
}

func (s *AchievementService) Delete(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	mongoID := c.Params("id")

	ref, err := s.PgRepo.GetByMongoID(c.Context(), mongoID)
	if err != nil || ref == nil {
		return helper.NotFound(c, "Prestasi tidak ditemukan")
	}

	student, err := s.StudentRepo.GetByUserID(c.Context(), userID)
	if err != nil || student == nil {
		return helper.BadRequest(c, "Mahasiswa tidak ditemukan", nil)
	}
	if ref.StudentID != student.ID {
		return helper.Forbidden(c, "Tidak dapat menghapus prestasi milik pengguna lain")
	}

	if ref.Status != utils.AchievementStatusDraft {
		return helper.BadRequest(c, "Hanya prestasi draft yang dapat dihapus", nil)
	}

	if err := s.MongoRepo.SoftDelete(c.Context(), ref.MongoAchievementID); err != nil {
		return helper.InternalServerError(c, "Gagal menghapus data MongoDB")
	}

	now := time.Now()
	ref.Status = utils.AchievementStatusDeleted
	ref.UpdatedAt = now

	if err := s.PgRepo.Update(c.Context(), ref); err != nil {
		return helper.InternalServerError(c, "Gagal memperbarui reference di PostgreSQL")
	}

	return helper.Success(c, "Prestasi draft berhasil dihapus", nil)
}
