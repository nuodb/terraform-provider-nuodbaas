---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nuodbaas_backuppolicies Data Source - nuodbaas"
subcategory: ""
description: |-
  Data source for listing NuoDB backup policies created using the DBaaS Control Plane
---

# nuodbaas_backuppolicies (Data Source)

Data source for listing NuoDB backup policies created using the DBaaS Control Plane

## Example Usage

```terraform
# Data source that returns the fully-qualified names of all backup policies
data "nuodbaas_backuppolicies" "policy_list" {}

# Data source that returns the fully-qualified names of backup policies within an organization
data "nuodbaas_backuppolicies" "org_policy_list" {
  filter = {
    organization = "org"
  }
}

# Data source that returns the fully-qualified names of backup policies satisfying label requirements
data "nuodbaas_backuppolicies" "label_policy_list" {
  filter = {
    labels = ["withkey", "key=expected", "key!=unexpected", "!withoutkey"]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Attributes) Filters to apply to policies (see [below for nested schema](#nestedatt--filter))

### Read-Only

- `policies` (Attributes List) The list of policies that satisfy the filter requirements (see [below for nested schema](#nestedatt--policies))

<a id="nestedatt--filter"></a>
### Nested Schema for `filter`

Optional:

- `labels` (List of String) List of filters to apply based on labels, which are composed using `AND`. Acceptable filter expressions are:
  * `key` - Only return items that have label with specified key
  * `key=value` - Only return items that have label with specified key set to value
  * `!key` - Only return items that do _not_ have label with specified key
  * `key!=value` - Only return items that do _not_ have label with specified key set to value
- `organization` (String) The organization to filter policies on


<a id="nestedatt--policies"></a>
### Nested Schema for `policies`

Read-Only:

- `name` (String) The name of the policy
- `organization` (String) The organization the policy belongs to
