package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestResponseMustHaveARequestID(t *testing.T) {
	mid := NewRequestId()

	api := fiber.New()
	SetMiddlewares(api, mid, nil, nil)
	api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := api.Test(req, 1000)
	assert.NotNil(t, resp.Header.Get("x-request-id"))
	assert.NotEqual(t, "", resp.Header.Get("x-request-id"))
}

func TestResponseMustHaveTheRequestIDFromRequest(t *testing.T) {
	mid := NewRequestId()

	api := fiber.New()
	SetMiddlewares(api, mid, nil, nil)
	api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("x-request-id", "123456")
	resp, _ := api.Test(req, 1000)
	assert.Equal(t, "123456", resp.Header.Get("x-request-id"))
}
