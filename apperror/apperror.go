package apperror

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError struct {
	Code    int
	Message string
	Err     error
	Stack   string
}

func (e *AppError) Error() string {
	if e.Code/100 == 5 {
		return fmt.Sprintf("Message:%s\nInternal Error: %v \nStack:\n%s", e.Message, e.Err, e.Stack)
	}
	return fmt.Sprintf("Message: %s | Internal Error: %v", e.Message, e.Err)
}

func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

func New(code int, message string, err error) *AppError {
	stack := captureStack()
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Stack:   stack,
	}
}

// modified https://github.com/pkg/errors/blob/5dd12d0cfe7f152f80558d591504ce685299311e/stack.go
func captureStack() string {
	const depth = 16
	var pcs [depth]uintptr

	// skip 4 frames apperror.captureStack x2, apperror.New, apperror.InternalServerError or another
	n := runtime.Callers(4, pcs[:])

	if n == 0 {
		return "stack trace is not available"
	}

	stackTrace := make([]string, 0, n)
	for i := range n {
		fn := runtime.FuncForPC(pcs[i])
		file, line := fn.FileLine(pcs[i])
		stackTrace = append(stackTrace, fmt.Sprintf("%s\n\tat %s:%d", fn.Name(), file, line))
	}

	return strings.Join(stackTrace, "\n")
}

func InternalServerError(err error, msg string) *AppError {
	return New(fiber.StatusInternalServerError, msg, err)
}

func BadRequestError(err error, msg string) *AppError {
	return New(fiber.StatusBadRequest, msg, err)
}

func UnauthorizedError(err error, msg string) *AppError {
	return New(fiber.StatusUnauthorized, msg, err)
}

func ForbiddenError(err error, msg string) *AppError {
	return New(fiber.StatusForbidden, msg, err)
}

func NotFoundError(err error, msg string) *AppError {
	return New(fiber.StatusNotFound, msg, err)
}

func ConflictError(err error, msg string) *AppError {
	return New(fiber.StatusConflict, msg, err)
}

func UnprocessableEntityError(err error, msg string) *AppError {
	return New(fiber.StatusUnprocessableEntity, msg, err)
}

func ServiceUnavailableError(err error, msg string) *AppError {
	return New(fiber.StatusServiceUnavailable, msg, err)
}

func ErrorHandler(c *fiber.Ctx, err error) error {

	// if is app error
	if IsAppError(err) {
		e := err.(*AppError)
		if err := c.Status(e.Code).JSON(fiber.Map{"error": e.Message}); err != nil {
			// if can't send error -- it should not be able
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		return nil
	}

	var e *fiber.Error
	if errors.As(err, &e) {
		if err := c.Status(e.Code).JSON(fiber.Map{"error": e.Error()}); err != nil {
			// if can't send error -- it should not be able
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		return nil
	}

	// other case return error that is not fiber error or app error
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
}

func ConvertAppErrorToGRPC(err error) error {
	if !IsAppError(err) {
		return status.Error(codes.Internal, "internal server error")
	}

	appErr := err.(*AppError)
	switch appErr.Code {
	case fiber.StatusInternalServerError:
		return status.Error(codes.Internal, appErr.Message)
	case fiber.StatusBadRequest:
		return status.Error(codes.InvalidArgument, appErr.Message)
	case fiber.StatusUnauthorized:
		return status.Error(codes.Unauthenticated, appErr.Message)
	case fiber.StatusForbidden:
		return status.Error(codes.PermissionDenied, appErr.Message)
	case fiber.StatusNotFound:
		return status.Error(codes.NotFound, appErr.Message)
	case fiber.StatusConflict:
		return status.Error(codes.AlreadyExists, appErr.Message)
	case fiber.StatusUnprocessableEntity:
		return status.Error(codes.InvalidArgument, appErr.Message)
	case fiber.StatusServiceUnavailable:
		return status.Error(codes.Unavailable, appErr.Message)
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
