package apikey

import (
	"net/http"
	"os"

	"github.com/gofiber/fiber"
)

const (
	// DefaultKeyIdentifier is the default api key identifier
	DefaultKeyIdentifier string = "key"
	// DefaultHeaderKeyIdentifier is the default api key identifier in request headers
	DefaultHeaderKeyIdentifier string = "x-api-key"
)

// Config is the configuration object for this middleware
type Config struct {
	// Skip this middleware
	Skip func(*fiber.Ctx) bool
	// Key is the api key
	Key string
	// ValidatorFunc is the function to validates the request
	ValidatorFunc func(*fiber.Ctx, Config) bool
}

// DefaultValidatorFunc is the default validator function used to validates the request
// This validator try to look for api key in:
// - URL's query params
// - Request headers
// If api key not found in both of those location, return false
func DefaultValidatorFunc(c *fiber.Ctx, cfg Config) bool {
	queryKey := c.Query(DefaultKeyIdentifier)
	headerKey := c.Get(DefaultHeaderKeyIdentifier)
	if queryKey == "" && headerKey == "" {
		return false
	}
	if queryKey != "" && queryKey == cfg.Key {
		return true
	}
	if headerKey != "" && headerKey == cfg.Key {
		return true
	}
	return false
}

var defaultConfig = Config{
	Key:           os.Getenv("API_KEY"),
	ValidatorFunc: DefaultValidatorFunc,
}

// New returns the middleware function
func New(config ...Config) func(*fiber.Ctx) {
	var cfg Config

	if len(config) == 0 {
		cfg = defaultConfig
	} else {
		cfg = config[0]
		if cfg.ValidatorFunc == nil {
			cfg.ValidatorFunc = DefaultValidatorFunc
		}
	}

	return func(c *fiber.Ctx) {
		if cfg.Skip != nil && cfg.Skip(c) {
			c.Next()
			return
		}
		pass := cfg.ValidatorFunc(c, cfg)
		if !pass {
			c.SendStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
