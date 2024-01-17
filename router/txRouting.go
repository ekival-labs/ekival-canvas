package router

import (
	"ekival-canvas/config"
	"ekival-canvas/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// RouterInit initializes the Fiber router with middleware and routes.
// It configures the router, sets up middleware, defines API versioned routes,
// and starts the server listening on the configured host.
func RouterInit() {
	router := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "EKIVAL-Canvas",
		AppName:       "Ekival Server 1.0.0",
	})

	// middlewares
	api := router.Group("/api", logger.New())
	api.Use(helmet.New())
	api.Use(cors.New())
	api.Use(etag.New())
	api.Use(idempotency.New())
	api.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	api.Get("/metrics", monitor.New())

	// api.Use(csrf.New(csrf.Config{
	// 	KeyLookup:      "header:X-Csrf-Token",
	// 	CookieName:     "csrf_",
	// 	CookieSameSite: "Lax",
	// 	Expiration:     1 * time.Hour,
	// 	KeyGenerator:   utils.UUIDv4,
	// }))

	// api.Use(compress.New(compress.Config{
	// 	Level: compress.LevelBestCompression, // 1
	// }))

	version := api.Group("/v1")

	tx := version.Group("/tx")
	tx.Post("/makerCreateAdaBuyOrder", handlers.MakerCreateAdaBuyOrderHandler)
	tx.Post("/makerCreateEkiBuyOrder", handlers.MakerCreateEkiBuyOrderHandler)

	// Running server
	log.Fatal(router.Listen(config.HOST))
}
