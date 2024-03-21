package handlers

import (
	"time"

	"go_cy_form_mailer/internals/datastructs"
	"go_cy_form_mailer/pkg/functions"

	"github.com/go-mail/mail"
	"github.com/gofiber/fiber/v3"
)

func Handle_mail(c fiber.Ctx) error {
	name := c.FormValue("NAME")
	phone := c.FormValue("PHONE")
	recaptchaResponse := c.FormValue("cf-turnstile-response")
	if name == "" || phone == "" || recaptchaResponse == "" {
		return c.Status(fiber.StatusBadRequest).Send(functions.GenerateJsonResponse(false, "Заполнены не все поля"))
	}

	cfTurnstileVerified, errMsg := functions.VerifyCloudflareTurnstile(recaptchaResponse, c)
	if !cfTurnstileVerified {
		return c.Status(fiber.StatusUnauthorized).Send(functions.GenerateJsonResponse(false, errMsg))
	}
	data := datastructs.FormData{
		Name:        name,
		Phone:       phone,
		CurrentTime: time.Now().Format("15:04:05 02-01-2006"),
		IP:          functions.GetClientIP(c),
		Admin:       false,
	}

	// Parse the HTML template and get the body content
	bodyContent, err := functions.PrepareMailTemplate(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(functions.GenerateJsonResponse(false, "Ошибка сервера"))
	}

	// Create the mail message
	m := mail.NewMessage()
	m.SetHeader("From", "noreply@cleanyear.ru")
	m.SetHeader("To", "***REMOVED***")
	m.SetHeader("Subject", "Форма на лендинге")

	// Set the HTML body using the template result
	m.SetBody("text/html", bodyContent.String())

	d := mail.NewDialer("***REMOVED***", 587, "noreply@cleanyear.ru", ***REMOVED***)

	if err := d.DialAndSend(m); err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(functions.GenerateJsonResponse(false, "Ошибка сервера"))
	}

	return c.Send(functions.GenerateJsonResponse(true, "Заявка успешно отправлена"))
}
