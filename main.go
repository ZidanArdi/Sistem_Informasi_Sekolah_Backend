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
	absensiRoute "backend/modules/absensi/route"
	perizinanRoute "backend/modules/perizinan/route"
	dashboardRoute "backend/modules/dashboard/route"
	userRoute "backend/modules/user/route"
	schoolRoute "backend/modules/school/route"

	_ "backend/docs"

	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	swagger "github.com/gofiber/swagger"
)

// @title API Sistem Informasi Sekolah
// @version 1.0
// @description Dokumentasi API backend Sistem Informasi Sekolah menggunakan Golang Fiber, GORM, PostgreSQL, dan JWT.
// @contact.name Praktikum Pemrograman III
// @contact.email praktikum@example.com
// @host 127.0.0.1:3000
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	// koneksi database
	config.ConnectDB()

	app := fiber.New()

	// middleware
	app.Use(cors.New())
	app.Use(logger.New())
	app.Static("/uploads", "./uploads")

	// Swagger UI route
	app.Get("/docs/*", swagger.HandlerDefault)

	// route
	authRoute.AuthRoute(app)

	publicAPI := app.Group("/api")
	siswaRoute.SiswaRoute(publicAPI)
	absensiRoute.AbsensiRoute(publicAPI)

	// Proxy regional data requests to prevent browser CORS and adblocker issues
	publicAPI.Get("/regions/*", func(c *fiber.Ctx) error {
		path := c.Params("*")
		url := "https://emsifa.github.io/api-wilayah-indonesia/api/" + path

		resp, err := http.Get(url)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menghubungi API wilayah: " + err.Error(),
			})
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Gagal membaca data wilayah: " + err.Error(),
			})
		}

		c.Set("Content-Type", "application/json")
		return c.Send(body)
	})

	api := app.Group("/api", middleware.JWTProtected)
	dashboardRoute.DashboardRoute(api)
	guruRoute.GuruRoute(api)
	kelasRoute.KelasRoute(api)
	mapelRoute.MapelRoute(api)
	jadwalRoute.JadwalRoute(api)
	nilaiRoute.NilaiRoute(api)
	perizinanRoute.PerizinanRoute(api)
	userRoute.UserRoute(api)
	schoolRoute.SchoolRoute(api)

	// test route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Backend Sistem Informasi Sekolah Berjalan",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
