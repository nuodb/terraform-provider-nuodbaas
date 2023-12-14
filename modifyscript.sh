#!/bin/sh

# Issue -> https://github.com/OpenAPITools/openapi-generator/issues/13166

for file in $(grep -l 'MapmapOfStringinterface{}' -R . | grep '[.]go$'); do
    echo $file
    sed 's/MapmapOfStringinterface{}/Object/g' "$file" > "${file}.tmp"
    mv "${file}.tmp" "$file"
done

echo "Conversion complete."