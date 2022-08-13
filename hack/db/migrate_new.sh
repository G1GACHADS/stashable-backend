#!/bin/bash

set -ex

if [ -z "$1" ]; then
    echo "Missing name of migration"
    echo "Usage: $0 NAME"
    exit 1
fi

migrate create -ext sql -dir db -seq $1