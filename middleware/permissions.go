package middleware

import(
	"uas/utils"
	"strings"
	"uas/helper"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return helper.Unauthorized(c, "Token tidak ditemukan atau format salah")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return helper.Unauthorized(c, "Token tidak valid atau expired")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role_id", claims.RoleID)
		c.Locals("permissions", claims.Permissions)

		return c.Next()
	}
}

	
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		perms, ok := c.Locals("permissions").([]string)
		if !ok {
			return helper.Forbidden(c, "Permissions tidak ditemukan")
		}

		for _, p := range perms {
			if p == permission {
				return c.Next()
			}
		}

		return helper.Forbidden(c, "Anda tidak memiliki izin untuk aksi ini")
	}
}


