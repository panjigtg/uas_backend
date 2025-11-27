package config

import (
	"uas/helper"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2"
)

func InitApp() *fiber.App {
	helper.InitLogger()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return helper.InternalServerError(c, err.Error())
		},
	})

	app.Use(cors.New())

	return app
}
