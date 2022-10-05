package controller

import (
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dougefr/product-scrapper/cmd/api/middleware"
	mock_logger "github.com/dougefr/product-scrapper/domain/contract/mock/logger"
	"github.com/dougefr/product-scrapper/domain/usecase"
	mock_usecase "github.com/dougefr/product-scrapper/domain/usecase/mock"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPostProductsShouldReturn422WhenBodyIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())

	uc := mock_usecase.NewMockScrapProduct(ctrl)

	c := NewProduct(logger, uc)

	app := fiber.New()
	middleware.SetMiddlewares(app, nil, nil, middleware.NewHandlerError())
	SetControllersRoutes(app, app, nil, c)

	req := httptest.NewRequest("POST", "/products", strings.NewReader("invalid!"))
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req, 1)
	defer resp.Body.Close()

	assert.Equal(t, 422, resp.StatusCode)
}

func TestPostProductsShouldReturn500WhenUcResultsAnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())

	uc := mock_usecase.NewMockScrapProduct(ctrl)
	uc.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(usecase.ScrapProductOutput{}, errors.New("fake error"))

	c := NewProduct(logger, uc)

	app := fiber.New()
	middleware.SetMiddlewares(app, nil, nil, middleware.NewHandlerError())
	SetControllersRoutes(app, app, nil, c)

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`
		{
		  "url": "https://www.amazon1.com.br/Colch%C3%A3o-Ensacadas-Espuma-Viscoel%C3%A1stica-Ortop%C3%A9dico/dp/B08LH5L7TK/ref=cm_cr_arp_d_product_top?ie=UTF8&th=1"
		}`))
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req, 1)
	defer resp.Body.Close()

	assert.Equal(t, 500, resp.StatusCode)
}

func TestPostProductsShouldReturn200WhenUcWasExecutedWithSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)

	uc := mock_usecase.NewMockScrapProduct(ctrl)
	uc.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(usecase.NewScrapProductOutput(
		"title",
		"description",
		"image",
		10.1,
		"url",
	), nil)

	c := NewProduct(logger, uc)

	app := fiber.New()
	middleware.SetMiddlewares(app, nil, nil, middleware.NewHandlerError())
	SetControllersRoutes(app, app, nil, c)

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`
		{
		  "url": "https://www.amazon1.com.br/Colch%C3%A3o-Ensacadas-Espuma-Viscoel%C3%A1stica-Ortop%C3%A9dico/dp/B08LH5L7TK/ref=cm_cr_arp_d_product_top?ie=UTF8&th=1"
		}`))
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req, 1)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var expected = "{\"title\":\"title\",\"description\":\"description\",\"image\":\"image\",\"price\":10.1,\"url\":\"url\"}"
	assert.Equal(t, expected, string(body))
}
