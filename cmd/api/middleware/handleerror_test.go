package middleware

import (
	"errors"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/dougefr/product-scrapper/cmd/api/common"
	"github.com/dougefr/product-scrapper/domain/businesserr"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandlerBusinessError(t *testing.T) {
	tableTest := map[businesserr.BusinessError]int{
		businesserr.InvalidProductURL:                                      422,
		businesserr.InvalidProductPrice:                                    422,
		businesserr.InvalidProductDescription:                              422,
		businesserr.InvalidProductImage:                                    422,
		businesserr.InvalidProductTitle:                                    422,
		businesserr.ScrapperNotImplemented:                                 422,
		businesserr.NewBusinessError("ErrFakeError", "this error is fake"): 400,
	}

	assert.Equal(t, len(businesserr.AllBusinessErrors), len(tableTest), "a business error is missing on the table test above")

	for be, statusCode := range tableTest {
		t.Run("test handle business error "+be.Code()+" - "+strconv.Itoa(statusCode), func(t *testing.T) {
			api := fiber.New()
			SetMiddlewares(api, nil, nil, NewHandlerError())
			api.Get("/", func(ctx *fiber.Ctx) error {
				return be
			})

			req := httptest.NewRequest("GET", "/", nil)
			resp, _ := api.Test(req, 10000)
			assert.Equal(t, statusCode, resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, "{\"message\":\""+be.Error()+"\",\"code\":\""+be.Code()+"\"}",
				string(body))
		})
	}
}

func TestHandlerCustomResponseError(t *testing.T) {
	tableTest := map[common.CustomResponseError]int{
		common.NewCustomResponseError(400, "message0"): 400,
		common.NewCustomResponseError(401, "message1"): 401,
		common.NewCustomResponseError(402, "message2"): 402,
		common.NewCustomResponseError(403, "message3"): 403,
		common.NewCustomResponseError(404, "message4"): 404,
		common.NewCustomResponseError(405, "message5"): 405,
		common.NewCustomResponseError(406, "message6"): 406,
		common.NewCustomResponseError(407, "message7"): 407,
	}

	for ce, statusCode := range tableTest {
		t.Run("test handle business error "+ce.Error()+" - "+strconv.Itoa(statusCode), func(t *testing.T) {
			api := fiber.New()
			SetMiddlewares(api, nil, nil, NewHandlerError())
			api.Get("/", func(ctx *fiber.Ctx) error {
				return ce
			})

			req := httptest.NewRequest("GET", "/", nil)
			resp, _ := api.Test(req, 10000)
			assert.Equal(t, statusCode, resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, "{\"message\":\""+ce.Message()+"\",\"code\":\"UnknownError\"}",
				string(body))
		})
	}
}

func TestHandlerInternalError(t *testing.T) {
	api := fiber.New()
	SetMiddlewares(api, nil, nil, NewHandlerError())
	api.Get("/", func(ctx *fiber.Ctx) error {
		return errors.New("fake error")
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := api.Test(req, 10000)
	assert.Equal(t, 500, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "{\"message\":\"internal error\",\"code\":\"UnknownError\"}",
		string(body))
}

func TestHandlerInternalNoError(t *testing.T) {
	api := fiber.New()
	SetMiddlewares(api, nil, nil, NewHandlerError())
	api.Get("/", func(ctx *fiber.Ctx) error {
		return nil
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := api.Test(req, 10000)
	assert.Equal(t, 200, resp.StatusCode)
}
