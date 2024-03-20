package main

import (
	"time"

	SetupRouters "form_mailer/internals/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func main() {
	app := fiber.New(fiber.Config{
		CaseSensitive:      true,
		StrictRouting:      true,
		DisableKeepalive:   true,
		EnableIPValidation: false,
		ReadTimeout:        2 * time.Second,
		TrustedProxies:     []string{"172.18.0.0/16"},
		ProxyHeader:        "X-Forwarded-For",
		ServerHeader:       "***REMOVED*** CY go form_mailer",
		AppName:            "***REMOVED*** CY form_mailer v0.1",
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(requestid.New())
	SetupRouters.SetupRouters(app)
	app.Listen("127.0.0.1:8080")

}
