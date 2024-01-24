# Terraform Provider for NuoDB Control Plane

[Documentation](plugin/docs/index.md)

## Installing

You can install the provider from the official repo using **TBD**

### Installing development build

You can also build the provider locally and install it into your terraform config.

1. Build the provider and install it into `$GOBIN`:
    ```sh
    cd plugin && go install .
    ```

2. Add the new executable into the `dev_overrides` section of your terraform configuration. Either edit `~/.terraformrc` or create a new terraform configuration.
    ```hcl
    provider_installation {

        dev_overrides {
            "hashicorp.com/edu/nuodbaas" = "<GOBIN_PATH>"
        }

        # For all other providers, install them directly from their origin provider
        # registries as normal.
        direct {}
    }
    ```
    Where `<GOBIN_PATH>` is the absolute path of `$GOBIN`.

3. Use the provider as you normally would. You can omit the `version` since that value will be ignored.
    ```hcl
    terraform {
        required_providers {
            nuodbaas = {
                source = "hashicorp.com/edu/nuodbaas"
            }
        }
    }
    ```

4. If you did not edit your `~/.terraformrc` and instead created a custom configuration, don't forget to set `TF_CLI_CONFIG_FILE` to that file.

    ```sh
    TF_CLI_CONFIG_FILE="config.tfrc" terraform plan
    ```
For more details about overriding providers, see the [Terraform documentation](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers).

## Local testing

1. Setup Kubernetes and have it accessible via kubectl.
Ensure that you have a CSI driver configured.

2. Clone this repository.

    ```sh
    git clone https://github.com/nuodb/nuodbaas-tf-plugin
    cd nuodbaas-tf-plugin
    ```

3. Install an instance of the control plane to test against.
   
    ```sh
    make deploy-cp
    ```

4. Run acceptance tests.

    ```sh
    make discover-test
    ```

    If you want to manually configure the Control Plane instance, set  `NUODB_CP_URL_BASE`, `NUODB_CP_ORGANIZATION`, `NUODB_CP_USER`, and `NUODB_CP_PASSWORD` environment variables and run:

    ```sh
    make testacc
    ```

To run a single test:

```sh
TESTARGS="-run='TestAccProjectResource'" make discover-test
```
