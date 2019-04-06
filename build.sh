#!/usr/bin/env bash

go get ./...

go build -o ./bin/bimock ./cmd/bimock
go build -o ./bin/simulator ./cmd/simulator