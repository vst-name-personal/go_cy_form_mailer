package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"text/template"

	"github.com/gofiber/fiber/v3"

	"go_cy_form_mailer/internals/datastructs"

	_ "embed"
)

//go:embed views/mail/mail.html
var mail_template []byte

func VerifyCloudflareTurnstile(recaptchaResponse string, c fiber.Ctx) (bool, string) {
	secret := "1x0000000000000000000000000000000AA"
	remoteIP := GetClientIP(c)
	cfURL := "https://challenges.cloudflare.com/turnstile/v0/siteverify"

	data := map[string]string{
		"secret":   secret,
		"response": recaptchaResponse,
		"remoteip": remoteIP,
	}
	formData := url.Values{}
	for key, value := range data {
		formData.Set(key, value)
	}
	resp, err := http.PostForm(cfURL, formData)
	if err != nil {
		return false, fmt.Sprintf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	if success, ok := response["success"].(bool); ok && success {
		return true, ""
	}

	errorCodes := response["error-codes"].([]interface{})
	for _, code := range errorCodes {
		if codeStr, ok := code.(string); ok {
			switch codeStr {
			case "timeout-or-duplicate":
				return false, "Заявка уже отправлена - перезагрузите страницу для повторной отправки"
			case "11060*":
				return false, "Превышено время ожидания ответа - попробуйте еще раз"
			case "102***", "103***", "104***":
				return false, "Ошибка параметров: пожалуйста, повторите попытку"
			default:
				return false, "Верификация не пройдена - Cloudflare Turnstile"
			}
		}
	}

	return false, ""
}

func GenerateJsonResponse(success bool, message string) []byte {
	response := map[string]interface{}{
		"success": success,
		"message": message,
	}
	jsonResponse, _ := json.Marshal(response)
	return jsonResponse
}

func GetClientIP(c fiber.Ctx) string {
	ip := c.Get("X-REAL-IP")
	if ip == "" {
		ip = c.Get("X-FORWARDED-FOR")
	}
	if ip == "" {
		ip = c.IP()
	}
	return ip
}

func PrepareMailTemplate(data *datastructs.FormData) (*bytes.Buffer, error) {
	// Parse the HTML template
	tmpl, err := template.New("emailTemplate").Parse(string(mail_template))
	if err != nil {
		return nil, err
	}

	// Execute the template with the data
	var bodyContent bytes.Buffer
	if err := tmpl.Execute(&bodyContent, data); err != nil {
		return nil, err
	}
	// Do something with bodyContent if needed

	return &bodyContent, nil
}
