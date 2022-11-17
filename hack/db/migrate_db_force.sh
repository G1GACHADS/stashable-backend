#!/bin/bash

set -ex

if [ -z "$1" ]; then
    echo "Missing version of migration"
    echo "Usage: $0 VERSION"
    exit 1
fi

migrate -path db -database \
    postgres://postgres:postgres@localhost/stashable\?sslmode=disable force $1
