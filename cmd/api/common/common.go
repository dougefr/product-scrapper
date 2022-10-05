package common

type (
	// CustomResponseError erro customizado de API
	CustomResponseError interface {
		error
		StatusCode() int
		Message() string
	}

	customResponseError struct {
		statusCode int
		message    string
	}

	// ResponseError resposta de erro
	ResponseError struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	}
)

// NewCustomResponseError cria um novo CustomResponseError
func NewCustomResponseError(statusCode int, message string) CustomResponseError {
	return customResponseError{
		statusCode: statusCode,
		message:    message,
	}
}

func (c customResponseError) Error() string {
	return c.message
}

func (c customResponseError) StatusCode() int {
	return c.statusCode
}

func (c customResponseError) Message() string {
	return c.message
}
