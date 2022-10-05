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

## Bibliotecas relevantes utilizadas

- https://github.com/MontFerret/ferret: biblioteca utilizada para fazer o web scraping
- https://github.com/go-redis/redis: cliente do Redis. Foi utilizado como solução de armazenamento pois é escalável, lida bem com tempo de expiração de cache e armazenamento não estruturado (NoSQL)
- https://github.com/gofiber/fiber: biblioteca amplamente utilizada para o desenvolvimento de APIs REST em golang devido boa performance
- https://github.com/golang/mock: biblioteca para a geração de mock para o desenvolvimento de testes unitários
- https://github.com/stretchr/testify: biblioteca com métodos de asserts para serem utilizados em testes unitários
- https://github.com/swaggo/swag: biblioteca para geração da documentão swagger do projeto
- https://go.uber.org/zap: biblioteca de log estruturados

## Arquitetura do projeto

O projeto está dividido em 3 camadas:

### cmd

O pacote cmd, encontrados em vários projeto Golang, é uma convenção onde são estruturados os entrypoints da aplicação, ou seja, as funções main. Um mesmo projeto pode conter vários entrypoints, com o intuito de compartilhar uma mesma code base. Ou seja, em um mesmo repositório podemos ter uma API Rest, um job e uma API gRPC, conforme a necessidade de negócio.

Neste pacote, é inicializada toda a infraestrutura necessária para a execução dos casos de usos, que existem no pacote domain.

### domain
O pacote domain concentra todo o domínio da aplicação, e existem 3 sub-pacotes com responsabilidades distintas:

#### entity

No pacote entity, as entidades encapsulam as regras de negócios mais centrais da aplicação e de alto nível. Essas regras são menos propensas a mudar quando algo externo muda. Por exemplo, não é esperado que essas entidades fossem afetadas se um mecanismo de segurança e autenticação fosse configurado em uma API.

A modelagem dessas entidades e value objects deve ser feito, de preferencia, de forma não anêmica e sem refletir necessariamente um modelo de armazenamento de dados (ex: tabelas de um banco de dados relacional) ou um contrato com um serviço externo. É encorajado aqui fazer uso ao máximo de conceitos de POO, dentro das limitações da linguagem.

#### usecase

No pacote usecase (caso de uso), concentramos todas as demais regras de negócio da aplicação, orquestrando o fluxo de dados entre um mundo externo abstrato e as entidades. A ideia aqui é que este pacote, juntamente com o entity, seja capaz de transmitir sobre do que se trata a aplicação. Ou seja, uma pessoa que abrir esses pacotes e encontrar entidades como Mercadoria e usecases como Cadastrar Mercadoria e Registrar Saída de Mercadoria, seja possível deduzir que se trata de um sistema de estoque.

Cada usecase é composto de um contrato de entrada, um contrato de saída (ou retorno), e um interactor, que é de fato a implementação do caso de uso. Esses interactors, por meio do princípio de Inversão de Dependências, interagem com o mundo externo por meio de interfaces genéricas, sem saber de fato quem as implementam.

#### contract

Neste pacote se encontram todas as interfaces que serão utilizadas pelo usecases para comunicar com o mundo no exterior.

### infra

O último pacote core da aplicação, o infra, é responsável por armazenar implementações do sub-pacote contract. Neste pacote as classes conhecem detalhes concretos de onde a aplicação está sendo executada e quais são as infraestruturas disponibilizadas para a execução dos usecases. Por exemplo: se um usecase possuir uma dependência com uma interface de Cache, podemos ter uma implementação dessa interface no pacote infra utilizando Redis ou Memcached.

