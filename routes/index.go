package routes

import (
	"uas/app/services"
	"github.com/gofiber/fiber/v2"
)

func GlobalRoutes(app *fiber.App, authService *services.AuthService) {
	api := app.Group("/api/v1")

	api.Post("/register", authService.Register)
	api.Post("/login", authService.Login)
}
