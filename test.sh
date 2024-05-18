#!/bin/bash

# 1. Just tests.
# go test -v -vet=all

# 2. Tests and code coverage report.
# go test -v -vet=all -coverprofile=coverage.out
# go tool cover -html=coverage.out

# 3. Tests, benchmarking and code coverage report.
go test -v ./... -vet=all -coverprofile=coverage.out -bench=.
go tool cover -html=coverage.out
