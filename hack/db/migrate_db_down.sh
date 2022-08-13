#!/bin/bash

set -ex

migrate -path db -database \
    postgres://postgres:postgres@localhost/stashable\?sslmode=disable down
