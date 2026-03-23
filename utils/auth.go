package utils

import (
	"errors"

	"github.com/YaLaewWa/tamlayka-shared/apperror"
	"github.com/gofiber/fiber/v2"
)

func AuthorizeRequest(c *fiber.Ctx, id string) error {
	token_id, ok := c.Locals("id").(string)
	if !ok || token_id != id {
		return apperror.UnauthorizedError(errors.New("unauthorized"), "user id not found in context")
	}
	return nil
}
