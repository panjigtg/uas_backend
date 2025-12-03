package routes

import (
	"uas/app/services"
	"uas/middleware"

    "github.com/gofiber/fiber/v2"
)

func StudentRoutes(r fiber.Router, svc *services.StudentService) {
	students := r.Group("/students")

	students.Use(middleware.AuthRequired())
	students.Use(middleware.RequirePermission("user:manage"))

    students.Get("/", svc.GetAll)
    students.Get("/:id", svc.GetByID)
    students.Put("/:id/advisor", svc.UpdateAdvisor)
    students.Get("/:id/achievements", svc.GetAchievements)
}
