package SetupRouters

import (
	"go_cy_form_mailer/pkg/routes/routes"

	"github.com/gofiber/fiber/v3"
)

func SetupRouters(app *fiber.App) {
	routes.Route_mail(app)

}
