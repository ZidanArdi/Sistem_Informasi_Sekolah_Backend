package main

import (
	"log"

	"backend/config"
	"backend/modules/siswa/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

)

func main() {

	// koneksi database
	config.ConnectDB()

	app := fiber.New()

	// middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// route
	route.SiswaRoute(app)

	// test route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Backend Sistem Informasi Sekolah Berjalan",
		})
	})

	log.Fatal(app.Listen(":3000"))
}