#!/usr/bin/env bash

set -eux
echo "Installing lint search engine..."
go get -u -v github.com/alecthomas/gometalinter/
gometalinter --config=linter.json ./... --install

echo "Looking for lint..."
gometalinter ./... --config=linter.json