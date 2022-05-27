
GOARCH ?= amd64
GOOS ?= linux
POSTGRESQL_URL ?= postgres://statuschecker:password@localhost:5432/statuschecker?sslmode=disable

.PHONY: test
test:
	go test -coverprofile coverage.out ./...

.PHONY: vet
vet:
	go vet ./...

build: generate vet
	GOARCH=${GOARCH} \
	GOOS=${GOOS} \
	GO111MODULE=on \
	go build -o bin/statuschecker main.go

image: build
	docker build -t statuschecker .

.PHONY: dbup
dbup:
	docker-compose up -d db

.PHONY: dbdown
dbdown:
	docker-compose down

dbinstall: dbup
	docker run --platform linux/amd64 --rm -v $(PWD):/src --network host migrate/migrate:4 -database ${POSTGRESQL_URL} -path /src/db/migrations up

CODE_GEN=oapi-codegen
.PHONY: generate
generate:
	docker run --rm -v $(PWD):/src -w /src kjconroy/sqlc:1.13.0 generate
	GO111MODULE=on \
	mockgen  -destination test/mocks/queries.go -package mocks github.com/sortednet/statuschecker/internal/statuschecker DbQuery
	mockgen  -destination test/mocks/httpClient.go -package mocks github.com/sortednet/statuschecker/internal/statuschecker HttpClient


webgen:
	${CODE_GEN} -package=web openapispec.yaml  > internal/web/statuschecker.gen.go

