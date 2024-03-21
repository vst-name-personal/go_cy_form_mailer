package handlers

import (
	env_var "go_cy_form_mailer/cmd/env"
	"go_cy_form_mailer/functions"
	"go_cy_form_mailer/internals/datastructs"
	"log"
	"time"

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

	// Parse the HTML template
	bodyContent, err := functions.PrepareMailTemplate(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(functions.GenerateJsonResponse(false, "500 - Ошибка сервера"))
	}

	// Define the mail message
	m := mail.NewMessage()
	m.SetHeader("From", get_var("mail_server_sender"))
	m.SetHeader("To", get_var("mail_server_receiver"))
	m.SetHeader("Subject", "Форма на лендинге")

	// Set the HTML body using the template result
	m.SetBody("text/html", bodyContent.String())
	d := mail.NewDialer(get_var("mail_server_domain"), 587, get_var("mail_server_sender"), get_var("mail_server_passwd"))
	d.LocalName = "go_cy_form_mailer"
	if err := d.DialAndSend(m); err != nil {
		log.Print(d)
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).Send(functions.GenerateJsonResponse(false, "500 - Ошибка сервера"))
	}
	return c.Send(functions.GenerateJsonResponse(true, "200 - Форма успешно отправлена"))
}
func get_var(str string) string {
	// Get singleton instance
	s := env_var.GetInstance()
	// Access data from singleton
	log.Printf(s.Data[str])
	return s.Data[str]
}
