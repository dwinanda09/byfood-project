package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const RequestIDKey = "request_id"

func DefaultMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var requestID string
			
			// Check if request already has X-Request-ID header
			if headerRequestID := c.Request().Header.Get("X-Request-ID"); headerRequestID != "" {
				requestID = headerRequestID
			} else {
				// Generate new UUID for request tracking
				requestID = uuid.New().String()
			}
			
			// Store in Echo context
			c.Set(RequestIDKey, requestID)
			
			// Add to Go context for deeper layer access
			ctx := context.WithValue(c.Request().Context(), RequestIDKey, requestID)
			c.SetRequest(c.Request().WithContext(ctx))
			
			// Add to response headers for client tracking
			c.Response().Header().Set("X-Request-ID", requestID)
			
			return next(c)
		}
	}
}