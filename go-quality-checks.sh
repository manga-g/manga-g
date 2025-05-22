#!/bin/bash

echo "--- Running gofmt ---"
# gofmt checks if the code is formatted according to Go's standards.
# -l lists files whose formatting differs.
# -s simplifies code (e.g., x[i:len(x)] to x[i:]).
# We'll check for unformatted files and report.
FORMATTED_FILES=$(gofmt -l -s .)
if [ -n "$FORMATTED_FILES" ]; then
    echo "The following files are not formatted correctly (run 'gofmt -w .'):"
    echo "$FORMATTED_FILES"
    # exit 1 # Optionally exit if not formatted
else
    echo "gofmt: All files are correctly formatted."
fi
echo ""

echo "--- Running go vet ---"
# go vet examines Go source code and reports suspicious constructs.
# Ensure go vet is part of your Go installation.
go vet ./...
echo ""

echo "--- Running gocyclo (top 20 most complex functions, sorted) ---"
if ! command -v gocyclo &> /dev/null
then
    echo "gocyclo not found, attempting to install..."
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
    if ! command -v gocyclo &> /dev/null
    then
        echo "Failed to install gocyclo or it's not in PATH. Please install it manually."
    else
        echo "gocyclo installed successfully."
    fi
fi
gocyclo . | sort -nr | head -20
echo ""

echo "--- Running ineffassign ---"
# ineffassign detects when assignments to existing variables are not used.
# Ensure ineffassign is installed: go install github.com/gordonklaus/ineffassign@latest
if ! command -v ineffassign &> /dev/null
then
    echo "ineffassign not found, attempting to install..."
    go install github.com/gordonklaus/ineffassign@latest
    if ! command -v ineffassign &> /dev/null
    then
        echo "Failed to install ineffassign or it's not in PATH. Please install it manually."
        # You might want to exit here or skip this check
    else
        echo "ineffassign installed successfully."
    fi
fi
ineffassign ./...
echo ""

echo "--- Running misspell ---"
# misspell finds commonly misspelled English words in comments and strings.
# Ensure misspell is installed: go install github.com/client9/misspell/cmd/misspell@latest
if ! command -v misspell &> /dev/null
then
    echo "misspell not found, attempting to install..."
    go install github.com/client9/misspell/cmd/misspell@latest
    if ! command -v misspell &> /dev/null
    then
        echo "Failed to install misspell or it's not in PATH. Please install it manually."
        # You might want to exit here or skip this check
    else
        echo "misspell installed successfully."
    fi
fi
misspell -error . # -error makes it exit with non-zero status on misspellings
echo ""

echo "--- Checking for LICENSE file ---"
if [ -f "LICENSE" ] || [ -f "LICENSE.md" ] || [ -f "LICENSE.txt" ] || [ -f "COPYING" ] || [ -f "COPYING.md" ] || [ -f "COPYING.txt" ]; then
    echo "LICENSE file found."
else
    echo "Warning: No LICENSE file found in the root directory. Go Report Card checks for files like LICENSE, LICENSE.md, COPYING, etc."
    # exit 1 # Optionally exit if no license
fi
echo ""

echo "--- All checks based on your Go Report Card complete ---" 