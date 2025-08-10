package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type SecurityConfig struct {
	APIKeyHeader    string
	AllowedAPIKeys  []string
	RateLimitRPS    int
	RateLimitBurst  int
	EnableAPIKey    bool
	EnableRateLimit bool
	TrustedProxies  []string
	MaxRequestSize  string
}

type SecurityMiddleware struct {
	config      SecurityConfig
	logger      *zap.Logger
	rateLimiter *rate.Limiter
}

func NewSecurityMiddleware(config SecurityConfig, logger *zap.Logger) *SecurityMiddleware {
	var limiter *rate.Limiter
	if config.EnableRateLimit {
		limiter = rate.NewLimiter(rate.Limit(config.RateLimitRPS), config.RateLimitBurst)
	}

	return &SecurityMiddleware{
		config:      config,
		logger:      logger,
		rateLimiter: limiter,
	}
}

// SecurityHeaders adds essential security headers
func (sm *SecurityMiddleware) SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Security headers
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")

			// HSTS header for HTTPS
			if c.Request().TLS != nil {
				c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}

			// Remove server information
			c.Response().Header().Set("Server", "")

			return next(c)
		}
	}
}

// APIKeyAuth validates API key authentication
func (sm *SecurityMiddleware) APIKeyAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !sm.config.EnableAPIKey {
				return next(c)
			}

			// Skip authentication for health check endpoint
			if c.Request().URL.Path == "/health" || c.Request().URL.Path == "/metrics" {
				return next(c)
			}

			apiKey := c.Request().Header.Get(sm.config.APIKeyHeader)
			if apiKey == "" {
				sm.logger.Warn("Missing API key",
					zap.String("path", c.Request().URL.Path),
					zap.String("remote_addr", c.Request().RemoteAddr),
				)
				return echo.NewHTTPError(http.StatusUnauthorized, "API key required")
			}

			// Validate API key
			validKey := false
			for _, validAPIKey := range sm.config.AllowedAPIKeys {
				if apiKey == validAPIKey {
					validKey = true
					break
				}
			}

			if !validKey {
				sm.logger.Warn("Invalid API key",
					zap.String("path", c.Request().URL.Path),
					zap.String("remote_addr", c.Request().RemoteAddr),
					zap.String("api_key_prefix", fmt.Sprintf("%s...", apiKey[:min(len(apiKey), 8)])),
				)
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API key")
			}

			return next(c)
		}
	}
}

// RateLimiter implements rate limiting
func (sm *SecurityMiddleware) RateLimiter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !sm.config.EnableRateLimit || sm.rateLimiter == nil {
				return next(c)
			}

			// Skip rate limiting for health check
			if c.Request().URL.Path == "/health" {
				return next(c)
			}

			if !sm.rateLimiter.Allow() {
				sm.logger.Warn("Rate limit exceeded",
					zap.String("path", c.Request().URL.Path),
					zap.String("remote_addr", c.Request().RemoteAddr),
					zap.String("user_agent", c.Request().UserAgent()),
				)
				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			}

			return next(c)
		}
	}
}

// RequestValidator validates request structure and content
func (sm *SecurityMiddleware) RequestValidator() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Validate content type for POST/PUT requests
			if c.Request().Method == http.MethodPost || c.Request().Method == http.MethodPut {
				contentType := c.Request().Header.Get("Content-Type")
				if !strings.Contains(contentType, "application/json") {
					sm.logger.Warn("Invalid content type",
						zap.String("content_type", contentType),
						zap.String("path", c.Request().URL.Path),
					)
					return echo.NewHTTPError(http.StatusBadRequest, "Content-Type must be application/json")
				}
			}

			// Validate request size
			if c.Request().ContentLength > 0 {
				maxSize := int64(1024 * 1024) // 1MB default
				if c.Request().ContentLength > maxSize {
					sm.logger.Warn("Request size too large",
						zap.Int64("content_length", c.Request().ContentLength),
						zap.Int64("max_size", maxSize),
						zap.String("path", c.Request().URL.Path),
					)
					return echo.NewHTTPError(http.StatusRequestEntityTooLarge, "Request body too large")
				}
			}

			return next(c)
		}
	}
}

// CORS configuration for production
func (sm *SecurityMiddleware) CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://localhost:3000",
			"https://*.byfood.com",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-API-Key",
		},
		ExposeHeaders: []string{
			"X-Request-ID",
		},
		AllowCredentials: false,
		MaxAge:           86400, // 24 hours
	})
}

// min helper function for API key prefix logging
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
