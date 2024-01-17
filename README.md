# Terraform Provider for NuoDB Control Plane

[Documentation](plugin/docs/index.md)

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
