package middleware

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/random"
)

type (
	// RequestIDConfig defines the config for RequestID middleware.
	RequestIDConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Generator defines a function to generate an ID.
		// Optional. Default value random.String(32).
		Generator func() string

		// GenerateFromContext defines a function to generate an ID.
		// Diff with Generator function, user can get metadata from echo.Context,
		// like cookie, header, params etc.
		GenerateFromContext func(c echo.Context) string
	}
)

var (
	// DefaultRequestIDConfig is the default RequestID middleware config.
	DefaultRequestIDConfig = RequestIDConfig{
		Skipper:             DefaultSkipper,
		Generator:           generator,
		GenerateFromContext: generateFromContext,
	}
)

// RequestID returns a X-Request-ID middleware.
func RequestID() echo.MiddlewareFunc {
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

// RequestIDWithConfig returns a X-Request-ID middleware with config.
func RequestIDWithConfig(config RequestIDConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRequestIDConfig.Skipper
	}
	if config.Generator == nil {
		config.Generator = generator
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				if config.GenerateFromContext != nil {
					rid = config.GenerateFromContext(c)
				} else {
					rid = config.Generator()
				}
			}
			res.Header().Set(echo.HeaderXRequestID, rid)

			return next(c)
		}
	}
}

func generator() string {
	return random.String(32)
}

func generateFromContext(c echo.Context) string {
	return generator()
}
