package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/YaLaewWa/tamlayka-shared/apperror"

	"github.com/gofiber/fiber/v2"
)

type checkTokenResponseBody struct {
	Sub string `json:"sub"`
}

type checkTokenRequestBody struct {
	AccessToken string `json:"access_token"`
}

func checkAccessToken(accessToken string) (string, error) {
	req := checkTokenRequestBody{AccessToken: accessToken}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", apperror.InternalServerError(err, "unable to marshal body")
	}
	url := "https://www.googleapis.com/oauth2/v3/tokeninfo"
	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", apperror.InternalServerError(err, "unable to make a request")
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check the response status code
	if resp.StatusCode == http.StatusOK {

		var result checkTokenResponseBody
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return "", apperror.InternalServerError(err, "Failed to decode json")
		}
		return result.Sub, nil
	} else {
		return "", apperror.UnauthorizedError(err, "Unauthorized header")
	}
}

func Auth(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")

	if authHeader == "" {
		return apperror.UnauthorizedError(errors.New("request without authorization header"), "Authorization header is required")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return apperror.UnauthorizedError(errors.New("invalid authorization header"), "Authorization header is invalid")
	}

	token := authHeader[7:]
	id, err := checkAccessToken(token)
	if err != nil {
		return err
	}

	ctx.Locals("id", id)

	return ctx.Next()
}
