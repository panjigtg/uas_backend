package middleware

import (
	"strings"

	"uas/utils"
	"uas/helper"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return helper.Unauthorized(c, "Token tidak ditemukan")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return helper.Unauthorized(c, "Token tidak valid")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role_id", claims.RoleID)
		c.Locals("permissions", claims.Permissions)

		return c.Next()
	}
}
