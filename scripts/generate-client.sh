#!/bin/sh

set -e

# Pin version to 7.1.0 to avoid errors due to properties not declared in
# OpenAPI spec (status.conditions for projects and databases):
# https://github.com/OpenAPITools/openapi-generator/issues/17156
OPENAPI_GEN_VERSION="v7.1.0"
OPENAPI_GEN_IMAGE="openapitools/openapi-generator-cli:$OPENAPI_GEN_VERSION"

NUODB_CP_VERSION="v2.3.0"
NUODB_CP_SPEC_URL="https://raw.githubusercontent.com/nuodb/nuodb-cp-releases/$NUODB_CP_VERSION/openapi.yaml"

SCRIPT_DIR="$(dirname "$0")"

cd "$SCRIPT_DIR/.."
CLIENT_DIR=client

# Generate Golang client
[ -e "$CLIENT_DIR" ] && rm -r "$CLIENT_DIR"
docker run --rm -v "$(pwd)/$CLIENT_DIR:/out" "$OPENAPI_GEN_IMAGE" \
    generate -i "$NUODB_CP_SPEC_URL" -g go -o /out -p withGoMod=false --package-name nuodbaas \
    --git-host github.com --git-user-id nuodb --git-repo-id terraform-provider-nuodbaas/client

# Post process generated code as workaround for openapi-generator bug:
# https://github.com/OpenAPITools/openapi-generator/issues/13166
echo "Post processing generated code..."
for file in $(grep -l 'MapmapOfStringinterface{}' -R "$CLIENT_DIR" | grep '[.]go$'); do
    echo "$file"
    sed 's/MapmapOfStringinterface{}/Object/g' "$file" > "${file}.tmp"
    mv "${file}.tmp" "$file"
done

# Delete test files, which generate useless noise in test runs
find "$CLIENT_DIR" -name \*_test.go -exec rm {} \;
