# Deployment scripts for Terraform tests

This directory contains scripts for configuring, creating, and destroying DBaaS Control Plane instances used to run the Terraform tests in this repository.
Each sub-directory represents a deployment driver and must have the following scripts:

- `setup.sh`: Creates an instance of the NuoDB Control Plane to test against.
- `env.sh`: Displays the environment variables that can be set to enable testing against the environment.
- `collect.sh`: Collects diagnostic logging from the environment.
- `teardown.sh`: Tears down the created instance of the NuoDB Control Plane.
