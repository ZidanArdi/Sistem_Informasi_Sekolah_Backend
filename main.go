package main

import (
	"log"

	"backend/config"
	"backend/middleware"
	authRoute "backend/modules/auth/route"
	guruRoute "backend/modules/guru/route"
	jadwalRoute "backend/modules/jadwal/route"
	kelasRoute "backend/modules/kelas/route"
	mapelRoute "backend/modules/mapel/route"
	nilaiRoute "backend/modules/nilai/route"
	siswaRoute "backend/modules/siswa/route"

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
	authRoute.AuthRoute(app)

	publicAPI := app.Group("/api")
	siswaRoute.SiswaRoute(publicAPI)

	api := app.Group("/api", middleware.JWTProtected)
	guruRoute.GuruRoute(api)
	kelasRoute.KelasRoute(api)
	mapelRoute.MapelRoute(api)
	jadwalRoute.JadwalRoute(api)
	nilaiRoute.NilaiRoute(api)

	// test route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Backend Sistem Informasi Sekolah Berjalan",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
