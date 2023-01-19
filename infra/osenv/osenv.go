package osenv

import (
	"os"

	"github.com/dougefr/product-scrapper/domain/contract/env"
)

// Env implementação de env.Env
type Env struct {
	// RedisAddrAttr conexão com o REDIS
	RedisAddrAttr string

	// ChromiumAddrAttr conexão com o CHROMIUM
	ChromiumAddrAttr string

	// NameAttr nome do ambiente
	NameAttr string
}

// NewEnv cria um novo Env
func NewEnv() env.Env {
	return Env{
		NameAttr:         os.Getenv("ENV_NAME"),
		RedisAddrAttr:    os.Getenv("REDIS_ADDR"),
		ChromiumAddrAttr: os.Getenv("CHROMIUM_ADDR"),
	}
}

// Name retorna o nome do ambiente
func (e Env) Name() string {
	return e.NameAttr
}

// CacheAddr retorna o endereço do cache
func (e Env) CacheAddr() string {
	return e.RedisAddrAttr
}

// BrowserAddr retorna o endereço do browser
func (e Env) BrowserAddr() string {
	return e.ChromiumAddrAttr
}
