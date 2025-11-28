package routes

import (
	"uas/app/services"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authService *services.AuthService) {
	api := app.Group("/api/v1/auth")

	api.Post("/register", authService.Register)
	api.Post("/login", authService.Login)
	api.Post("/refresh", authService.Refresh)
	api.Post("/logout", authService.Logout)
}
