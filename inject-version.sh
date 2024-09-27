#!/bin/sh

# Inject current version into example provider configuration
cd "$(dirname "$0")"
version="$(./get-version.sh)"
file=examples/provider/provider.tf
sed "s/{{ version }}/$version/" "$file.tmpl" > "$file"
