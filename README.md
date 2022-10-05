# product-scrapper-api

## Requisitos para rodar

- Compilador 1.19 do Golang: https://go.dev/dl/
- Redis
- Chromium

## Como executar

Criar um arquivo `.env`, conforme o exemplo `.env.example`, alterando corretamente as variáveis de ambiente, e executar o comando:

```shell
make run
```
A aplicação ficará disponibilizada na porta 8888, e o swagger pode ser acessado em http://localhost:8888/swagger

### Aplicação rodando com docker compose

Ao invés de executar a aplicação separadamente, é possível executá-la com o docker compose: `docker compose up --build`

### Infraestrutura rodando com docker

A aplicação precisa se conectar a um servidor Redis e Chromium para funcionar corretamente. Para inicializar essa infraestrutura com docker, deve-se executar os seguintes comandos:

- Redis: `docker run -d -p 6379:6379 redis`
- Chromium: `docker run -d -p 9222:9222 montferret/chromium`

## Makefile

No Makefile se encontram comandos úteis para lidar com o dia a dia do projeto:
- `make run`: executa o projeto
- `make test`: roda os testes unitários do projeto
- `make fmt`: forma o código-fonte do projeto
- `make lint`: inspeciona o código-fonte do projeto
- `make sec`: inspeciona o segurança do código-fonte do projeto

## Bibliotecas relevantes utilizadas:

- https://github.com/MontFerret/ferret: biblioteca utilizada para fazer o web scraping
- https://github.com/go-redis/redis: cliente do Redis. Foi utilizado como solução de armazenamento pois é escalável, lida bem com tempo de expiração de cache e armazenamento não estruturado (NoSQL)
- https://github.com/gofiber/fiber: biblioteca amplamente utilizada para o desenvolvimento de APIs REST em golang devido boa performance
- https://github.com/golang/mock: biblioteca para a geração de mock para o desenvolvimento de testes unitários
- https://github.com/stretchr/testify: biblioteca com métodos de asserts para serem utilizados em testes unitários
- https://github.com/swaggo/swag: biblioteca para geração da documentão swagger do projeto
- https://go.uber.org/zap: biblioteca de log estruturados