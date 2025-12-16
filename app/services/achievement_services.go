package services

import (
	"time"
	"uas/app/models"
	"uas/app/repository"
	"uas/helper"
	"uas/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AchievementService struct {
	StudentRepo  repository.StudentRepository
	MongoRepo    repository.AchievementMongoRepository
	PgRepo       repository.AchievementReferenceRepository
	lecturerRepo repository.LecturerRepository
	UserRepo     repository.UserRepository
}

func NewAchievementService(
	stdRepo repository.StudentRepository,
	mongoRepo repository.AchievementMongoRepository,
	pgRepo repository.AchievementReferenceRepository,
	lecturerRepo repository.LecturerRepository,
	usrRepo repository.UserRepository,
) *AchievementService {
	return &AchievementService{
		StudentRepo:  stdRepo,
		MongoRepo:    mongoRepo,
		PgRepo:       pgRepo,
		lecturerRepo: lecturerRepo,
		UserRepo:     usrRepo,
	}
}

// List achievements
// @Summary      List prestasi
// @Description  Mengambil daftar prestasi sesuai hak akses user
// @Tags         Achievements
// @Security     BearerAuth
// @Param        page   query int false "Page number"
// @Param        limit  query int false "Limit per page"
// @Success      200 {object} models.MetaInfo
// @Failure      401 {object} models.MetaInfo
// @Failure      403 {object} models.MetaInfo
// @Router       /achievements [get]
func (s *AchievementService) List(c *fiber.Ctx) error {
	role := c.Locals("role_id").(string)
	userID := c.Locals("user_id").(string)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var (
		refs         []models.AchievementReference
		total        int
		err          error
		emptyMessage string
	)

	switch role {

	case "Mahasiswa":
		student, err := s.StudentRepo.GetByUserID(c.Context(), userID)
		if err != nil || student == nil {
			return helper.BadRequest(c, "Mahasiswa tidak ditemukan", nil)
		}

		refs, total, err = s.PgRepo.FindByStudentIDPaginated(
			c.Context(),
			student.ID,
			limit,
			offset,
		)
		emptyMessage = "Belum ada prestasi"

	case "Admin":
		refs, total, err = s.PgRepo.FindAllPaginated(
			c.Context(),
			limit,
			offset,
		)
		emptyMessage = "Belum ada prestasi"

	case "Dosen Wali":
		lecturer, err := s.lecturerRepo.GetByUserID(c.Context(), userID)
		if err != nil || lecturer == nil {
			return helper.Forbidden(c, "Anda bukan dosen wali")
		}

		advisees, err := s.StudentRepo.FindByAdvisorID(c.Context(), lecturer.ID)
		if err != nil {
			return helper.InternalServerError(c, "Gagal mengambil data mahasiswa bimbingan")
		}

		if len(advisees) == 0 {
			return helper.Success(
				c,
				"Belum ada prestasi mahasiswa bimbingan",
				[]fiber.Map{},
			)
		}

		studentIDs := make([]string, 0, len(advisees))
		for _, st := range advisees {
			studentIDs = append(studentIDs, st.ID)
		}

		refs, total, err = s.PgRepo.FindForAdvisorPaginated(
			c.Context(),
			studentIDs,
			limit,
			offset,
		)
		emptyMessage = "Belum ada prestasi yang disubmit"

	default:
		return helper.Forbidden(c, "Role tidak memiliki akses")
	}

	if err != nil {
		return helper.InternalServerError(c, "Gagal mengambil data prestasi")
	}

	if len(refs) == 0 {
		return helper.Success(c, emptyMessage, []fiber.Map{})
	}

	list := make([]fiber.Map, 0, len(refs))

	for _, ref := range refs {
		ach, err := s.MongoRepo.FindByID(c.Context(), ref.MongoAchievementID)
		if err != nil || ach == nil {
			continue
		}

		eventDate := ""
		if ach.Details != nil {
			if raw, ok := ach.Details["eventDate"]; ok && raw != nil {
				eventDate = utils.FormatDate(raw)
			}
		}

		thumbnail := ""
		if len(ach.Attachments) > 0 {
			thumbnail = ach.Attachments[0].FileURL
		}

		list = append(list, fiber.Map{
			"id":           ach.ID.Hex(),
			"title":        ach.Title,
			"type":         ach.AchievementType,
			"status":       ref.Status,
			"event_date":   eventDate,
			"thumbnail":    thumbnail,
			"updated_at":   ach.UpdatedAt,
			"student_code": ref.StudentCode,
			"student_name": ref.StudentName,
		})
	}

	meta := models.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData:     total,
		TotalPages: (total + limit - 1) / limit,
	}

	return helper.Paginated(c, "Daftar prestasi ditemukan", list, meta)
}

// Get achievement detail
// @Summary      Detail prestasi
// @Description  Mengambil detail prestasi berdasarkan ID
// @Tags         Achievements
// @Security     BearerAuth
// @Param        id   path string true "Achievement ID"
// @Success      200 {object} models.MetaInfo
// @Failure      404 {object} models.MetaInfo
// @Router       /achievements/{id} [get]
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

	eventDate := utils.FormatDate(details["eventDate"])

	location, _ := details["location"].(string)
	organizer, _ := details["organizer"].(string)

	delete(details, "eventDate")
	delete(details, "location")
	delete(details, "organizer")

	return helper.Success(c, "Detail prestasi ditemukan", fiber.Map{
		"id":          ach.ID.Hex(),
		"title":       ach.Title,
		"type":        ach.AchievementType,
		"description": ach.Description,
		"status":      ref.Status,

		"event_date": eventDate,
		"location":   location,
		"organizer":  organizer,

		"attachments": ach.Attachments,
		"tags":        ach.Tags,
		"details":     details,

		"created_at": ach.CreatedAt,
		"updated_at": ach.UpdatedAt,
	})
}

// Create achievement
// @Summary      Buat prestasi
// @Description  Membuat data prestasi baru
// @Tags         Achievements
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body models.AchievementCreateInput true "Payload prestasi"
// @Success      201 {object} models.MetaInfo
// @Failure      400 {object} models.MetaInfo
// @Router       /achievements [post]
func (s *AchievementService) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	student, err := s.StudentRepo.GetByUserID(c.Context(), userID)
	if err != nil || student == nil {
		return helper.BadRequest(c, "Mahasiswa tidak ditemukan", nil)
	}

	var input models.AchievementCreateInput
	if err := c.BodyParser(&input); err != nil {
		return helper.BadRequest(c, "Format request tidak valid", err.Error())
	}

	if input.AchievementType == "" || input.Title == "" {
		return helper.BadRequest(c, "achievement_type dan title wajib diisi", nil)
	}

	user, err := s.UserRepo.GetByID(userID)
	if err != nil || user == nil {
		return helper.InternalServerError(c, "Gagal mengambil data user")
	}

	sanitized := utils.SanitizeMongoMap(input.Details)

	filtered := utils.FilterDetails(input.AchievementType, sanitized)
	now := time.Now()

	achievement := models.AchievementMongo{
		StudentID:       student.StudentID,
		AchievementType: input.AchievementType,
		Title:           input.Title,
		Description:     input.Description,
		Details:         filtered,
		Tags:            input.Tags,
		Attachments:     []models.AchievementFile{},
		Points:          0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	mongoID, err := s.MongoRepo.Create(c.Context(), &achievement)
	if err != nil {
		return helper.InternalServerError(c, "Gagal menyimpan prestasi")
	}

	ref := models.AchievementReference{
		ID:                 uuid.NewString(),
		StudentID:          student.ID,
		StudentCode:        student.StudentID,
		StudentName:        user.FullName,
		MongoAchievementID: mongoID,
		Status:             utils.AchievementStatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.PgRepo.Create(c.Context(), &ref); err != nil {
		_ = s.MongoRepo.SoftDelete(c.Context(), mongoID)
		return helper.InternalServerError(c, "Gagal menyimpan reference prestasi")
	}

	return helper.Created(c, "Prestasi berhasil dibuat", fiber.Map{
		"id":     mongoID,
		"status": ref.Status,
	})
}

// Submit achievement
// @Summary      Submit prestasi
// @Description  Mengirim prestasi untuk diverifikasi
// @Tags         Achievements
// @Security     BearerAuth
// @Param        id path string true "Achievement ID"
// @Success      200 {object} models.MetaInfo
// @Failure      400 {object} models.MetaInfo
// @Router       /achievements/{id}/submit [post]
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

// Delete achievement
// @Summary      Hapus prestasi
// @Description  Menghapus prestasi draft
// @Tags         Achievements
// @Security     BearerAuth
// @Param        id path string true "Achievement ID"
// @Success      200 {object} models.MetaInfo
// @Failure      403 {object} models.MetaInfo
// @Router       /achievements/{id} [delete]
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

// Verify achievement
// @Summary      Verifikasi prestasi
// @Description  Menerima dan memverifikasi prestasi
// @Tags         Achievements
// @Security     BearerAuth
// @Param        id path string true "Achievement ID"
// @Success      200 {object} models.MetaInfo
// @Failure      400 {object} models.MetaInfo
// @Router       /achievements/{id}/verify [post]
func (s *AchievementService) Verify(c *fiber.Ctx) error {
	id := c.Params("id")
	advisorID := c.Locals("user_id").(string)

	ref, err := s.PgRepo.GetByMongoID(c.Context(), id)
	if err != nil || ref == nil {
		return helper.NotFound(c, "Prestasi tidak ditemukan")
	}

	if ref.Status != utils.AchievementStatusSubmitted {
		return helper.BadRequest(c, "Prestasi belum dikirim atau sudah diverifikasi", nil)
	}

	now := time.Now()
	ref.Status = utils.AchievementStatusVerified
	ref.VerifiedAt = &now
	ref.VerifiedBy = &advisorID

	if err := s.PgRepo.Update(c.Context(), ref); err != nil {
		return helper.InternalServerError(c, "Gagal memverifikasi prestasi")
	}

	return helper.Success(c, "Prestasi berhasil diverifikasi", fiber.Map{
		"status":      ref.Status,
		"verified_at": ref.VerifiedAt,
		"verified_by": ref.VerifiedBy,
	})
}

// Reject achievement
// @Summary      Tolak prestasi
// @Description  Menolak prestasi dengan catatan
// @Tags         Achievements
// @Security     BearerAuth
// @Accept       json
// @Param        id   path string true "Achievement ID"
// @Param        body body object true "Catatan penolakan"
// @Success      200 {object} models.MetaInfo
// @Router       /achievements/{id}/reject [post]
func (s *AchievementService) Reject(c *fiber.Ctx) error {
	id := c.Params("id")
	advisorID := c.Locals("user_id").(string)

	var body struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&body); err != nil {
		return helper.BadRequest(c, "Format request salah", nil)
	}

	ref, err := s.PgRepo.GetByMongoID(c.Context(), id)
	if err != nil || ref == nil {
		return helper.NotFound(c, "Prestasi tidak ditemukan")
	}

	if ref.Status != utils.AchievementStatusSubmitted {
		return helper.BadRequest(c, "Prestasi belum dikirim atau sudah diproses", nil)
	}

	student, err := s.StudentRepo.FindByID(c.Context(), ref.StudentID)
	if err != nil || student == nil {
		return helper.NotFound(c, "Mahasiswa tidak ditemukan")
	}

	// if student.AdvisorID == nil || *student.AdvisorID != advisorID {
	//     return helper.Forbidden(c, "Anda bukan dosen wali mahasiswa ini")
	// }

	now := time.Now()
	ref.Status = utils.AchievementStatusRejected
	ref.RejectionNote = &body.Note
	ref.VerifiedAt = &now
	ref.VerifiedBy = &advisorID

	// Simpan ke database
	if err := s.PgRepo.Update(c.Context(), ref); err != nil {
		return helper.InternalServerError(c, "Gagal menolak prestasi")
	}

	// Response
	return helper.Success(c, "Prestasi berhasil ditolak", fiber.Map{
		"status":         ref.Status,
		"rejection_note": ref.RejectionNote,
		"rejected_at":    ref.VerifiedAt,
		"rejected_by":    ref.VerifiedBy,
	})
}

// Upload achievement attachments
// @Summary      Upload lampiran
// @Description  Upload file prestasi
// @Tags         Achievements
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Param        id    path string true "Achievement ID"
// @Param        files formData file true "Files"
// @Success      200 {object} models.MetaInfo
// @Router       /achievements/{id}/attachments [post]
func (s *AchievementService) UploadAttachments(c *fiber.Ctx) error {
	id := c.Params("id")

	ach, err := s.MongoRepo.FindByID(c.Context(), id)
	if err != nil || ach == nil {
		return helper.NotFound(c, "Prestasi tidak ditemukan")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return helper.BadRequest(c, "Gagal membaca form file", nil)
	}

	files := form.File["files"]
	if len(files) == 0 {
		return helper.BadRequest(c, "Tidak ada file yang diupload", nil)
	}

	uploaded := make([]models.AchievementFile, 0)

	for _, file := range files {
		savePath := "./uploads/achievements" + file.Filename
		if err := c.SaveFile(file, savePath); err != nil {
			return helper.InternalServerError(c, "Gagal menyimpan file")
		}

		uploadedFile := models.AchievementFile{
			FileName:   file.Filename,
			FileURL:    savePath,
			FileType:   file.Header.Get("Content-Type"),
			UploadedAt: time.Now(),
		}
		uploaded = append(uploaded, uploadedFile)
	}

	ach.Attachments = append(ach.Attachments, uploaded...)

	if err := s.MongoRepo.Update(c.Context(), ach); err != nil {
		return helper.InternalServerError(c, "Gagal menambahkan attachment")
	}

	return helper.Success(c, "Attachment berhasil diupload", uploaded)
}

// Achievement history
// @Summary      Riwayat prestasi
// @Description  Mengambil history status prestasi
// @Tags         Achievements
// @Security     BearerAuth
// @Param        id path string true "Achievement ID"
// @Success      200 {object} models.MetaInfo
// @Router       /achievements/{id}/history [get]
func (s *AchievementService) History(c *fiber.Ctx) error {
	id := c.Params("id")

	ref, err := s.PgRepo.GetByMongoID(c.Context(), id)
	if err != nil || ref == nil {
		return helper.NotFound(c, "Prestasi tidak ditemukan")
	}

	history := make([]fiber.Map, 0)

	// 1. Draft (selalu ada)
	history = append(history, fiber.Map{
		"status":    "draft",
		"timestamp": ref.CreatedAt,
	})

	// 2. Submitted
	if ref.SubmittedAt != nil {
		history = append(history, fiber.Map{
			"status":    "submitted",
			"timestamp": ref.SubmittedAt,
		})
	}

	// 3. Verified / Rejected
	if ref.VerifiedAt != nil {
		if ref.Status == utils.AchievementStatusVerified {
			history = append(history, fiber.Map{
				"status":    "verified",
				"timestamp": ref.VerifiedAt,
			})
		}

		if ref.Status == utils.AchievementStatusRejected {
			history = append(history, fiber.Map{
				"status":    "rejected",
				"timestamp": ref.VerifiedAt,
				"note":      ref.RejectionNote,
			})
		}
	}

	return helper.Success(c, "History prestasi ditemukan", history)
}

// Update achievement
// @Summary      Update prestasi
// @Description  Update data prestasi (hanya status draft)
// @Tags         Achievements
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string true "Achievement ID"
// @Param        body body models.AchievementCreateInput true "Payload update prestasi"
// @Success      200 {object} models.MetaInfo
// @Failure      400 {object} models.MetaInfo
// @Failure      403 {object} models.MetaInfo
// @Router       /achievements/{id} [put]
func (s *AchievementService) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	mongoID := c.Params("id")

	ref, err := s.PgRepo.GetByMongoID(c.Context(), mongoID)
	if err != nil || ref == nil {
		return helper.NotFound(c, "Prestasi tidak ditemukan")
	}

	// hanya owner
	student, err := s.StudentRepo.GetByUserID(c.Context(), userID)
	if err != nil || student == nil {
		return helper.BadRequest(c, "Mahasiswa tidak ditemukan", nil)
	}
	if ref.StudentID != student.ID {
		return helper.Forbidden(c, "Tidak dapat mengubah prestasi milik orang lain")
	}

	// hanya draft
	if ref.Status != utils.AchievementStatusDraft {
		return helper.BadRequest(c, "Prestasi hanya bisa diubah saat draft", nil)
	}

	var input models.AchievementCreateInput
	if err := c.BodyParser(&input); err != nil {
		return helper.BadRequest(c, "Format request tidak valid", nil)
	}

	ach, err := s.MongoRepo.FindByID(c.Context(), mongoID)
	if err != nil || ach == nil {
		return helper.NotFound(c, "Data prestasi tidak ditemukan")
	}

	// sanitize & filter details
	sanitized := utils.SanitizeMongoMap(input.Details)
	filtered := utils.FilterDetails(input.AchievementType, sanitized)

	ach.Title = input.Title
	ach.Description = input.Description
	ach.AchievementType = input.AchievementType
	ach.Details = filtered
	ach.Tags = input.Tags
	ach.UpdatedAt = time.Now()

	if err := s.MongoRepo.Update(c.Context(), ach); err != nil {
		return helper.InternalServerError(c, "Gagal update prestasi")
	}

	ref.UpdatedAt = time.Now()
	if err := s.PgRepo.Update(c.Context(), ref); err != nil {
		return helper.InternalServerError(c, "Gagal update reference prestasi")
	}

	return helper.Success(c, "Prestasi berhasil diperbarui", fiber.Map{
		"id": mongoID,
	})
}
