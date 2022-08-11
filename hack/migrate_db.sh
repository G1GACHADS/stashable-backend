#!/bin/bash

migrate -path db -database \
    postgres://postgres:postgres@localhost/storage_system\?sslmode=disable up
