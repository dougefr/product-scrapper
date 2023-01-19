#!make
-include .env
export

path=./...
GOPATH=$(shell go env GOPATH)

setup-test: .make.setup-test
.make.setup-test:
	go install github.com/kyoh86/richgo@latest
	go install github.com/golang/mock/mockgen@v1.7.0-rc.1
	touch .make.setup-test

setup-lint: .make.setup-lint
.make.setup-lint:
	go install golang.org/x/lint/golint@latest
	touch .make.setup-lint

setup-sec: .make.setup-sec
.make.setup-sec:
	wget -O - -q https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s latest
	touch .make.setup-sec

setup-swag: .make.setup-swag
.make.setup-swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	touch .make.setup-swag

setup: .make.setup .make.setup-test .make.setup-lint .make.setup-sec .make.setup-swag
.make.setup:
	go install golang.org/x/tools/cmd/goimports@latest
	touch .make.setup

download: setup
	go mod download

run: setup swag
	go run cmd/api/main.go

sec: setup-sec
	./bin/gosec $(path)

fmt: setup
	go fmt $(path)
	find . -name \*.go -exec $(GOPATH)/bin/goimports -w {} \;

lint: setup-lint mock
	golint -set_exit_status -min_confidence 0.9 $(path)
	@echo "Golint found no problems on your code!"
	go vet $(path)

test: mock
	richgo test $(path) $(args)

fullcover: mock
	go test -coverprofile=coverage.out $(path)
	go tool cover -func=coverage.out

swag: setup-swag
	swag init -g cmd/api/main.go

mock: setup-test
	mockgen -source=domain/usecase/scrap_product.go -destination=domain/usecase/mock/scrap_product.go
	mockgen -source=domain/contract/logger/logger.go -destination=domain/contract/mock/logger/logger.go
	mockgen -source=domain/contract/env/env.go -destination=domain/contract/mock/env/env.go
	mockgen -source=domain/contract/cache/cache.go -destination=domain/contract/mock/cache/cache.go
	mockgen -source=domain/contract/scrapper/scrapper.go -destination=domain/contract/mock/scrapper/scrapper.go

build: swag
	GOOS=linux GOARCH=amd64 go build -o dist/api ./cmd/api/main.go