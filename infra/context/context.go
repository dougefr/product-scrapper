package context

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type key int

const (
	// RequestID ID atribuido a requisição
	RequestID key = iota
	// Path da requisição
	Path key = iota
	// Method da requisição
	Method key = iota
)

// ConvertFiberCtxToCtx converte o context do gearbox em um context padrão do golang
func ConvertFiberCtxToCtx(fiberctx *fiber.Ctx) context.Context {
	ctx := context.WithValue(fiberctx.Context(), RequestID, fiberctx.GetRespHeader("request-id"))
	ctx = context.WithValue(ctx, Path, string(fiberctx.Context().Path()))
	ctx = context.WithValue(ctx, Method, string(fiberctx.Context().Method()))

	return ctx
}
