# Data Source: postman_workspace
Represents a workspace
## Example usage
```hcl
data "postman_workspace" "example" {
  search_name = "Work"
  type = "team"
}
```
## Argument Reference
* `search_name` - **(Optional, String)** The search string to apply to the name of the workspace. Uses contains.
* `name` - **(Optional, String)** The filter string to apply to the name of the workspace. Uses equality.
* `type` - **(Optional, String)** The type of the workspace to filter on.  Allowed values: `personal`, `private`, `public`, `team`, `partner`.
## Attribute Reference
* `id` - **(String)** Guid
