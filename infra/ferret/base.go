package ferret

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/drivers/http"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/dougefr/product-scrapper/domain/contract/env"
	logger2 "github.com/dougefr/product-scrapper/domain/contract/logger"
	"github.com/dougefr/product-scrapper/domain/entity"
)

type (
	base struct {
		logger  logger2.Logger
		program *runtime.Program
		url     string
		env     env.Env
	}

	result struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Image       string `json:"image"`
		Price       string `json:"price"`
	}
)

func newBase(
	ctx context.Context,
	logger logger2.Logger,
	url string,
	env env.Env,
	querySelectorTitle string,
	querySelectorDescription string,
	querySelectorImage string,
	querySelectorPrice string,
) (b *base, err error) {
	logger.Info(ctx, "initializing scrapper", logger2.Body{
		"url":                      url,
		"querySelectorDescription": querySelectorDescription,
		"querySelectorTitle":       querySelectorTitle,
		"querySelectorImage":       querySelectorImage,
		"querySelectorPrice":       querySelectorPrice,
	})

	query := fmt.Sprintf(`
		LET doc = DOCUMENT('%s', {
			driver: "cdp", 
			userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.87 Safari/537.36"
		})
		LET title = ELEMENT(doc, '%s')
		LET description = ELEMENT(doc, '%s')
		LET image = ELEMENT(doc, '%s')
		LET price = ELEMENT(doc, '%s')
		
		RETURN {
			title: TRIM(title?.innerText),
			description: TRIM(description?.innerText),
			image: image?.attributes?.src,
			price: TRIM(price?.innerText)
		}
	`,
		url,
		querySelectorTitle,
		querySelectorDescription,
		querySelectorImage,
		querySelectorPrice)

	logger.Info(ctx, "compiling scrapper query", logger2.Body{
		"url":                      url,
		"querySelectorDescription": querySelectorDescription,
		"querySelectorTitle":       querySelectorTitle,
		"querySelectorImage":       querySelectorImage,
		"querySelectorPrice":       querySelectorPrice,
		"query":                    query,
	})

	comp := compiler.New()
	program, err := comp.Compile(query)
	if err != nil {
		logger.Error(ctx, "error when compiling scrapper query", logger2.Body{
			"url":                      url,
			"querySelectorDescription": querySelectorDescription,
			"querySelectorTitle":       querySelectorTitle,
			"querySelectorImage":       querySelectorImage,
			"querySelectorPrice":       querySelectorPrice,
			"query":                    query,
			"err":                      err,
		})

		err = fmt.Errorf("error when compiling scrapper query: %w", err)
		return
	}

	b = &base{
		logger:  logger,
		program: program,
		url:     url,
		env:     env,
	}

	logger.Info(ctx, "scrapper initialized with success", logger2.Body{
		"url":                      url,
		"querySelectorDescription": querySelectorDescription,
		"querySelectorTitle":       querySelectorTitle,
		"querySelectorImage":       querySelectorImage,
		"querySelectorPrice":       querySelectorPrice,
		"query":                    query,
	})

	return
}

func (b base) ScrapProduct(ctx context.Context) (product entity.Product, err error) {
	b.logger.Info(ctx, "scrap started", logger2.Body{
		"url": b.url,
	})

	ctx = drivers.WithContext(ctx, cdp.NewDriver(cdp.WithAddress(b.env.BrowserAddr())))
	ctx = drivers.WithContext(ctx, http.NewDriver(), drivers.AsDefault())

	out, err := b.program.Run(ctx)
	if err != nil {
		b.logger.Error(ctx, "error when scrap website", logger2.Body{
			"url": b.url,
			"err": err,
		})

		err = fmt.Errorf("error when scrap website: %w", err)
		return
	}

	var res result
	err = json.Unmarshal(out, &res)
	if err != nil {
		b.logger.Error(ctx, "error when unmarshalling scrap result", logger2.Body{
			"url": b.url,
			"err": err,
		})

		err = fmt.Errorf("error when unmarshalling scrap result: %w", err)
		return
	}

	strPrice := res.Price
	strPrice = strings.ReplaceAll(strPrice, " ", "")
	strPrice = strings.ReplaceAll(strPrice, "R$", "")
	strPrice = strings.ReplaceAll(strPrice, "&nbsp;", "")
	strPrice = strings.ReplaceAll(strPrice, ".", "")
	strPrice = strings.ReplaceAll(strPrice, ",", ".")

	price, err := strconv.ParseFloat(strPrice, 32)
	if err != nil {
		b.logger.Error(ctx, "error when converting price from string to float", logger2.Body{
			"url": b.url,
			"err": err,
		})

		err = nil // scrap will not be stopped if an error occur in getting price
	}

	product = entity.NewProduct(
		res.Title,
		res.Description,
		res.Image,
		price,
		b.url,
	)

	b.logger.Info(ctx, "scrap finished", logger2.Body{
		"url":         b.url,
		"title":       product.Title,
		"description": product.Description,
		"image":       product.Image,
		"price":       product.Price,
	})

	return
}
