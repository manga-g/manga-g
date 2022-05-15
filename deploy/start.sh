#!/usr/bin/env bash

cd ../
# Check if go is good to go
go version

# Try to link the go mod
echo "Checking if code dependencies are up to date"
go mod tidy && go mod vendor

# Just tries to run the app
echo "Running the app"
go run ../cmd/core/main.go