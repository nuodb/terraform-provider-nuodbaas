---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nuodbaas_backup Data Source - nuodbaas"
subcategory: ""
description: |-
  Data source for exposing information about NuoDB backups created using the DBaaS Control Plane
---

# nuodbaas_backup (Data Source)

Data source for exposing information about NuoDB backups created using the DBaaS Control Plane

## Example Usage

```terraform
# Data source that returns the attributes of a specific backup
data "nuodbaas_backup" "backup_details" {
  organization = nuodbaas_backup.backup.organization
  project      = nuodbaas_backup.backup.project
  database     = nuodbaas_backup.backup.database
  name         = nuodbaas_backup.backup.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `database` (String) The database that the backup belongs to
- `name` (String) The name of the backup
- `organization` (String) The organization that the backup belongs to
- `project` (String) The project that the backup belongs to

### Read-Only

- `import_source` (Attributes) (see [below for nested schema](#nestedatt--import_source))
- `labels` (Map of String) User-defined labels attached to the resource that can be used for filtering
- `status` (Attributes) (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--import_source"></a>
### Nested Schema for `import_source`

Read-Only:

- `backup_handle` (String) The existing backup handle to import
- `backup_plugin` (String) The plugin used to create the backup to import


<a id="nestedatt--status"></a>
### Nested Schema for `status`

Read-Only:

- `backup_handle` (String) The handle for the backup
- `backup_plugin` (String) The plugin used to manage the backup
- `created_by_policy` (String) The fully-qualified name of the backup policy that the backup was created by
- `creation_time` (String) The time that the backup was taken
- `message` (String) Message summarizing the state of the backup
- `ready_to_use` (Boolean) Whether the backup is ready to be used to restore a database
- `retained_as` (List of String) The matching retention cycles by this backup
- `state` (String) The state of the backup:
  * `Pending` - The backup is pending completion
  * `Succeeded` - The backup completed successfully and is available for use
  * `Failed` - The backup failed and is unusable
  * `Deleting` - The backup has been marked for deletion, which is in progress