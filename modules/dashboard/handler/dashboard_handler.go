package handler

import (
	"strings"

	"backend/helpers"
	"backend/modules/dashboard/service"

	"github.com/gofiber/fiber/v2"
)

func sendError(c *fiber.Ctx, defaultCode int, err error) error {
	msg := err.Error()
	bizCode := ""
	if strings.Contains(msg, ": ") {
		parts := strings.SplitN(msg, ": ", 2)
		if strings.HasPrefix(parts[0], "ERR_") {
			bizCode = parts[0]
			msg = parts[1]
		}
	}

	code := defaultCode
	if bizCode == "ERR_FORBIDDEN" {
		code = fiber.StatusForbidden
	} else if bizCode == "ERR_PROFILE_NOT_FOUND" {
		code = fiber.StatusNotFound
	} else if bizCode == "ERR_INTERNAL_SERVER" {
		code = fiber.StatusInternalServerError
	}

	if bizCode == "" {
		bizCode = "ERR_INTERNAL_SERVER"
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"code":    bizCode,
		"message": msg,
		"data":    nil,
	})
}

func GetAdminDashboard(c *fiber.Ctx) error {
	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    "ERR_UNAUTHORIZED",
			"message": "Token tidak ditemukan atau tidak valid",
			"data":    nil,
		})
	}
	userID, ok := userIDLocal.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    "ERR_UNAUTHORIZED",
			"message": "Format UserID dalam token tidak valid",
			"data":    nil,
		})
	}

	data, err := service.GetAdminDashboard(userID)
	if err != nil {
		return sendError(c, fiber.StatusForbidden, err)
	}

	return helpers.SuccessResponse(c, "Berhasil memuat dasbor admin", data)
}

func GetGuruDashboard(c *fiber.Ctx) error {
	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    "ERR_UNAUTHORIZED",
			"message": "Token tidak ditemukan atau tidak valid",
			"data":    nil,
		})
	}
	userID, ok := userIDLocal.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    "ERR_UNAUTHORIZED",
			"message": "Format UserID dalam token tidak valid",
			"data":    nil,
		})
	}

	data, err := service.GetGuruDashboard(userID)
	if err != nil {
		return sendError(c, fiber.StatusForbidden, err)
	}

	return helpers.SuccessResponse(c, "Berhasil memuat dasbor guru", data)
}

func GetSiswaDashboard(c *fiber.Ctx) error {
	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    "ERR_UNAUTHORIZED",
			"message": "Token tidak ditemukan atau tidak valid",
			"data":    nil,
		})
	}
	userID, ok := userIDLocal.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    "ERR_UNAUTHORIZED",
			"message": "Format UserID dalam token tidak valid",
			"data":    nil,
		})
	}

	data, err := service.GetSiswaDashboard(userID)
	if err != nil {
		return sendError(c, fiber.StatusForbidden, err)
	}

	return helpers.SuccessResponse(c, "Berhasil memuat dasbor siswa", data)
}

