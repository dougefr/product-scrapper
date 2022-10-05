package middleware

import (
	"github.com/dougefr/product-scrapper/cmd/api/common"
	"github.com/dougefr/product-scrapper/domain/businesserr"
	"github.com/gofiber/fiber/v2"
)

type (
	// ErrorHandler interface do middleware que trata erros retornados pelos controllers
	ErrorHandler interface {
		HandleError(ctx *fiber.Ctx) error
	}

	errorHandler struct {
	}
)

// NewHandlerError cria um novo ErrorHandler
func NewHandlerError() ErrorHandler {
	return errorHandler{}
}

// HandleError atribui um request ID a requisição
func (h errorHandler) HandleError(ctx *fiber.Ctx) error {
	err := ctx.Next()
	if err == nil {
		return nil
	}

	if be, ok := err.(businesserr.BusinessError); ok {
		switch be {

		case
			businesserr.InvalidProductTitle,
			businesserr.InvalidProductImage,
			businesserr.InvalidProductPrice,
			businesserr.InvalidProductDescription,
			businesserr.InvalidProductURL,
			businesserr.ScrapperNotImplemented:
			return SendBusinessErrorJson(ctx, 422, be)

		default:
			return SendBusinessErrorJson(ctx, 400, be)
		}
	}

	if re, ok := err.(common.CustomResponseError); ok {
		return SendErrorJson(ctx, re.StatusCode(), re.Message())
	}

	return SendErrorJson(ctx, 500, "internal error")
}

// SendBusinessErrorJson envia um json com o código do erro de negócio
func SendBusinessErrorJson(fiberctx *fiber.Ctx, status int, bs businesserr.BusinessError) error {
	fiberctx.Status(status)
	body := common.ResponseError{
		Code:    bs.Code(),
		Message: bs.Error(),
	}

	return fiberctx.JSON(body)
}

// SendErrorJson envia um json com uma mensagem de erro genérico
func SendErrorJson(fiberctx *fiber.Ctx, status int, msg string) error {
	fiberctx.Status(status)
	body := common.ResponseError{
		Code:    "UnknownError",
		Message: msg,
	}

	return fiberctx.JSON(body)
}
