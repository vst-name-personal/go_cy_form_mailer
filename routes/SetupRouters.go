package SetupRouters

import (
	"go_cy_form_mailer/handlers"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

func SetupRouters(app *fiber.App) {
	app.Get("/api/v1/live", healthcheck.NewHealthChecker())
	app.Post("/api/v1/send-mail", handlers.Handle_mail)
	log.Println("Routers initialized")
}
