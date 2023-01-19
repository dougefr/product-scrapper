package main

import (
	"context"

	"github.com/dougefr/product-scrapper/cmd/api/controller"
	"github.com/dougefr/product-scrapper/cmd/api/middleware"
	"github.com/dougefr/product-scrapper/domain/contract/logger"
	"github.com/dougefr/product-scrapper/domain/entity"
	"github.com/dougefr/product-scrapper/domain/usecase"
	"github.com/dougefr/product-scrapper/infra/osenv"
	"github.com/dougefr/product-scrapper/infra/ferret"
	"github.com/dougefr/product-scrapper/infra/redis"
	"github.com/dougefr/product-scrapper/infra/zap"
	"github.com/gofiber/fiber/v2"

	_ "github.com/dougefr/product-scrapper/docs"
)

// @title Product Scrapper API
// @description Product Scrapper API
// @BasePath /product-scrapper-api/v1
func main() {
	// Infras
	log := zap.NewLogger()
	env := osenv.NewEnv()
	redisClient := redis.NewRedis[entity.Product](env, log)
	scrapperFactory := ferret.NewFactory(log, env)

	// Usecases
	scrapProductInfoUC := usecase.NewScrapProduct(log, redisClient, scrapperFactory)

	// Controllers
	healthCheckCtrl := controller.NewHealthCheck(log)
	productCtrl := controller.NewProduct(log, scrapProductInfoUC)

	// Middlewares
	requestIdMid := middleware.NewRequestId()
	logMid := middleware.NewLog(log)
	errorHandler := middleware.NewHandlerError()

	// App
	app := fiber.New()
	api := app.Group("/product-scrapper-api/v1")
	middleware.SetMiddlewares(app, requestIdMid, logMid, errorHandler)
	controller.SetControllersRoutes(app, api, healthCheckCtrl, productCtrl)

	// Start service
	err := app.Listen(":8888")
	if err != nil {
		log.Fatal(context.Background(), "fatal error creating http server", logger.Body{
			"err": err,
		})
	}
}
