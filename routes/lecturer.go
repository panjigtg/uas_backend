package routes

import (
	"github.com/gofiber/fiber/v2"
	"uas/app/services"
	"uas/middleware"
)

func LecturerRoutes(r fiber.Router, lecturerService *services.LecturerService) {
	lecturers := r.Group("/lecturers")

	lecturers.Use(middleware.AuthRequired())

	lecturers.Get("/", middleware.RequirePermission("achievement:verify"), lecturerService.List)
	lecturers.Get("/:id/advisees", middleware.RequirePermission("achievement:verify"), lecturerService.GetMyAdvisees)
}
