package routes

import (
	"github.com/gofiber/fiber/v2"
	"uas/app/services"
)

type RouteContainer struct {
	AuthService	 		*services.AuthService
	UserService 		*services.UserService
	StudentService 		*services.StudentService
	AchievementService 	*services.AchievementService
}

func RegisterRoutes(app *fiber.App, c *RouteContainer) {
	// Group utama /api/v1
	api := app.Group("/api/v1")

	// Daftarkan masing-masing router
	AuthRoutes(api, c.AuthService)
	UserRoutes(api, c.UserService)
	StudentRoutes(api, c.StudentService)
	AchievementRoutes(api, c.AchievementService)
}
