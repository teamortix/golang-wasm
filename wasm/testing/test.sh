#!/bin/sh
# Runs Go tests by adding the correct executable to your $PATH and setting GOOS and GOARCH. 
# Providing which file to test is optional.

PATH="${PATH}:$(realpath $(dirname $BASH_SOURCE))"

GOOS=js GOARCH=wasm go test $1

