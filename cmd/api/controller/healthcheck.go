package controller

import (
	"github.com/dougefr/product-scrapper/domain/contract/logger"
	"github.com/gofiber/fiber/v2"
)

type (
	// HealthCheck controller para HealthCheck
	HealthCheck interface {
		HealthCheck(ctx *fiber.Ctx) error
	}

	healthCheck struct {
		logger logger.Logger
	}
)

// NewHealthCheck cria um novo HealthCheck
func NewHealthCheck(
	logger logger.Logger,
) HealthCheck {

	return healthCheck{
		logger: logger,
	}
}

func (h healthCheck) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.SendString("online")
}
