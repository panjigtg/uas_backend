package helper

import (
	"github.com/gofiber/fiber/v2"
)

func GetPagination(c *fiber.Ctx) (page, limit, offset int) {
	page = c.QueryInt("page", 1)
	limit = c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset = (page - 1) * limit
	return
}


