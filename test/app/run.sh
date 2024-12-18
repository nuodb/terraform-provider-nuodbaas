#!/bin/sh

# Build and stage Terraform provider
cd "$(dirname "$0")"
WS_ROOT="$(cd ../.. && pwd)"
make -C "$WS_ROOT" package

# Download Terraform binary if necessary
if [ -z "$TERRAFORM" ] || [ ! -x "$TERRAFORM" ]; then
    make -C "$WS_ROOT" bin/terraform
    TERRAFORM="$WS_ROOT/bin/terraform"
fi

# Set config to use local provider
export TF_CLI_CONFIG_FILE=./conf.tfrc
cat << EOF > "$TF_CLI_CONFIG_FILE"
provider_installation {
    filesystem_mirror {
        path    = "$WS_ROOT/dist/pkg_mirror"
        include = ["registry.terraform.io/nuodb/nuodbaas"]
    }
    direct {
        exclude = ["registry.terraform.io/nuodb/nuodbaas"]
    }
}
EOF

# Invoke a command and record whether it failed
check_err() {
    if ! $@ ; then
        errors="$errors\n* Command '$*' failed"
    fi
}

errors=""

# Get organization from user name
export TF_VAR_org_name="${NUODB_CP_USER%/*}"

# Initialize Terraform workspace and apply configuration
check_err "$TERRAFORM" init
check_err "$TERRAFORM" apply -auto-approve

# Destroy resources
check_err "$TERRAFORM" destroy -auto-approve

# Check that no errors occurred
if [ -n "$errors" ]; then
    printf "There were errors:$errors"
    exit 1
fi
