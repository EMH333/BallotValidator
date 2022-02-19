#!/bin/sh

# cd to "root" directory (one above scripts)
cd "$(dirname "$0")"/.. || exit

mkdir -p build

go build -o build/BallotCleaner cmd/main/main.go|| exit

./build/BallotCleaner "$@"
