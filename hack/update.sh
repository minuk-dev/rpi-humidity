#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail
go get -u ./... || true
go mod tidy || true
