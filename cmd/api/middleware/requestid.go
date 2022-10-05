package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const xRequestID = "x-request-id"

type (
	// RequestId interface do middleware que atribui um ID a requisicão
	RequestId interface {
		AddRequestId(ctx *fiber.Ctx) error
	}

	requestId struct {
	}
)

// NewRequestId cria um novo RequestId
func NewRequestId() RequestId {
	return requestId{}
}

// AddRequestId atribui um request ID a requisição
func (r requestId) AddRequestId(ctx *fiber.Ctx) (err error) {
	var requestID string

	if ctx.Get(xRequestID) != "" {
		requestID = ctx.Get(xRequestID)
	} else {
		requestID = uuid.NewString()
	}

	ctx.Set(xRequestID, requestID)
	ctx.Set("request-id", requestID)
	return ctx.Next()
}
