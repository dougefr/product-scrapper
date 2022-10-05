package businesserr

type (
	// BusinessError ...
	BusinessError interface {
		error
		Code() string
	}

	businessError struct {
		error string
		code  string
	}
)

// AllBusinessErrors mapeia todos os BusinessError criados apartir do NewBusinessError
var AllBusinessErrors []BusinessError

// NewBusinessError cria um novo BusinessError
func NewBusinessError(code, error string) BusinessError {
	b := businessError{
		error: error,
		code:  code,
	}

	AllBusinessErrors = append(AllBusinessErrors, b)

	return b
}

// Error ...
func (b businessError) Error() string {
	return b.error
}

// Code ...
func (b businessError) Code() string {
	return b.code
}

// Business errors que os usecases podem retornar
var (
	// InvalidProductTitle o título do produto é invalido
	InvalidProductTitle = NewBusinessError("ErrInvalidProductTitle", "product title is invalid")

	// InvalidProductImage a imagem do produto é inválida
	InvalidProductImage = NewBusinessError("ErrInvalidProductImage", "product image is invalid")

	// InvalidProductPrice o preço do produto é inválido
	InvalidProductPrice = NewBusinessError("ErrInvalidProductPrice", "product price is invalid")

	// InvalidProductDescription a descrição do produto é inválida
	InvalidProductDescription = NewBusinessError("ErrInvalidProductDescription", "product description is invalid")

	// InvalidProductURL a URL do produto é inválida
	InvalidProductURL = NewBusinessError("ErrInvalidProductURL", "product URL is invalid")

	// ScrapperNotImplemented o scrapper não foi implementado
	ScrapperNotImplemented = NewBusinessError("ErrScrapperNotImplemented", "there is no scrapper implemented for this site")
)
