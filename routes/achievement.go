package routes

import (
	"github.com/gofiber/fiber/v2"
	"uas/middleware"
	"uas/app/services"
)

func AchievementRoutes(r fiber.Router, achievementService *services.AchievementService) {
	achievement := r.Group("/achievements")

	achievement.Use(middleware.AuthRequired())
	
	achievement.Get("/", middleware.RequirePermission("achievement:read"), achievementService.List,)
	achievement.Get("/:id", middleware.RequirePermission("achievement:read"), achievementService.Detail,)
	achievement.Post("/", middleware.RequirePermission("achievement:create"), achievementService.Create)
	achievement.Post("/:id/submit", middleware.RequirePermission("achievement:update"), achievementService.Submit)
	achievement.Delete("/:id", middleware.RequirePermission("achievement:update"), achievementService.Delete)

}
