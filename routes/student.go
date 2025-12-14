package routes

import (
	"uas/app/services"
	"uas/middleware"

    "github.com/gofiber/fiber/v2"
)

func StudentRoutes(r fiber.Router, studentServices *services.StudentService) {
	students := r.Group("/students")

	students.Use(middleware.AuthRequired())
	students.Use(middleware.RequirePermission("user:manage"))

    students.Get("/", studentServices.GetAll)
    students.Get("/:id", studentServices.GetByID)
    students.Put("/:id/advisor", studentServices.UpdateAdvisor)
    students.Get("/:id/achievements", studentServices.GetAchievements)
}
