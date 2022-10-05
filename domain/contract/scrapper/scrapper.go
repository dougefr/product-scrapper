package scrapper

import (
	"context"

	"github.com/dougefr/product-scrapper/domain/entity"
)

type (
	// Scrapper interface que representa a infraestrutura de um scrapper
	Scrapper interface {
		ScrapProduct(ctx context.Context) (entity.Product, error)
	}

	// Factory interface que representa uma factory para Scrapper
	Factory interface {
		GetScrapper(ctx context.Context, url string) (Scrapper, error)
	}
)
