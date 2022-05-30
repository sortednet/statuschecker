#!/bin/bash -eux

make image GOOS=linux GOARCH=amd64 IMAGE_NAME=sortednet/statuschecker:1.0.0