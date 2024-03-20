package routes

import (
	"form_mailer/internals/handlers"

	"github.com/gofiber/fiber/v3"
)

func Route_mail(app *fiber.App) {
	app.Post("/api/v1/send-mail", handlers.Handle_mail)

}
