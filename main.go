package main

import (
    "context"
    "log"

    "github.com/nuodb/terraform-provider-nuodbaas/internal/provider"

    "github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
    opts := providerserver.ServeOpts{
        Address: "registry.terraform.io/nuodb/nuodbaas",
    }

    err := providerserver.Serve(context.Background(), provider.New(), opts)
    if err != nil {
        log.Fatal(err.Error())
    }
}
