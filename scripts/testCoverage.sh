#!/bin/bash
# This script is used to run the test coverage report for the project

cd "$(dirname "$0")"/.. || exit

go test -coverprofile=coverage.out ./... || exit

go tool cover -html coverage.out -o coverage.html || exit
