package scrapper

import (
	"context"

	"github.com/dougefr/product-scrapper/domain/contract/env"
	logger2 "github.com/dougefr/product-scrapper/domain/contract/logger"
)

type amazon struct {
	*base
}

func newAmazon(
	ctx context.Context,
	logger logger2.Logger,
	url string,
	env env.Env,
) (a *amazon, err error) {

	b, err := newBase(
		ctx,
		logger,
		url,
		env,
		"span[id=\"productTitle\"]",
		"div[id=\"feature-bullets\"]",
		"img[id=\"landingImage\"]",
		"span.a-price > span.a-offscreen")

	a = &amazon{
		base: b,
	}

	return
}
