package main

import (
	"log"
	"time"

	env_var "go_cy_form_mailer/cmd/env"
	SetupRouters "go_cy_form_mailer/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/addon/retry"
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
	s := env_var.GetInstance()
	log.Println(s.Data["mail_server_domain"])
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(requestid.New())
	retry.NewExponentialBackoff(retry.Config{
		InitialInterval: 2 * time.Second,
		MaxBackoffTime:  4 * time.Second,
		Multiplier:      2.0,
		MaxRetryCount:   3,
	})
	SetupRouters.SetupRouters(app)
	app.Listen(":8080")
}
