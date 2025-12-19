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
	achievement.Patch("/:id", middleware.RequirePermission("achievement:update"), achievementService.Update,)
	achievement.Post("/:id/submit", middleware.RequirePermission("achievement:update"), achievementService.Submit)
	achievement.Delete("/:id", middleware.RequirePermission("achievement:update"), achievementService.Delete)
	achievement.Post("/:id/verify", middleware.RequirePermission("achievement:verify"), achievementService.Verify,)
	achievement.Post("/:id/reject", middleware.RequirePermission("achievement:verify"), achievementService.Reject,)
	achievement.Post("/:id/attachments", middleware.RequirePermission("achievement:update"), achievementService.UploadAttachments)
	achievement.Get("/:id/history", middleware.RequirePermission("achievement:read"), achievementService.History)
}
