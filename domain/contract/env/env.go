package env

type (
	// Env retorna informações sobre o ambiente que a aplicação está sendo executada
	Env interface {
		// Name nome do ambiente
		Name() string

		// CacheAddr endereço do cache
		CacheAddr() string

		// BrowserAddr endereço do cache
		BrowserAddr() string
	}
)
