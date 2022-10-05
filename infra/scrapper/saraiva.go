package scrapper

import (
	"context"

	"github.com/dougefr/product-scrapper/domain/contract/env"
	logger2 "github.com/dougefr/product-scrapper/domain/contract/logger"
)

type saraiva struct {
	*base
}

func newSaraiva(
	ctx context.Context,
	logger logger2.Logger,
	url string,
	env env.Env,
) (s *saraiva, err error) {

	b, err := newBase(
		ctx,
		logger,
		url,
		env,
		"h1.title",
		"div[id=\"descricao\"]",
		"#product-1-item > img",
		"p.price-destaque")

	s = &saraiva{
		base: b,
	}

	return
}
