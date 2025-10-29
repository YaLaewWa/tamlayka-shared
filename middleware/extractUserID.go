package middleware

import (
	"errors"
	"strings"

	"github.com/YaLaewWa/tamlayka-shared/apperror"
	"github.com/YaLaewWa/tamlayka-shared/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ExtractUserIDMiddleware struct {
	jwtSecret []byte
}

func NewExtractUserIDMiddleware(jwtSecret string) *ExtractUserIDMiddleware {
	return &ExtractUserIDMiddleware{
		jwtSecret: []byte(jwtSecret),
	}
}

func (m *ExtractUserIDMiddleware) ExtractUserID(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")

	if authHeader == "" {
		return apperror.UnauthorizedError(errors.New("request without authorization header"), "Authorization header is required")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return apperror.UnauthorizedError(errors.New("invalid authorization header"), "Authorization header is invalid")
	}

	token := authHeader[7:]
	claims, err := jwt.DecodeJWT(token, m.jwtSecret)
	if err != nil {
		return apperror.UnauthorizedError(err, "Invalid token")
	}

	id, err := uuid.Parse(claims.ID)
	if err != nil {
		return apperror.BadRequestError(err, "Invalid token")
	}

	ctx.Locals("userID", id)

	return ctx.Next()
}
