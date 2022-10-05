package controller

import (
	"net/http/httptest"
	"testing"

	"github.com/dougefr/product-scrapper/cmd/api/middleware"
	mock_logger "github.com/dougefr/product-scrapper/domain/contract/mock/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckShouldReturn200WhenDatabaseIsFine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)

	c := NewHealthCheck(logger)

	api := fiber.New()
	middleware.SetMiddlewares(api, nil, nil, middleware.NewHandlerError())
	SetControllersRoutes(api, api, c, nil)

	req := httptest.NewRequest("GET", "/healthcheck", nil)
	resp, _ := api.Test(req, 10000)

	assert.Equal(t, 200, resp.StatusCode)
}
