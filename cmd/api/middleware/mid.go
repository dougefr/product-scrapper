package middleware

import "github.com/gofiber/fiber/v2"

// SetMiddlewares atribui middlewares ao app
func SetMiddlewares(
	app *fiber.App,
	requestId RequestId,
	log Log,
	errorHandler ErrorHandler) {

	if requestId != nil {
		app.Use(requestId.AddRequestId)
	}

	if log != nil {
		app.Use(log.LogRequest)
	}

	if errorHandler != nil {
		app.Use(errorHandler.HandleError)
	}
}
