package middleware

import (
	"net/http"

	"byfood-library/internal/domain/entities"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ErrorHandler struct {
	logger *zap.Logger
}

func NewErrorHandler(logger *zap.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
	Code      int    `json:"code"`
}

// CustomHTTPErrorHandler handles errors in a centralized way
func (eh *ErrorHandler) CustomHTTPErrorHandler(err error, c echo.Context) {
	requestID := GetRequestID(c)
	
	var code int
	var message string
	var errorType string

	// Handle different error types
	switch err {
	case entities.ErrBookNotFound:
		code = http.StatusNotFound
		errorType = "NOT_FOUND"
		message = "The requested book could not be found"
	case entities.ErrInvalidTitle:
		code = http.StatusBadRequest
		errorType = "VALIDATION_ERROR"
		message = "Book title is required and cannot be empty"
	case entities.ErrInvalidAuthor:
		code = http.StatusBadRequest
		errorType = "VALIDATION_ERROR"
		message = "Book author is required and cannot be empty"
	case entities.ErrInvalidYear:
		code = http.StatusBadRequest
		errorType = "VALIDATION_ERROR"
		message = "Publication year must be between 1000 and 2034"
	case entities.ErrInvalidUUID:
		code = http.StatusBadRequest
		errorType = "VALIDATION_ERROR"
		message = "Invalid UUID format provided"
	case entities.ErrDatabaseError:
		code = http.StatusInternalServerError
		errorType = "INTERNAL_ERROR"
		message = "A database error occurred. Please try again later"
	default:
		// Handle Echo HTTP errors
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			errorType = "HTTP_ERROR"
			message = he.Message.(string)
		} else {
			// Unknown error
			code = http.StatusInternalServerError
			errorType = "INTERNAL_ERROR"
			message = "An unexpected error occurred. Please try again later"
		}
	}

	// Log the error with context
	eh.logger.Error("HTTP error occurred",
		zap.String("request_id", requestID),
		zap.Int("status_code", code),
		zap.String("error_type", errorType),
		zap.String("path", c.Request().URL.Path),
		zap.String("method", c.Request().Method),
		zap.Error(err),
	)

	// Send structured error response
	errorResponse := ErrorResponse{
		Error:     errorType,
		Message:   message,
		RequestID: requestID,
		Code:      code,
	}

	// If response has already been committed, just log the error
	if c.Response().Committed {
		eh.logger.Error("Response already committed, cannot send error response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		return
	}

	// Send the error response
	if err := c.JSON(code, errorResponse); err != nil {
		eh.logger.Error("Failed to send error response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
	}
}

// Recovery middleware for panic handling
func (eh *ErrorHandler) Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					requestID := GetRequestID(c)
					eh.logger.Error("Panic recovered",
						zap.String("request_id", requestID),
						zap.Any("panic", r),
						zap.String("path", c.Request().URL.Path),
						zap.String("method", c.Request().Method),
					)

					err := &echo.HTTPError{
						Code:    http.StatusInternalServerError,
						Message: "Internal server error occurred",
					}
					eh.CustomHTTPErrorHandler(err, c)
				}
			}()
			return next(c)
		}
	}
}