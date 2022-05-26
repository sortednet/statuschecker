
GOARCH ?= amd64
GOOS ?= linux
POSTGRESQL_URL ?= postgres://statuschecker:password@localhost:5432/statuschecker?sslmode=disable

test:
	go test ./...

build: generate vet
	GOARCH=${GOARCH} \
	GOOS=${GOOS} \
	GO111MODULE=on \
	go build -o bin/statuschecker main.go

vet:
	go vet ./...

dbup:
	docker-compose up -d db

dbdown:
	docker-compose down

dbinstall: dbup
	docker run --platform linux/amd64 --rm -v $(PWD):/src --network host migrate/migrate:4 -database ${POSTGRESQL_URL} -path /src/db/migrations up

generate:
	docker run --rm -v $(PWD):/src -w /src kjconroy/sqlc:1.13.0 generate