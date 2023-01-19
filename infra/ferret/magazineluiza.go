package ferret

import (
	"context"

	"github.com/dougefr/product-scrapper/domain/contract/env"
	logger2 "github.com/dougefr/product-scrapper/domain/contract/logger"
)

type magazineLuiza struct {
	*base
}

func newMagazineLuiza(
	ctx context.Context,
	logger logger2.Logger,
	url string,
	env env.Env,
) (m *magazineLuiza, err error) {

	b, err := newBase(
		ctx,
		logger,
		url,
		env,
		"h1.title",
		"div[id=\"descricao\"]",
		"#product-1-item > img",
		"p.price-destaque")

	m = &magazineLuiza{
		base: b,
	}

	return
}
