package entity

import "github.com/dougefr/product-scrapper/domain/businesserr"

// Product entidade Product
type Product struct {
	Title       string
	Description string
	Image       string
	Price       float64
	URL         string
}

// NewProduct cria um novo Product
func NewProduct(
	title string,
	description string,
	image string,
	price float64,
	URL string,
) Product {

	return Product{
		Title:       title,
		Image:       image,
		Price:       price,
		Description: description,
		URL:         URL,
	}
}

// IsValid valida o Product
func (p Product) IsValid() businesserr.BusinessError {
	if p.Title == "" {
		return businesserr.InvalidProductTitle
	}

	if p.Image == "" {
		return businesserr.InvalidProductImage
	}

	if p.Price <= 0 {
		return businesserr.InvalidProductPrice
	}

	if p.Description == "" {
		return businesserr.InvalidProductDescription
	}

	if p.URL == "" {
		return businesserr.InvalidProductURL
	}

	return nil
}
