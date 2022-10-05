package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/dougefr/product-scrapper/domain/businesserr"
	mock_cache "github.com/dougefr/product-scrapper/domain/contract/mock/cache"
	mock_logger "github.com/dougefr/product-scrapper/domain/contract/mock/logger"
	mock_scrapper "github.com/dougefr/product-scrapper/domain/contract/mock/scrapper"
	"github.com/dougefr/product-scrapper/domain/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const validURL = "https://www.amazon.com.br/Colch%C3%A3o-Ensacadas-Espuma-Viscoel%C3%A1stica-Ortop%C3%A9dico/dp/B08LH5L7TK/ref=cm_cr_arp_d_product_top?ie=UTF8&th=1"

func TestScrapProduct_ExecuteShouldResultInBusinessErrorWhenScrapperFactoryReturnABusinessError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())

	cache := mock_cache.NewMockCache[entity.Product](ctrl)
	cache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)

	factory := mock_scrapper.NewMockFactory(ctrl)
	factory.EXPECT().GetScrapper(gomock.Any(), gomock.Any()).Return(nil, businesserr.ScrapperNotImplemented)

	uc := NewScrapProduct(logger, cache, factory)
	_, err := uc.Execute(context.Background(), NewScrapProductInput(validURL))

	assert.Equal(t, businesserr.ScrapperNotImplemented, err)
}

func TestScrapProduct_ExecuteShouldResultInAnErrorWhenScrapperFactoryReturnAnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())

	cache := mock_cache.NewMockCache[entity.Product](ctrl)
	cache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)

	factory := mock_scrapper.NewMockFactory(ctrl)
	factory.EXPECT().GetScrapper(gomock.Any(), gomock.Any()).Return(nil, errors.New("fake error"))

	uc := NewScrapProduct(logger, cache, factory)
	_, err := uc.Execute(context.Background(), NewScrapProductInput(validURL))

	assert.EqualError(t, err, "error when getting scrapper from factory: fake error")
}

func TestScrapProduct_ExecuteShouldResultInAnErrorWhenScrapperReturnAnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())

	cache := mock_cache.NewMockCache[entity.Product](ctrl)
	cache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)

	scrapper := mock_scrapper.NewMockScrapper(ctrl)
	scrapper.EXPECT().ScrapProduct(gomock.Any()).Return(entity.Product{}, errors.New("fake error"))

	factory := mock_scrapper.NewMockFactory(ctrl)
	factory.EXPECT().GetScrapper(gomock.Any(), gomock.Any()).Return(scrapper, nil)

	uc := NewScrapProduct(logger, cache, factory)
	_, err := uc.Execute(context.Background(), NewScrapProductInput(validURL))

	assert.EqualError(t, err, "error when getting product from scrapper: fake error")
}

func TestScrapProduct_ExecuteShouldResultInAnErrorWhenScrapperReturnAnInvalidProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
	logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())

	cache := mock_cache.NewMockCache[entity.Product](ctrl)
	cache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)

	scrapper := mock_scrapper.NewMockScrapper(ctrl)
	scrapper.EXPECT().ScrapProduct(gomock.Any()).Return(entity.Product{}, nil)

	factory := mock_scrapper.NewMockFactory(ctrl)
	factory.EXPECT().GetScrapper(gomock.Any(), gomock.Any()).Return(scrapper, nil)

	uc := NewScrapProduct(logger, cache, factory)
	_, err := uc.Execute(context.Background(), NewScrapProductInput(validURL))

	assert.Equal(t, businesserr.InvalidProductTitle, err)
}

func TestScrapProduct_ExecuteShouldResultInAProductInfoWhenScrapperReturnAValidProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).Times(3)

	cache := mock_cache.NewMockCache[entity.Product](ctrl)
	cache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)
	cache.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	scrapper := mock_scrapper.NewMockScrapper(ctrl)
	scrapper.EXPECT().ScrapProduct(gomock.Any()).Return(entity.NewProduct(
		"title",
		"description",
		"image",
		1.1,
		validURL,
	), nil)

	factory := mock_scrapper.NewMockFactory(ctrl)
	factory.EXPECT().GetScrapper(gomock.Any(), gomock.Any()).Return(scrapper, nil)

	uc := NewScrapProduct(logger, cache, factory)
	output, err := uc.Execute(context.Background(), NewScrapProductInput(validURL))

	assert.Nil(t, err)
	assert.Equal(t, NewScrapProductOutput(
		"title",
		"description",
		"image",
		1.1,
		validURL,
	), output)
}

func TestScrapProduct_ExecuteShouldResultInAProductInfoWhenScrapperReturnAValidProductButCacheFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).Times(3)
	logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)

	cache := mock_cache.NewMockCache[entity.Product](ctrl)
	cache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("fake error"))
	cache.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("fake error"))

	scrapper := mock_scrapper.NewMockScrapper(ctrl)
	scrapper.EXPECT().ScrapProduct(gomock.Any()).Return(entity.NewProduct(
		"title",
		"description",
		"image",
		1.1,
		validURL,
	), nil)

	factory := mock_scrapper.NewMockFactory(ctrl)
	factory.EXPECT().GetScrapper(gomock.Any(), gomock.Any()).Return(scrapper, nil)

	uc := NewScrapProduct(logger, cache, factory)
	output, err := uc.Execute(context.Background(), NewScrapProductInput(validURL))

	assert.Nil(t, err)
	assert.Equal(t, NewScrapProductOutput(
		"title",
		"description",
		"image",
		1.1,
		validURL,
	), output)
}

func TestScrapProduct_ExecuteShouldResultInAProductInfoWhenThereIsAValueInCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)

	cache := mock_cache.NewMockCache[entity.Product](ctrl)
	product := entity.NewProduct(
		"title",
		"description",
		"image",
		1.1,
		validURL,
	)
	cache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&product, nil)

	factory := mock_scrapper.NewMockFactory(ctrl)

	uc := NewScrapProduct(logger, cache, factory)
	output, err := uc.Execute(context.Background(), NewScrapProductInput(validURL))

	assert.Nil(t, err)
	assert.Equal(t, NewScrapProductOutput(
		"title",
		"description",
		"image",
		1.1,
		validURL,
	), output)
}

func TestScrapProduct_ExecuteShouldResultInAProductInfoWhenThereIsAnInvalidValueInCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).Times(3)
	logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)

	cache := mock_cache.NewMockCache[entity.Product](ctrl)
	product := entity.NewProduct(
		"",
		"description",
		"image",
		1.1,
		validURL,
	)
	cache.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&product, nil)
	cache.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("fake error"))

	scrapper := mock_scrapper.NewMockScrapper(ctrl)
	scrapper.EXPECT().ScrapProduct(gomock.Any()).Return(entity.NewProduct(
		"title",
		"description",
		"image",
		1.1,
		validURL,
	), nil)

	factory := mock_scrapper.NewMockFactory(ctrl)
	factory.EXPECT().GetScrapper(gomock.Any(), gomock.Any()).Return(scrapper, nil)

	uc := NewScrapProduct(logger, cache, factory)
	output, err := uc.Execute(context.Background(), NewScrapProductInput(validURL))

	assert.Nil(t, err)
	assert.Equal(t, NewScrapProductOutput(
		"title",
		"description",
		"image",
		1.1,
		validURL,
	), output)
}
