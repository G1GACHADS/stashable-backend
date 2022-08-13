#!/bin/bash

set -e

migrate -path db -database \
    postgres://postgres:postgres@localhost/stashable\?sslmode=disable down
