#!/bin/bash

set -eo pipefail

go clean
go clean -testcache
go build ./...
go test ./...