#!/bin/sh

# cd to "root" directory (one above scripts)
cd "$(dirname "$0")"/.. || exit

go run main.go load.go util.go 1.go
