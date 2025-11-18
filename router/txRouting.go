package router

import (
	"log"

	"ekival-canvas/config"
	"ekival-canvas/handlers/ada_buy_handlers"
	"ekival-canvas/handlers/offers_handlers"
	"ekival-canvas/handlers/submit_handlers"
	"ekival-canvas/handlers/token_buy_handlers"

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
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "EKIVAL-Canvas",
		AppName:       "Ekival Canvas 1.0.0",
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

	adaBuy := tx.Group("/ada-buy")
	adaBuy.Post("/maker-create-order", ada_buy_handlers.MakerCreateAdaBuyOrderHandler)
	adaBuy.Post("/taker-commit-to-order", ada_buy_handlers.TakerCommitToAdaBuyOrderHandler)

	tokenBuy := tx.Group("/token-buy")
	tokenBuy.Post("/maker-create-order", token_buy_handlers.MakerCreateEkiBuyOrderHandler)

	offer := tx.Group("/offer")
	offer.Post("/create-offer", offers_handlers.MakerCreateOfferHandler)

	tx.Post("/tx-submit", submit_handlers.TxSubmitHandler)

	// Running server
	log.Fatal(router.Listen(config.HOST))
}
