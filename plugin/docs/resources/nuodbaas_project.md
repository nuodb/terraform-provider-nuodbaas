---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nuodbaas_project Resource - plugin"
subcategory: ""
description: |-
  
---

# nuodbaas_project (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the project
- `organization` (String) Name of the organization for which project is created
- `sla` (String) The SLA for the project. Cannot be updated once the project is created.
- `tier` (String) The Tier for the project. Cannot be updated once the project is created.

### Optional

- `maintenance` (Attributes) (see [below for nested schema](#nestedatt--maintenance))

### Read-Only

- `resource_version` (String)

<a id="nestedatt--maintenance"></a>
### Nested Schema for `maintenance`

Optional:

- `expires_in` (String) The time until the project or database is disabled, e.g. 1d
- `is_disabled` (Boolean)