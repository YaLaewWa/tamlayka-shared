package utils

import (
	"errors"

	"github.com/YaLaewWa/tamlayka-shared/apperror"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ParseIdParam(id string) (uuid.UUID, error) {
	if err := uuid.Validate(id); err != nil {
		return uuid.Nil, apperror.BadRequestError(err, "invalid id")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, apperror.InternalServerError(err, "failed to parse id")
	}

	return parsedID, nil
}

func ParsePaginationQuery(c *fiber.Ctx) (int, int, error) {
	limit := c.QueryInt("limit", 10)
	if limit > 50 || limit <= 0 {
		return 0, 0, apperror.BadRequestError(errors.New("limit has to be between 1 and 50"), "limit has to be between 1 and 50")
	}
	page := c.QueryInt("page", 1)
	if page <= 0 {
		return 0, 0, apperror.BadRequestError(errors.New("page has to be greater than 0"), "page has to be greater than 0")
	}
	return limit, page, nil
}
