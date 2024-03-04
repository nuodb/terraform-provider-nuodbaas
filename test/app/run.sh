#!/usr/bin/env bash

cd $( dirname $0 )

WS_ROOT=$( realpath ../.. )
export TF_CLI_CONFIG_FILE=./conf.tfrc
TERRAFORM_BIN_PATH=bin/terraform

make -C $WS_ROOT package $TERRAFORM_BIN_PATH
: ${TERRAFORM_BIN:=$WS_ROOT/$TERRAFORM_BIN_PATH}

cat << EOF > $TF_CLI_CONFIG_FILE
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

USER_ORG=$( echo $NUODB_CP_USER | cut -d '/' -f 1 )

$TERRAFORM_BIN init
$TERRAFORM_BIN apply -auto-approve -var "org_name=${USER_ORG}"

exit_code=$( $TERRAFORM_BIN output exit )

$TERRAFORM_BIN destroy -auto-approve -var "org_name=${USER_ORG}"

exit $exit_code
