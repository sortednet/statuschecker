#!/bin/bash -eux

curl -X POST -H "Content-Type: application/json" localhost:8080/service -d '{"name":"nourl"}'
