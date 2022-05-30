
GOARCH ?= amd64
GOOS ?= linux
GO111MODULE=on
POSTGRESQL_URL ?= postgres://statuschecker:password@localhost:5432/statuschecker?sslmode=disable
IMAGE_NAME ?= statuschecker


test:
	go test -coverprofile coverage.out ./...

vet:
	go vet ./...

clean:
	rm -f bin/statuschecker

build: vet
	GOARCH=${GOARCH} \
	GOOS=${GOOS} \
	GO111MODULE=on \
	go build -o bin/statuschecker main.go

buildAll: generate build test

image:
	make clean build GOARCH=amd64 GOOS=linux  # because the image is a amd64/linux image
	docker build -t ${IMAGE_NAME} .

run: image
	docker-compose up -d statuschecker

stop:
	docker-compose down

dbup:
	docker-compose up -d db

dbinstall: dbup
	docker run --platform linux/amd64 --rm -v $(PWD):/src --network host migrate/migrate:4 -database ${POSTGRESQL_URL} -path /src/db/migrations up

generate: generate-db generate-web generate-mocks

generate-db:
	docker run --rm -v $(PWD):/src -w /src kjconroy/sqlc:1.13.0 generate

generate-mocks: install-mockgen
	mockgen -destination test/mocks/queries.go -package mocks github.com/sortednet/statuschecker/internal/statuschecker DbQuery
	mockgen -destination test/mocks/httpClient.go -package mocks github.com/sortednet/statuschecker/internal/statuschecker HttpClient


generate-web: install-oapi-codegen
	oapi-codegen -package=web openapispec.yaml  > internal/web/statuschecker.gen.go

install-oapi-codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.11.0

install-mockgen:
	go install github.com/golang/mock/mockgen@v1.6.0
