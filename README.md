# statuschecker
Status Checker server checks the status of registered services.
The services are checked by testing a HTTP endpoint registered with each service

# Design

The status checker will allow the registration of services and will periodically poll the given health 
endpoint for each service. Registered services will be persisted to a database. Each services current status is 
be cached in memory.

A HTTP api will have the registration and status retrieval functions

# Code

## Main Code
* openapispec.yaml  -  Defines the web API
* db/migrations     -  Defines the database migration scripts for managing the DB schema
* db/queries        -  Defines the SQL used by the app 
* internal/statuschecker/service.go   -  main implementation
* internal/web/controller.go - web API implementation

## Generated code

Much of the boilerplate code is generated.

* Database access code
* Web service stubs and API model
* Mocks

Generation tools

* github.com/golang/mock/mockgen@v1.6.0
* github.com/deepmap/oapi-codegen/cmd/oapi-codegenv1.11.0
* sqlc


### Database

migrate and sqlc are used for database interaction. The database code is all generated from the DDL nd SQL.

To generate the code
```
make generate-db
```
The generated code is written to internal/store


### Web API
The web api code (models and server stub) is generated from the openapispec.yaml file
The server stub generated is for the middleware.Echo server.

To generate
```
make generate-web
```

# Configuration

Configuration is through viper.
The defaults are configured in config/config.yaml. The config file path may be supplied on the command line by setting `--config filepath`.

Environment variables may be used to override the defaults in the file  
E.g. 
```
APP_POLLINTERVAL=5m bin/statuschecker
```

# TODO

1. Unit test web controller
2. echo using zap logging



# Development

## Quick Start
```bash
make image # compile, test and create the docker image
make dbup  # start the database
make dbinstall  # install the schema in the database. NB, wait a few seconds after dbup for the database to start
make run  # Run the application
curl -v localhost:8080/ready # should get a 200 if app is ready to serve
curl -v localhost:8080/metrics # get the metrics for the app (prometheus compatible)
make stop # stops both the app and the database
```
## Tools

* go 1.18
* docker
* make


## Database

The database can be run locally for development.

The database is controlled using docker-compose via the makefile 
```bash
make dbup
make dbdown
```

The schema is created using `migrate`.

```bash
make dbinstall
```


## Testing
To test, start the DB, run the app and curl the API (see quickstart above).  
There are some example curl scripts in test/scripts 


# Operations - running the app

Prometheus metrics are available on `hostname:port/metrics`
```
curl localhost:8080/metrics
```

Alive check (status 200 if alive) on : `hostname:port/alive`
```
curl -v localhost:8080/alive
```

Ready check (status 200 if ready) on : `hostname:port/ready`
```
curl -v localhost:8080/ready
```
