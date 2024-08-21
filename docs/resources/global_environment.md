# Resource: postman_global_environment
Represents the global environment of a workspace
## Example usage
```hcl
resource "postman_workspace" "Workspace" {
  name = "My Workspace"
  type = "personal"
  description = "This is my workspace"
}
resource "postman_global_environment" "example" {
  workspace_id = postman_workspace.Workspace.id
  var {
    key = "url_base"
    value = "https://postman-echo.com"
    type = "default"
    enabled = true
  }
}
```
## Argument Reference
* `workspace_id` - **(Required, ForceNew, String)** The id of the parent workspace.
* `var` - **(Optional, list{var})** Configuration block for a variable.  Can be specified multiple times for each var.  Each block supports the fields documented below.
## var
* `key` - **(Required, String)** The name of the variable.
* `value` - **(Optional, String)** The value of the variable.
* `type` - **(Optional, String)** The type of variable. Allowed values: `default`, `secret`. Default: `default`
* `enabled` - **(Optional, Boolean)** Whether the variable is enabled. Default: `true`
## Attribute Reference
* `id` - **(String)** Same as `workspace_id`
## Import
Global environments can be imported using a proper value of `id` as described above
