#!/bin/sh

# cd to "root" directory (one above scripts)
cd "$(dirname "$0")"/.. || exit

mkdir -p build

go build -o build/greek cmd/greek/greek.go|| exit

./build/greek "$@"
