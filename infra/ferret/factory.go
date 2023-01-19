package ferret

import (
	"context"
	"strings"

	"github.com/dougefr/product-scrapper/domain/businesserr"
	"github.com/dougefr/product-scrapper/domain/contract/env"
	logger2 "github.com/dougefr/product-scrapper/domain/contract/logger"
	"github.com/dougefr/product-scrapper/domain/contract/scrapper"
)

type factory struct {
	logger logger2.Logger
	env    env.Env
}

// NewFactory cria uma nova scrapper.Factory
func NewFactory(
	logger logger2.Logger,
	env env.Env,
) scrapper.Factory {

	return factory{
		logger: logger,
		env:    env,
	}
}

func (f factory) GetScrapper(ctx context.Context, url string) (scrapper.Scrapper, error) {
	if strings.Contains(url, "www.magazineluiza.com.br") {
		return newMagazineLuiza(ctx, f.logger, url, f.env)
	}

	if strings.Contains(url, "www.zattini.com.br") {
		return newZattini(ctx, f.logger, url, f.env)
	}

	if strings.Contains(url, "www.amazon.com") {
		return newAmazon(ctx, f.logger, url, f.env)
	}

	if strings.Contains(url, "www.saraiva.com.br") {
		return newSaraiva(ctx, f.logger, url, f.env)
	}

	return nil, businesserr.ScrapperNotImplemented
}
