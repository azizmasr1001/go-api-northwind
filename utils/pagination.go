package utils

import "github.com/gofiber/fiber/v2"

func GetPagination(ctx *fiber.Ctx) (int, int) {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)
	return page, limit
}
