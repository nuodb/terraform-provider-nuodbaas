# Terraform Provider for NuoDB Control Plane

[Documentation](docs/index.md)

## Installing

You can install the provider from the official repo using **TBD**

### Installing development build

You can also build the provider locally and install it into your terraform config.

1. Build the provider:
    ```sh
    make package
    ```

2. Add the provider package as a [`filesystem_mirror`](https://developer.hashicorp.com/terraform/cli/config/config-file#filesystem_mirror) in your terraform configuration. Either edit `~/.terraformrc` or create a new terraform configuration.
    ```hcl
    provider_installation {
        filesystem_mirror {
            path    = "<REPO PATH>/dist/pkg_mirror"
            include = ["registry.terraform.io/nuodb/nuodbaas"]
        }
        direct {
            exclude = ["registry.terraform.io/nuodb/nuodbaas"]
        }
    }
    ```
    Where `<REPO PATH>` is the absolute path of this git repository.

3. In your terraform workspace, delete the `.terraform` dirrectory, which contains cached versions of providers. This needs to be done any time you rebuild the package.
    ```sh
    rm -r .terraform
    ```

4. Use the provider as you normally would.
    ```hcl
    terraform {
        required_providers {
            nuodbaas = {
                source = "registry.terraform.io/nuodb/nuodbaas"
                version = "0.1.0"
            }
        }
    }
    ```

5. If you did not edit your `~/.terraformrc` and instead created a custom configuration, don't forget to set `TF_CLI_CONFIG_FILE` to that file.

    ```sh
    TF_CLI_CONFIG_FILE="config.tfrc" terraform init
    ```
For more details about terraform cli configuration, see the [documentation](https://developer.hashicorp.com/terraform/cli/config/config-file).

## Local testing

1. Setup Kubernetes and have it accessible via kubectl.
Ensure that you have a CSI driver configured.

2. Clone this repository.

    ```sh
    git clone git@github.com:nuodb/terraform-provider-nuodbaas.git
    cd terraform-provider-nuodbaas
    ```

3. Install an instance of the control plane to test against.
   
    ```sh
    make deploy-cp
    ```

4. Run acceptance tests.

    ```sh
    eval "$(make extract-creds)"
    make testacc
    ```

    The `extract-creds` target extracts the credentials for the `system/admin` user from the Kubernetes cluster and prints them to standard output. If you want to manually configure the Control Plane instance, set  `NUODB_CP_URL_BASE`, `NUODB_CP_ORGANIZATION`, `NUODB_CP_USER`, and `NUODB_CP_PASSWORD` environment variables and run:

    ```sh
    make testacc
    ```

To run a single test:

```sh
TESTARGS="-run='TestAccProjectResource'" make testacc
```
