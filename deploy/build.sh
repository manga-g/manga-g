#!/usr/bin/env bash

## Build Context for the Project
echo "Checking project root directory"
cd ../ || return

# Get Dependencies
echo "Checking for code dependencies"
go mod tidy && go mod vendor

# Moving back to the Deploy Directory
cd deploy || return
echo "Starting build from $PWD"

# Build the Project Binary
go build -o ./cmd/core/manga-g
echo "Build Script has completed"