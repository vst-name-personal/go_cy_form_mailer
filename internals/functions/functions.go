package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	_ "embed"

	"github.com/gofiber/fiber/v3"
)

func VerifyCloudflareTurnstile(recaptchaResponse string, c fiber.Ctx) (bool, string) {
	secret := "0x4AAAAAAAQrT-mRFHElTVWWMhXP383seew"
	remoteIP := GetClienIP(c)
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

func GetClienIP(c fiber.Ctx) string {
	ip := c.Get("X-REAL-IP")
	if ip == "" {
		ip = c.Get("X-FORWARDED-FOR")
	}
	if ip == "" {
		ip = c.IP()
	}
	return ip
}
