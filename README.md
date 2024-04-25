# NuoDB DBaaS Provider for Terraform

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/nuodb/terraform-provider-nuodbaas/tree/main.svg?style=shield&circle-token=64fc942b186124a09d998787035b14627614064f)](https://dl.circleci.com/status-badge/redirect/gh/nuodb/terraform-provider-nuodbaas/tree/main)

The NuoDB DBaaS Provider for Terraform allows NuoDB databases to be managed using Terraform within the [NuoDB Control Plane](https://github.com/nuodb/nuodb-cp-releases).

For more information, see the [NuoDB DBaaS Provider](https://registry.terraform.io/providers/nuodb/nuodbaas) page in the Terraform Registry.

## Usage requirements

* Terraform v1.5.x or greater
* Access to NuoDB Control Plane v2.3.x or greater

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

Once DBaaS access has been configured, you can use Terraform to manage NuoDB projects and databases.

The following Terraform configuration defines a NuoDB project (`nuodbaas_project`) and database (`nuodbaas_database`).

```hcl
terraform {
  required_providers {
    nuodbaas = {
      source  = "registry.terraform.io/nuodb/nuodbaas"
      version = "1.1.0"
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

# Expose nuosql arguments to connect to database
output "nuosql_args" {
  value = <<-EOT
  ${nuodbaas_database.db.name}@${nuodbaas_database.db.status.sql_endpoint}:443 \
      --user dba --password ${nuodbaas_database.db.dba_password} \
      --connection-property trustedCertificates='${nuodbaas_database.db.status.ca_pem}'
  EOT
  sensitive = true
}
```

The example above provides the minimum configuration attributes needed in order to create a NuoDB project and database.
A project belongs to an `organization` and has a `name`, `sla`, and `tier`, which define management policies and configuration attributes that are inherited by databases.
A database belongs to `project` and has a `name` and `dba_password`, which is the password of the DBA user that can be used to create less-privileged database users.

> [!IMPORTANT]
> The `organization` and `project` attributes of the database are resolved from the project resource.
This is important because it defines an [implicit dependency](https://developer.hashicorp.com/terraform/tutorials/configuration-language/dependencies#manage-implicit-dependencies) on the project, which forces the project to be created before the database.

For a complete list of project and database attributes, see the [Project](docs/resources/project.md) and [Database](docs/resources/database.md) resource documentation.

### Creating resources

Now that we understand what the configuration defines, we can use Terraform to create the actual resources.

1. Create a file named `example.tf` from the content above.
2. Run `terraform init` to initialize your Terraform workspace.
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

### Connecting to the database

The example configuration above also exposes an [`output`](https://developer.hashicorp.com/terraform/language/values/outputs) named `nuosql_args`, which makes available the arguments to supply to the `nuosql` tool in order to establish a client connection to the database.
The `nuosql` tool can be found in the [NuoDB Client Package](https://github.com/nuodb/nuodb-client) ([v20230228](https://github.com/nuodb/nuodb-client/releases/tag/v20230228) or greater is required in order to connect to DBaaS databases).

Once `nuosql` is installed and available on the system path, you can connect to the database by running the command:

```bash
eval "nuosql $(terraform output -raw nuosql_args)"
```

### Destroying resources

Once you are done using your database and project, you can delete them by running `terraform destroy`.

## Build requirements

* GNU Make 4.4
* Golang 1.19
* Java 11 (for integration testing)
* Minikube / Docker Desktop and Kubernetes 1.28.x (for end-to-end testing)

### Installing development builds

You can also build the provider locally and configure Terraform to use it.
The `make package` command builds and stages the provider in the `dist/pkg_mirror` directory so that it can be used as a local [`filesystem_mirror`](https://developer.hashicorp.com/terraform/cli/config/config-file#filesystem_mirror).

A locally built provider can be used as follows:

```bash
# Build provider, which is staged in dist/pkg_mirror
make package

# Create Terraform configuration
cat <<EOF > terraform.rc
provider_installation {
    filesystem_mirror {
        path    = "$(pwd)/dist/pkg_mirror"
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

## Local testing

There are two configurations that are used by the tests that are run as part of continuous integration.
The integration tests use a stripped-down CRUD-only Kubernetes environment consisting of a Kubernetes API server backed by `etcd`, along with the NuoDB Control Plane REST service.
The end-to-end tests use a real Kubernetes environment to run a full instance of the NuoDB Control Plane as described in [DBaaS Quick Start Guide](https://github.com/nuodb/nuodb-cp-releases/blob/main/docs/QuickStart.md).

### Integration testing

To run integration tests:

```bash
make integration-tests
```

This downloads and deploys the CRUD-only Kubernetes environment and NuoDB Control Plane REST service, runs the tests, and shuts down the test environment when finished.

### End-to-end testing

To run end-to-end tests:

1. Deploy the NuoDB Control Plane into the Kubernetes cluster that `kubectl` is configured to use.
    ```bash
    make deploy-cp
    ```
2. Configure the credentials needed to access the NuoDB Control Plane.
    ```bash
    eval "$(make extract-creds)"
    ```
3. Run the tests.
    ```bash
    make testacc
    ```
4. Clean up all resources created for the NuoDB Control Plane when finished.
    ```bash
    make undeploy-cp
    ```

### Running a single test

To run a single test:

```sh
TESTARGS="-run=TestFullLifecycle" make testacc
```
