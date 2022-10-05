package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/dougefr/product-scrapper/domain/businesserr"
	cache2 "github.com/dougefr/product-scrapper/domain/contract/cache"
	logger2 "github.com/dougefr/product-scrapper/domain/contract/logger"
	"github.com/dougefr/product-scrapper/domain/contract/scrapper"
	"github.com/dougefr/product-scrapper/domain/entity"
)

type (
	// ScrapProduct caso de uso responsável por fazer scrap dos produtos
	ScrapProduct interface {
		Execute(ctx context.Context, in ScrapProductInput) (ScrapProductOutput, error)
	}

	// ScrapProductInput contrato de entrada do usecase ScrapProduct
	ScrapProductInput struct {
		URL string
	}

	// ScrapProductOutput contrato de saída do usecase ScrapProduct
	ScrapProductOutput struct {
		Title       string
		Description string
		Image       string
		Price       float64
		URL         string
	}

	scrapProduct struct {
		logger          logger2.Logger
		cache           cache2.Cache[entity.Product]
		scrapperFactory scrapper.Factory
	}
)

// NewScrapProductInput cria um novo ScrapProductInput
func NewScrapProductInput(
	URL string,
) ScrapProductInput {

	return ScrapProductInput{
		URL: URL,
	}
}

// NewScrapProductOutput croa um novo ScrapProductOutput
func NewScrapProductOutput(
	title string,
	description string,
	image string,
	price float64,
	URL string,
) ScrapProductOutput {

	return ScrapProductOutput{
		Title:       title,
		Description: description,
		Image:       image,
		Price:       price,
		URL:         URL,
	}
}

// NewScrapProduct cria um novo ScrapProduct
func NewScrapProduct(
	logger logger2.Logger,
	cache cache2.Cache[entity.Product],
	scrapperFactory scrapper.Factory,
) ScrapProduct {

	return scrapProduct{
		logger:          logger,
		cache:           cache,
		scrapperFactory: scrapperFactory,
	}
}

func (s scrapProduct) Execute(ctx context.Context, in ScrapProductInput) (out ScrapProductOutput, err error) {
	s.logger.Info(ctx, "looking for data in cache", logger2.Body{
		"url": in.URL,
	})

	cachedData, err := s.cache.Get(ctx, in.URL)
	if err != nil {
		s.logger.Warn(ctx, "error when getting data from cache", logger2.Body{
			"err": err,
			"url": in.URL,
		})

		err = nil // o scrap não será interrompido se ocorrer um erro no cache
	}

	if cachedData != nil {
		if be := cachedData.IsValid(); be != nil {
			s.logger.Warn(ctx, "scrapper returned an invalid product info", logger2.Body{
				"err": be,
				"url": in.URL,
			})

			err = nil // o scrap não será interrompido se ocorrer um erro no cache
		} else {
			out = NewScrapProductOutput(
				cachedData.Title,
				cachedData.Description,
				cachedData.Image,
				cachedData.Price,
				cachedData.URL,
			)

			s.logger.Info(ctx, "scrap product usecase finished - returning data from cache", logger2.Body{
				"url":         out.URL,
				"title":       out.Title,
				"description": out.Description,
				"image":       out.Image,
				"price":       out.Price,
			})

			return
		}
	}

	s.logger.Info(ctx, "getting data from scrapper", logger2.Body{
		"url": in.URL,
	})

	sc, err := s.scrapperFactory.GetScrapper(ctx, in.URL)
	if err != nil {
		s.logger.Error(ctx, "error when getting scrapper from factory", logger2.Body{
			"err": err,
			"url": in.URL,
		})

		if be, ok := err.(businesserr.BusinessError); ok {
			err = be
			return
		}

		err = fmt.Errorf("error when getting scrapper from factory: %w", err)
		return
	}

	product, err := sc.ScrapProduct(ctx)
	if err != nil {
		s.logger.Error(ctx, "error when getting product from scrapper", logger2.Body{
			"err": err,
			"url": in.URL,
		})

		err = fmt.Errorf("error when getting product from scrapper: %w", err)
		return
	}

	if be := product.IsValid(); be != nil {
		s.logger.Error(ctx, "scrapper returned an invalid product info", logger2.Body{
			"err": be,
			"url": in.URL,
		})

		err = be
		return
	}

	err = s.cache.Set(ctx, in.URL, product, 1*time.Hour)
	if err != nil {
		s.logger.Warn(ctx, "error when setting scrapper result in cache", logger2.Body{
			"err": err,
			"url": in.URL,
		})

		err = nil // o scrap não será interrompido se ocorrer um erro no cache
	}

	out = NewScrapProductOutput(
		product.Title,
		product.Description,
		product.Image,
		product.Price,
		product.URL,
	)

	s.logger.Info(ctx, "scrap product usecase finished - returning data from scrapper", logger2.Body{
		"url":         out.URL,
		"title":       out.Title,
		"description": out.Description,
		"image":       out.Image,
		"price":       out.Price,
	})

	return
}
