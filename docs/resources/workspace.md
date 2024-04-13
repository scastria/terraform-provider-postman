# Resource: postman_workspace
Represents a workspace
## Example usage
```hcl
resource "postman_workspace" "example" {
  name = "My Workspace"
  type = "personal"
  description = "This is my workspace"
}
```
## Argument Reference
* `name` - **(Required, String)** The name of the workspace.
* `type` - **(Required, String)** The type of the workspace.  Allowed values: `personal`, `private`, `public`, `team`, `partner`.
* `description` - **(Optional, String)** The description of the workspace.
## Attribute Reference
* `id` - **(String)** Guid
## Import
Workspaces can be imported using a proper value of `id` as described above
