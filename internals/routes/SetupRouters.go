package SetupRouters

import (
	"form_mailer/internals/routes/routes"

	"github.com/gofiber/fiber/v3"
)

func SetupRouters(app *fiber.App) {
	routes.Route_mail(app)

}
