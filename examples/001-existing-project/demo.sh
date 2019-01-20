#!/bin/sh

set -e
set -x

go mod init github.com/LawnGnome/what-the-dep-example
go build
