package config

import (
	"uas/routes"

	"github.com/gofiber/fiber/v2"
)

func Bootstrap() *fiber.App {
	app := InitApp()


	db := InitDatabase()

	repositories := InitRepositories(
		db.Postgres,
		db.Mongo,
	)

	services := InitServices(repositories)

	routes.RegisterRoutes(app, services)

	return app
}
