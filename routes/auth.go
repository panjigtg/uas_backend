package routes

import (
	"github.com/gofiber/fiber/v2"
	"uas/middleware"
	"uas/app/services"
)

func AuthRoutes(r fiber.Router, authService *services.AuthService) {
	r.Post("/auth/login", authService.Login)
	r.Post("/auth/register", authService.Register)
	r.Post("/auth/refresh", authService.Refresh)
	r.Post("/auth/logout", authService.Logout)

	r.Get("/auth/profile", middleware.AuthRequired(), authService.Profile)
}
