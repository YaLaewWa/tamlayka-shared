package utils

import (
	"errors"

	"github.com/YaLaewWa/tamlayka-shared/apperror"
	"github.com/gofiber/fiber/v2"
)

func GetUserID(c *fiber.Ctx, id string) (string, error) {
	tokenId, ok := c.Locals("id").(string)
	if !ok || tokenId != id {
		return "", apperror.UnauthorizedError(errors.New("unauthorized"), "user id not found in context")
	}
	return tokenId, nil
}
