package routes

import (
	"github.com/gofiber/fiber/v2"
	"uas/app/services"
	"uas/middleware"
)

func ReportRoutes(r fiber.Router, reportService *services.ReportService) {
	reports := r.Group("/reports")

	reports.Use(middleware.AuthRequired())
	reports.Get(
		"/statistics",
		middleware.RequirePermission("achievement:read"),
		reportService.Statistics,
	)
	
	reports.Get(
		"/student/:id",
		middleware.RequirePermission("achievement:read"),
		reportService.StudentReport,
	)
}