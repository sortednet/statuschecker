#!/bin/bash -eux

curl -X POST -H "Content-Type: application/json" localhost:8080/service -d '{"name":"checker", "url":"http://localhost:8080/health"}'
curl -X POST -H "Content-Type: application/json" localhost:8080/service -d '{"name":"google", "url":"http://google.com"}'