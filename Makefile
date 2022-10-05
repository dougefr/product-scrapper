#!make
-include .env
export

path=./...
GOPATH=$(shell go env GOPATH)

setup-test: .make.setup-test
.make.setup-test:
	GO111MODULE=off go get -u github.com/kyoh86/richgo
	GO111MODULE=off go get -u github.com/golang/mock/mockgen
	GO111MODULE=off go get -u golang.org/x/tools/cmd/cover
	touch .make.setup-test

setup-lint: .make.setup-lint
.make.setup-lint:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	touch .make.setup-lint

setup-sec: .make.setup-sec
.make.setup-sec:
	wget -O - -q https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s latest
	touch .make.setup-sec

setup-swag: .make.setup-swag
.make.setup-swag:
	GO111MODULE=off go get -u github.com/swaggo/swag/cmd/swag
	touch .make.setup-swag

setup: .make.setup .make.setup-test .make.setup-lint .make.setup-sec .make.setup-swag
.make.setup:
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
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
	$(GOPATH)/bin/golint -set_exit_status -min_confidence 0.9 $(path)
	@echo "Golint found no problems on your code!"
	go vet $(path)

test: mock
	$(GOPATH)/bin/richgo test $(path) $(args)

fullcover: mock
	go test -coverprofile=coverage.out $(path)
	go tool cover -func=coverage.out

swag: setup-swag
	$(GOPATH)/bin/swag init -g cmd/api/main.go

mock: setup-test
	$(GOPATH)/bin/mockgen -source=domain/usecase/scrap_product.go -destination=domain/usecase/mock/scrap_product.go
	$(GOPATH)/bin/mockgen -source=domain/contract/logger/logger.go -destination=domain/contract/mock/logger/logger.go
	$(GOPATH)/bin/mockgen -source=domain/contract/env/env.go -destination=domain/contract/mock/env/env.go
	$(GOPATH)/bin/mockgen -source=domain/contract/cache/cache.go -destination=domain/contract/mock/cache/cache.go
	$(GOPATH)/bin/mockgen -source=domain/contract/scrapper/scrapper.go -destination=domain/contract/mock/scrapper/scrapper.go

build: swag
	GOOS=linux GOARCH=amd64 go build -o dist/api ./cmd/api/main.go