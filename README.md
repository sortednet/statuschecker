# statuschecker
Example Status Checker server written in Golang

# Design

The status checker will allow the registration of services and will periodically poll the given health 
endpoint for each service. Registered services will be persisted to a database. The status will be cached in memory.

A HTTP api will have the registration and status retrieval functions


## Code

### Main Code
openapispec.yaml  -  Defines the web API
db/migrations     -  Defines the database migration scripts for managing the DB schema
db/queries        -  Defines the SQL used by the app 
internal/statuschecker/service.go   -  main implementation
internal/web/controller.go - web API implementation

### Generated code

Code is generated by 
```
make generate
```

NB:
local tools are currently relied on 
```bash
go install github.com/golang/mock/mockgen@v1.6.0
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen
```

#### Database

migrate and sqlc are used for database interaction. The database code is all generated from the DDL nd SQL.

To generate the code
```
make generate
```
The generated code is written to internal/store


#### Web API
The web api code (models and server stub) is generated from the openapispec.yaml file
The server stub generated is for the middleware.Echo server.

To generate
```
make oapicodegen
```

??? Mock the Queries type for testing or create a repository interface ???

API written in openapispec with the server stub generated (Echo)

Very small service layer that 
* manages the cache 
* sets promauto values
* get service status

Poller that calls out to each service and updates the cache




# TODO

1. Unit test controller
2. Add a prometheus scrape point for testing
3. Zap logging time format
4. echo using zap logging
5. CI
   1. Dockerfile for app
   2. DockerCompose to also start app 
      1. Dockerfile for build
   3. docker-compose to build image as well as run it
   4. Build script ensuring GOARCH, GOOS etc are all setup for linux (as the image is linux/amd64)




# Development

## Tools

Need an image with
sqlc
migrate
mockgen
openapi gen

## Local Environment

go install github.com/golang/mock/mockgen@v1.6.0
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen

### Database

The database can be run locally for development

```bash
make dbup
make dbdown
```

### Testing

There are some example curl scripts in test/scripts 