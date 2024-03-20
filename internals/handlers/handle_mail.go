package handlers

import (
	"time"

	"bytes"
	_ "embed"
	"html/template"

	"form_mailer/internals/functions"

	"github.com/go-mail/mail"
	"github.com/gofiber/fiber/v3"
)

//go:embed views/mail/mail.html
var mail_template []byte

var Handle_mail = func(c fiber.Ctx) error {
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
	// Parse the embedded HTML template content
	// Parse the HTML template
	tmpl, err := template.New("emailTemplate").Parse(string(mail_template))
	if err != nil {
		panic("Error parsing mail html template")
	}

	// Prepare the data for the template
	data := struct {
		Name        string
		Phone       string
		CurrentTime string
		IP          string
	}{
		Name:        name,
		Phone:       phone,
		CurrentTime: time.Now().Format("15:04:05 02-01-2006"),
		IP:          functions.GetClienIP(c),
	}

	// Execute the template with the data
	var bodyContent bytes.Buffer
	err = tmpl.Execute(&bodyContent, data)
	if err != nil {
		panic("Error executing mail template")
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
