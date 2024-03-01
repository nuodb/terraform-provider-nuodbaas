# NuoDB DBaaS Provider for Terraform

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/nuodb/terraform-provider-nuodbaas/tree/main.svg?style=shield&circle-token=64fc942b186124a09d998787035b14627614064f)](https://dl.circleci.com/status-badge/redirect/gh/nuodb/terraform-provider-nuodbaas/tree/main)

The NuoDB DBaaS Provider for Terraform allows NuoDB databases to be managed using Terraform within the [NuoDB Control Plane](https://github.com/nuodb/nuodb-cp-releases).

## Usage requirements

* Terraform v1.5.x
* Access to NuoDB Control Plane v2.3.x

## Installation

The initial release of NuoDB DBaaS provider is in development and has not yet been published to the [Terraform Registry](https://registry.terraform.io/).

In order to use the NuoDB DBaaS provider, download the version of the provider for your platform from the [Releases](https://github.com/nuodb/terraform-provider-nuodbaas/releases) page and set up a local [`filesystem_mirror`](https://developer.hashicorp.com/terraform/cli/config/config-file#filesystem_mirror).
For example, on Linux x86:

```bash
# Create directory for filesystem mirror
mkdir -p pkg_mirror/registry.terraform.io/nuodb/nuodbaas

# Download NuoDB DBaaS provider
curl -L https://github.com/nuodb/terraform-provider-nuodbaas/releases/download/v0.2.0/terraform-provider-nuodbaas_0.2.0_linux_amd64.zip \
    --output-dir pkg_mirror/registry.terraform.io/nuodb/nuodbaas

# Create Terraform configuration
cat <<EOF > terraform.rc
provider_installation {
    filesystem_mirror {
        path    = "$(pwd)/pkg_mirror"
        include = ["registry.terraform.io/nuodb/nuodbaas"]
    }
    direct {
        exclude = ["registry.terraform.io/nuodb/nuodbaas"]
    }
}
EOF

# Use Terraform configuration
export TF_CLI_CONFIG_FILE="$(pwd)/terraform.rc"
```

Alternatively, the Terraform configuration file can be placed in `~/.terraformrc` so that it gets picked up automatically without `TF_CLI_CONFIG_FILE` being set.

## Configuring DBaaS access

Access to an instance of the NuoDB Control Plane is required in order to use the NuoDB DBaaS provider.
If you do not have access to a running instance of the NuoDB Control Plane, you can set one up yourself by following the instructions in [DBaaS Quick Start Guide](https://github.com/nuodb/nuodb-cp-releases/blob/main/docs/QuickStart.md).

The access credentials and URL for the NuoDB Control Plane can be supplied to the NuoDB DBaaS provider as configuration attributes or as environment variables.
In the example below, we use environment variables:

```bash
export NUODB_CP_USER=org/user
export NUODB_CP_PASSWORD=secret
export NUODB_CP_URL_BASE=https://example.dbaas.nuodb.com/api
```

See the [Documentation](docs/index.md) for information on provider configuration.

## Getting started

Once the Terraform provider and DBaaS access have been configured, you can use Terraform to manage NuoDB projects and databases.

The following Terraform configuration defines a NuoDB project (`nuodbaas_project`) and database (`nuodbaas_database`).

```hcl
terraform {
  required_providers {
    nuodbaas = {
      source  = "registry.terraform.io/nuodb/nuodbaas"
    }
  }
}

provider "nuodbaas" {
  # Credentials and URL supplied by environment variables
}

# Create a project
resource "nuodbaas_project" "proj" {
  organization = "org"
  name         = "proj"
  sla          = "dev"
  tier         = "n0.nano"
}

# Create a database within the project
resource "nuodbaas_database" "db" {
  organization = nuodbaas_project.proj.organization
  project      = nuodbaas_project.proj.name
  name         = "db"
  dba_password = "secret"
}
```

The example above provides the minimum configuration attributes needed in order to create a NuoDB project and database.
A project belongs to an `organization` and has a `name`, `sla`, and `tier`, which define management policies and configuration attributes that are inherited by databases.
A database belongs to `project` and has a `name` and `dba_password`, which is the password of the DBA user that can be used to create less-privileged database users.

> **NOTE**: The `organization` and `project` attributes of the database are resolved from the project resource.
This is important because it defines an [implicit dependency](https://developer.hashicorp.com/terraform/tutorials/configuration-language/dependencies#manage-implicit-dependencies) on the project, which forces the project to be created before the database.

For a complete list of project and database attributes, see the [Project](docs/resources/project.md) and [Database](docs/resources/database.md) resource documentation.

### Creating resources

Now that we understand what the configuration defines, we can use Terraform to create the actual resources.

1. Create a file named `example.tf` from the content above.
2. Run `terraform init` to initialize your Terraform workspace.
Since the local filesystem mirror is being used, this command will display a warning about missing checksums that can be ignored.
3. Run `terraform apply` to create the project and database.
In this step, you will be prompted to confirm that you want to create the resources.
Resource creation will be aborted unless you enter `yes`.

By default, the `terraform apply` invocation will block until the project and database both become available, which may take a few minutes.

### Inspecting resources

Once the database is running, you can connect to it by using the `sql_endpoint` and `ca_pem` read-only attributes, which are available after the database resource has been created.
To display all resource attributes, you can run `terraform show` within your workspace.

```console
$ terraform show
...
resource "nuodbaas_database" "db" {
    dba_password = (sensitive value)
    labels       = {}
    name         = "db"
    organization = "org"
    project      = "proj"
    properties   = {
        product_version = "5.0"
        tier_parameters = {}
    }
    status       = {
        ca_pem       = <<-EOT
            -----BEGIN CERTIFICATE-----
            MIICxDCCAaygAwIBAgIJAJdWQkAR7tVqMA0GCSqGSIb3DQEBCwUAMBcxFTATBgNV
            BAMMDGNhLm51b2RiLmNvbTAeFw0yNDAzMDEyMDIyMDlaFw0yNTAzMDEyMDIyMDla
            MBcxFTATBgNVBAMMDGNhLm51b2RiLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEP
            ADCCAQoCggEBAIvrbA9ZNSE4NvWblJJ1dlcM3mBNxnSti9FaFsf20uxydOKn1gpC
            Qwo6hQvXOD+nI5N1a/YSKbp+HyDgU/4XDCVetjpxLzem049v+1Fa00EkM2aMoQQF
            ke4SFkgDUoxqTHZ2lqHEDJwUu9zhwoxbsoxKpGNvmsP3kJvevoFiTKlIfBW8kw/H
            2PXXOjkopvIM7gOr9vIIZV40zJZpp7e0i9MWNWNsL8XJqIYUFCTit0W8CbqTHoO0
            WiYxSeiuV0rtFA1wMpXz5vSp3l9SskbKbQQmJFsst50s80KqZUnIDsiKXBOfZF6l
            pz8K2NOs9XP35CzpVXHF3pnLGVsaY+ZFMtkCAwEAAaMTMBEwDwYDVR0TAQH/BAUw
            AwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAEifAXDhWDIUwF8/GxlXbO5AY60ZS6S9n
            hcaLmVDAYK+JmBEcpmdDQUegOaajsO0zN25DR7Obkfw0Nq5imVYCaVEEmY4j6tL/
            KAOBw6NxT1sy32IHKcREFd4qlgE6pgaET87S4uL6Ej4QlKR0ksB+bBTA6NKgmX+W
            qW8fjXs/q2ey45hUigSCfXspO5bS8OX0iAxQA+5Sw1Ys4GpaM4Y/1FscPCh99hzv
            rpnZ3nuGnXnKy7oceAaG/bBxnMK4BEbm4TIO1MHHdTZL5+gcnuDMtKTxjdpgziud
            ECB524GGqjDOBQ/Ixa6bRT2hM2LaqTcPzlPICl5NNTyhoQlSDIGdYQ==
            -----END CERTIFICATE-----
        EOT
        message      = "database is available"
        ready        = true
        shutdown     = false
        sql_endpoint = "proj-f4ccefc29c8d.dbaas.nuodb.com"
        state        = "Available"
    }
    tier         = "n0.nano"
}
```

### Destroying resources

Once you are done using your database and project, you can delete by running `terraform destroy`.

## Build requirements

* GNU Make 4.4
* Golang 1.21
* Java 11 (for integration testing)
* Minikube / Docker Desktop and Kubernetes 1.28.x (for end-to-end testing)

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
                version = "0.2.0"
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

    The `extract-creds` target extracts the credentials for the `system/admin` user from the Kubernetes cluster and prints them to standard output. If you want to manually configure the Control Plane instance, set  `NUODB_CP_URL_BASE`, `NUODB_CP_USER`, and `NUODB_CP_PASSWORD` environment variables and run:

    ```sh
    make testacc
    ```

To run a single test:

```sh
TESTARGS="-run='TestFullLifecycle'" make testacc
```

Certain tests respect the `-test.short` flag to skip non-critical scenarios. When testing against a real control plane, this option can significantly speed up test execution.
```sh
TESTARGS="-test.short" make testacc
```
