package middleware

import (
	"time"

	"github.com/dougefr/product-scrapper/domain/contract/logger"
	"github.com/dougefr/product-scrapper/infra/context"
	"github.com/gofiber/fiber/v2"
)

type (
	// Log interface do middleware responsável por logar a entrada e saída de uma requisição
	Log interface {
		LogRequest(ctx *fiber.Ctx) error
	}

	log struct {
		logger logger.Logger
	}
)

// NewLog cria um novo Log
func NewLog(logger logger.Logger) Log {
	return log{
		logger: logger,
	}
}

// LogRequest realiza o log da requisição
func (l log) LogRequest(fiberctx *fiber.Ctx) (err error) {
	ctx := context.ConvertFiberCtxToCtx(fiberctx)
	l.logger.Info(ctx, "begin")
	start := time.Now()
	err = fiberctx.Next()
	duration := time.Since(start)

	l.logger.Info(ctx, "end", logger.Body{
		"runtime": duration.Milliseconds(),
	})

	return
}
