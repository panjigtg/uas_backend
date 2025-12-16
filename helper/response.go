package helper

import (
	"encoding/json"
	"uas/app/models"

	"github.com/gofiber/fiber/v2"
)

func logResponse(c *fiber.Ctx, statusCode int, payload interface{}) {
	level := Log.Info()

	if statusCode >= 500 {
		level = Log.Error()
	} else if statusCode >= 400 {
		level = Log.Warn()
	}

	b, _ := json.Marshal(payload)

	level.
		Str("method", c.Method()).
		Str("path", c.Path()).
		Int("status", statusCode).
		Int("size", len(b)).
		Send()
}


func Success(c *fiber.Ctx, message string, data interface{}) error {
	response := models.MetaInfo{
		Status:  "success",
		Message: message,
		Data:    data,
	}

	logResponse(c, fiber.StatusOK, response)

	return c.Status(fiber.StatusOK).JSON(response)
}



func Created(c *fiber.Ctx, message string, data interface{}) error {
	response := models.MetaInfo{
		Status:  "success",
		Message: message,
		Data:    data,
	}

	logResponse(c, fiber.StatusCreated, response)

	return c.Status(fiber.StatusCreated).JSON(response)
}


func BadRequest(c *fiber.Ctx, message string, errors interface{}) error {
	response := models.MetaInfo{
		Status:  "error",
		Message: message,
		Errors:  errors,
	}

	logResponse(c, fiber.StatusBadRequest, response)

	return c.Status(fiber.StatusBadRequest).JSON(response)
}



func Unauthorized(c *fiber.Ctx, message string) error {
	response := models.MetaInfo{
		Status: "error",
		Message: message,
	}

	logResponse(c, fiber.StatusUnauthorized, response)

	return c.Status(fiber.StatusUnauthorized).JSON(response)
}


func Forbidden(c *fiber.Ctx, message string) error {
	response := models.MetaInfo{
		Status: "error",
		Message: message,
	}

	logResponse(c, fiber.StatusForbidden, response)

	return c.Status(fiber.StatusForbidden).JSON(response)
}


func NotFound(c *fiber.Ctx, message string) error {
	response := models.MetaInfo{
		Status:  "error",
		Message: message,
	}

	logResponse(c, fiber.StatusNotFound, response)

	return c.Status(fiber.StatusNotFound).JSON(response)
}



func InternalServerError(c *fiber.Ctx, message string) error {
	response := models.MetaInfo{
		Status: "error",
		Message: message,
	}

	logResponse(c, fiber.StatusInternalServerError, response)

	return c.Status(fiber.StatusInternalServerError).JSON(response)
}

func Paginated(c *fiber.Ctx, message string, data interface{}, meta models.PaginationMeta) error {
	response := models.MetaInfo{
		Status: "success",
		Message: message,
		Meta: meta,
		Data: data,
	}

	logResponse(c, fiber.StatusOK, response)
	return c.Status(fiber.StatusOK).JSON(response)
}	