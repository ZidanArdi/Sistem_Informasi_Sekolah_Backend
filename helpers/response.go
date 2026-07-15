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

// SwaggerSuccessResponse represents a generic 200 OK success response
type SwaggerSuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Berhasil memproses data"`
	Data    interface{} `json:"data"`
}

// SwaggerCreatedResponse represents a generic 201 Created success response
type SwaggerCreatedResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Data berhasil dibuat"`
	Data    interface{} `json:"data"`
}

// SwaggerError400Response represents a 400 Bad Request error response
type SwaggerError400Response struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Input tidak valid"`
}

// SwaggerError401Response represents a 401 Unauthorized error response
type SwaggerError401Response struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Token tidak ditemukan atau tidak valid"`
}

// SwaggerError403Response represents a 403 Forbidden error response
type SwaggerError403Response struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Akses ditolak: Anda tidak memiliki akses untuk fitur ini"`
	Code    string `json:"code" example:"ERR_PERMISSION_DENIED"`
}

// SwaggerError404Response represents a 404 Not Found error response
type SwaggerError404Response struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Data tidak ditemukan"`
}

// SwaggerError500Response represents a 500 Internal Server Error response
type SwaggerError500Response struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Terjadi kesalahan internal pada server"`
}
