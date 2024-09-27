#!/bin/sh

# Scrape the version from the main.go file so that we do not specify the same
# value in multiple places.
cd "$(dirname "$0")"
sed -n 's|^\t*version string *= "\([^"]*\)" // {{version}}|\1|p' main.go
