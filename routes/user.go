package routes

import (
	"github.com/gofiber/fiber/v2"
	"uas/middleware"
	"uas/app/services"
)

func UserRoutes(r fiber.Router, userService *services.UserService) {
	users := r.Group("/users")

	users.Use(middleware.AuthRequired())
	users.Use(middleware.RequirePermission("user:manage"))

	users.Get("/", userService.GetAll)
	users.Get("/:id", userService.GetByID)
	users.Post("/", userService.Create)
	users.Put("/:id", userService.Update)
	users.Delete("/:id", userService.Delete)
	users.Put("/:id/role", userService.UpdateRole)
}
