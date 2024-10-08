---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nuodbaas_backuppolicy Resource - nuodbaas"
subcategory: ""
description: |-
  Resource for managing NuoDB backup policies created using the DBaaS Control Plane
---

# nuodbaas_backuppolicy (Resource)

Resource for managing NuoDB backup policies created using the DBaaS Control Plane

## Example Usage

```terraform
# A backup policy with minimal configuration
resource "nuodbaas_backuppolicy" "basic" {
  organization = "org"
  name         = "basic"
  frequency    = "@weekly"
  selector = {
    scope = "org"
  }
}

# A backup policy with explicit configuration for various attributes
resource "nuodbaas_backuppolicy" "pol" {
  organization = "org"
  name         = "pol"
  labels = {
    "provisioned-by" : "terraform"
  }
  frequency = "@daily"
  selector = {
    scope = "org"
    slas  = ["qa", "prod"]
    tiers = ["n0.small", "n1.small"]
    labels = {
      "rpo" : "1d"
    }
  }
  retention = {
    hourly  = 24
    daily   = 7
    weekly  = 4
    monthly = 12
    yearly  = 3
  }
  suspended = false
  properties = {
    propagate_policy_labels   = true
    propagate_database_labels = true
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `frequency` (String) The frequency to schedule backups at, in cron format
- `name` (String) The name of the backup policy
- `organization` (String) The organization that the backup policy belongs to
- `selector` (Attributes) (see [below for nested schema](#nestedatt--selector))

### Optional

- `labels` (Map of String) User-defined labels attached to the resource that can be used for filtering
- `properties` (Attributes) (see [below for nested schema](#nestedatt--properties))
- `retention` (Attributes) (see [below for nested schema](#nestedatt--retention))
- `suspended` (Boolean) Whether backups from the policy are suspended

### Read-Only

- `status` (Attributes) (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--selector"></a>
### Nested Schema for `selector`

Required:

- `scope` (String) The scope that the backup policy applies to

Optional:

- `labels` (Map of String) The user-defined labels to filter databases on
- `slas` (List of String) The SLAs to filter databases on
- `tiers` (List of String) The tiers to filter databases on


<a id="nestedatt--properties"></a>
### Nested Schema for `properties`

Optional:

- `propagate_database_labels` (Boolean) Whether to propagate the user-defined labels from the matching database to backup resources created by this policy
- `propagate_policy_labels` (Boolean) Whether to propagate the user-defined labels from the backup policy to backup resources created by this policy


<a id="nestedatt--retention"></a>
### Nested Schema for `retention`

Optional:

- `daily` (Number) The number of daily backups to retain
- `hourly` (Number) The number of hourly backups to retain
- `monthly` (Number) The number of monthly backups to retain
- `settings` (Attributes) (see [below for nested schema](#nestedatt--retention--settings))
- `weekly` (Number) The number of weekly backups to retain
- `yearly` (Number) The number of yearly backups to retain

<a id="nestedatt--retention--settings"></a>
### Nested Schema for `retention.settings`

Optional:

- `day_of_week` (String) The day of the week used to promote backup to weekly
- `month` (String) The month of the year used to promote backup to yearly
- `promote_latest_to_daily` (Boolean) Whether to promote the latest backup within the day if multiple backups exist for that day
- `promote_latest_to_hourly` (Boolean) Whether to promote the latest backup within the hour if multiple backups exist for that hour
- `promote_latest_to_monthly` (Boolean) Whether to promote the latest backup within the month if multiple backups exist for that month
- `relative_to_last` (Boolean) Whether to apply the backup rotation scheme relative to the last successful backup instead to the current time



<a id="nestedatt--status"></a>
### Nested Schema for `status`

Read-Only:

- `last_missed_backups` (Attributes List) The last database backups that were not scheduled by this policy (see [below for nested schema](#nestedatt--status--last_missed_backups))
- `last_missed_schedule_time` (String) The time that backups were last missed by this policy
- `last_schedule_time` (String) The time that backups were last taken by this policy
- `next_schedule_time` (String) The time that backups are next scheduled by this policy

<a id="nestedatt--status--last_missed_backups"></a>
### Nested Schema for `status.last_missed_backups`

Read-Only:

- `database` (String) The fully-qualified database name for which a backup was missed by this policy
- `message` (String) A human readable message indicating details about the missed backup by this policy
- `missed_time` (String) The time that a backup was missed by this policy
- `reason` (String) A programmatic identifier indicating the reason for missing a backup by this policy

## Import

Import is supported using the following syntax:

```shell
# An existing backup policy can be imported by specifying the organization
# and policy name, separated by "/"
terraform import nuodbaas_backuppolicy.pol org/pol
```
