package helpers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(c *fiber.Ctx, code int, message string) error {
	var bizCode string
	if strings.Contains(message, ": ") {
		parts := strings.SplitN(message, ": ", 2)
		if strings.HasPrefix(parts[0], "ERR_") {
			bizCode = parts[0]
			message = parts[1]
		}
	}

	// Override HTTP status code if it's a specific permission or first login block business error code
	if bizCode == "ERR_PERMISSION_DENIED" || bizCode == "ERR_FIRST_LOGIN" {
		code = fiber.StatusForbidden // 403
	}

	res := fiber.Map{
		"success": false,
		"message": message,
	}
	if bizCode != "" {
		res["code"] = bizCode
	}
	return c.Status(code).JSON(res)
}
