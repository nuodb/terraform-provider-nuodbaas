#!/bin/sh

# Patch types.go to fix names for enum constants. This is a workaround for https://github.com/deepmap/oapi-codegen/issues/1614.
file=openapi/types.go
sed 's/\([A-Za-z0-9_]*\)\(  *\)\(ProjectStatusModelState\) = "\1"/\3\1\2\3 = "\1"/' "$file" > "${file}.tmp"

# Print any changes generated and update file
changes="$(diff -u "$file" "${file}.tmp")"
if [ -n "$changes" ]; then
    echo "Changes generated to $file:"
    echo "$changes"
    mv "${file}.tmp" "$file"
fi
