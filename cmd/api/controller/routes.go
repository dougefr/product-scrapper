package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// SetControllersRoutes adiciona as rotas na API
func SetControllersRoutes(
	app *fiber.App,
	api fiber.Router,
	healthCheck HealthCheck,
	productCtrl Product) {

	if healthCheck != nil {
		app.Get("/healthcheck", healthCheck.HealthCheck)
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	if productCtrl != nil {
		api.Post("/products", productCtrl.GetProduct)
	}
}
