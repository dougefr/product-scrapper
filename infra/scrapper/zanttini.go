package scrapper

import (
	"context"

	"github.com/dougefr/product-scrapper/domain/contract/env"
	logger2 "github.com/dougefr/product-scrapper/domain/contract/logger"
)

type zattini struct {
	*base
}

func newZattini(
	ctx context.Context,
	logger logger2.Logger,
	url string,
	env env.Env,
) (z *zattini, err error) {

	b, err := newBase(
		ctx,
		logger,
		url,
		env,
		"h1[data-productname]",
		"p[itemprop=\"description\"]",
		"section.photo > figure > img",
		"div.default-price > * > strong")

	z = &zattini{
		base: b,
	}

	return
}
