#!/bin/sh

# cd to "root" directory (one above scripts)
cd "$(dirname "$0")"/.. || exit

mkdir -p build

go build -o build/lastYear cmd/lastYear/lastYear.go|| exit

./build/lastYear "$@"
