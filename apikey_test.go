package apikey

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber"
)

func Test_Skip(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?key=secret", nil)
	app := fiber.New()
	app.Use(New(Config{Key: "secret", Skip: func(c *fiber.Ctx) bool { return true }}))
	app.Get("/", func(c *fiber.Ctx) {
		c.Send("ok")
	})
	resp, _ := app.Test(req)

	if http.StatusOK != resp.StatusCode {
		t.Error("middleware should be skipped")
	}
}

func Test_MissingKeyInRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	app := fiber.New()
	app.Use(New(Config{Key: "secret"}))
	app.Get("/", func(c *fiber.Ctx) {
		c.Send("ok")
	})
	resp, _ := app.Test(req)

	if http.StatusUnauthorized != resp.StatusCode {
		t.Error("request should be blocked")
	}
}

func Test_KeyInQueryParams(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?key=secret", nil)
	app := fiber.New()
	app.Use(New(Config{Key: "secret"}))
	app.Get("/", func(c *fiber.Ctx) {
		c.Send("ok")
	})
	resp, _ := app.Test(req)

	if http.StatusOK != resp.StatusCode {
		t.Error("request should be allowed")
	}
}

func Test_KeyInHeaders(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("x-api-key", "secret")

	app := fiber.New()
	app.Use(New(Config{Key: "secret"}))
	app.Get("/", func(c *fiber.Ctx) {
		c.Send("ok")
	})
	resp, _ := app.Test(req)

	if http.StatusOK != resp.StatusCode {
		t.Error("request should be allowed")
	}
}

func Test_CustomValidatorFunc(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("x-api-key", "secret")

	app := fiber.New()
	app.Use(New(Config{
		Key: "secret",
		ValidatorFunc: func(c *fiber.Ctx, cfg Config) bool {
			return false
		},
	}))
	app.Get("/", func(c *fiber.Ctx) {
		c.Send("ok")
	})
	resp, _ := app.Test(req)

	if http.StatusUnauthorized != resp.StatusCode {
		t.Error("request should be blocked")
	}
}
