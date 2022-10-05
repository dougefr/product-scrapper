package controller

import (
	"github.com/dougefr/product-scrapper/cmd/api/common"
	"github.com/dougefr/product-scrapper/domain/contract/logger"
	"github.com/dougefr/product-scrapper/domain/usecase"
	"github.com/dougefr/product-scrapper/infra/context"
	"github.com/gofiber/fiber/v2"
)

type (
	// Product controller para Product
	Product interface {
		GetProduct(ctx *fiber.Ctx) error
	}

	product struct {
		logger             logger.Logger
		scrapProductInfoUC usecase.ScrapProduct
	}

	scrapProductInfoBody struct {
		URL string `json:"url"`
	}

	scrapProductInfoResponse struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Image       string  `json:"image"`
		Price       float64 `json:"price"`
		URL         string  `json:"url"`
	}
)

func newScrapProductInfoResponse(
	title string,
	description string,
	image string,
	price float64,
	URL string,
) scrapProductInfoResponse {

	return scrapProductInfoResponse{
		Title:       title,
		Description: description,
		Image:       image,
		Price:       price,
		URL:         URL,
	}
}

// NewProduct cria um novo Product
func NewProduct(
	logger logger.Logger,
	scrapProductInfo usecase.ScrapProduct,
) Product {

	return product{
		logger:             logger,
		scrapProductInfoUC: scrapProductInfo,
	}
}

// GetProduct obtém os dados de um produto de acordo com a URL passada pelo body
// @Summary obtém os dados de um produto de acordo com a URL passada pelo body
// @Description obtém os dados de um produto de acordo com a URL passada pelo body
// @Accept json
// @Produce json
// @Param ingredient body scrapProductInfoBody true "dados do website do produto"
// @Success 200 {object} scrapProductInfoResponse
// @Failure 422 {object} common.ResponseError "businesserr.InvalidProductTitle,businesserr.InvalidProductImage,businesserr.InvalidProductPrice,businesserr.InvalidProductDescription,businesserr.InvalidProductURL"
// @Router       /products [post]
func (p product) GetProduct(ctx *fiber.Ctx) (err error) {
	body := scrapProductInfoBody{}
	err = ctx.BodyParser(&body)
	if err != nil {
		p.logger.Error(context.ConvertFiberCtxToCtx(ctx), "error when parsing body", logger.Body{
			"err": err,
		})
		return common.NewCustomResponseError(422, "unable to parse payload as JSON")
	}

	output, err := p.scrapProductInfoUC.Execute(context.ConvertFiberCtxToCtx(ctx),
		usecase.NewScrapProductInput(body.URL))
	if err != nil {
		p.logger.Error(context.ConvertFiberCtxToCtx(ctx), "error when executing findAllIngredient usecase", logger.Body{
			"err": err,
		})
		return
	}

	return ctx.JSON(newScrapProductInfoResponse(
		output.Title,
		output.Description,
		output.Image,
		output.Price,
		output.URL,
	))
}
