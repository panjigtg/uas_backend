package config

import (
	"uas/routes"

	"github.com/gofiber/fiber/v2"
)

func Bootstrap() *fiber.App {
	app := InitApp()
	db := InitDatabase()

	container := BuildContainer(db.Postgres)

	routes.RegisterRoutes(app, &routes.RouteContainer{
		AuthService: container.AuthService,
		UserService: container.UserService,
		StudentService: container.StudentService,
	})

	return app
}

