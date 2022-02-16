#!/bin/sh

# cd to "root" directory (one above scripts)
cd "$(dirname "$0")"/.. || exit

mkdir -p build

go build -o build/BallotCleaner main.go load.go util.go 1.go 2.go 3.go 4.go || exit

# if 1 or more args then shift
if [ $# -gt 1 ]; then
    shift 1
fi

./build/BallotCleaner "$@"
