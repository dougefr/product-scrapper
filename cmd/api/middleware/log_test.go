package middleware

import (
	"context"
	"net/http/httptest"
	"testing"

	logger2 "github.com/dougefr/product-scrapper/domain/contract/logger"
	mock_logger "github.com/dougefr/product-scrapper/domain/contract/mock/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestShouldLogRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any(), "begin", gomock.Any())
	logger.EXPECT().Info(gomock.Any(), "end", gomock.Any()).Do(func(
		ctx context.Context,
		message string,
		valueMaps ...logger2.Body) {

		assert.GreaterOrEqual(t, int64(0), valueMaps[0]["runtime"])
	})

	mid := NewLog(logger)

	api := fiber.New()
	SetMiddlewares(api, nil, mid, nil)
	api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := api.Test(req, 10000)
	assert.Equal(t, 200, resp.StatusCode)
}
