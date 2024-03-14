# Application connectivity demo

The Terraform configuration files in this directory show how an application can be started along with a NuoDB database.
Terraform manages the dependency of the application on the database based on the presence of values that are interpolated from `nuodbaas_database` resource attributes, like the SQL endpoint and CA PEM.

# Requirements

- Terraform v1.5.x or greater
- Access to NuoDB Control Plane v2.3.x or greater
- Docker runtime environment

# Usage

1. Download the contents of this repository and navigate to this directory:
```
git clone https://github.com/nuodb/terraform-provider-nuodbaas.git
cd terraform-provider-nuodbaas-main/examples/demo
```
Alternatively, you can download and unpack the contents of this directory specifically as follows:
```
curl -s -L https://github.com/nuodb/terraform-provider-nuodbaas/archive/refs/heads/main.tar.gz | \
    tar -xv --strip-components=2 terraform-provider-nuodbaas-main/examples/demo
cd demo
```
2. Configure DBaaS credentials as described in [Configuring DBaaS access](/README.md#configuring-dbaas-access).
3. Run `terraform apply` to create the NuoDB project, database, and the application.
The application runs within a local Docker container.

Once the project, database, and application are created, you can interact with the database and inspect the output of the application container.

## Connecting to database

Connect to the database by using `nuosql`:

    eval "nuosql $(terraform output -raw nuosql_args)"

## Displaying application output

Display the application output by using `docker logs`:

    docker logs -f dbapp

## Checking for expected data

The following command can be used to verify that the `testdata` table has the expected number of rows in it:

    eval "echo 'select count(*) from testdata;' | nuosql $(terraform output -raw nuosql_args)" && ( grep Writing dbapp/out.log | wc -l )

## Cleaning up

When finished, run `terraform destroy` to destroy all resources.
