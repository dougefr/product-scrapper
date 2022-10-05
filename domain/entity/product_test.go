package entity

import (
	"fmt"
	"testing"

	"github.com/dougefr/product-scrapper/domain/businesserr"
	"github.com/stretchr/testify/assert"
)

func TestProduct_IsValid(t *testing.T) {
	tests := []struct {
		product Product
		wantErr error
	}{
		{
			product: NewProduct(
				"",
				"description",
				"image",
				1.0,
				"url",
			),
			wantErr: businesserr.InvalidProductTitle,
		},
		{
			product: NewProduct(
				"title",
				"",
				"image",
				1.0,
				"url",
			),
			wantErr: businesserr.InvalidProductDescription,
		},
		{
			product: NewProduct(
				"title",
				"description",
				"",
				1.0,
				"url",
			),
			wantErr: businesserr.InvalidProductImage,
		},
		{
			product: NewProduct(
				"title",
				"description",
				"image",
				0,
				"url",
			),
			wantErr: businesserr.InvalidProductPrice,
		},
		{
			product: NewProduct(
				"title",
				"description",
				"image",
				1.0,
				"",
			),
			wantErr: businesserr.InvalidProductURL,
		},
		{
			product: NewProduct(
				"title",
				"description",
				"image",
				1.0,
				"url",
			),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.wantErr), func(t *testing.T) {
			assert.Equal(t, tt.wantErr, tt.product.IsValid())
		})
	}
}
