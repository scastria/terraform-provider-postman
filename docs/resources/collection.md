# Resource: postman_collection
Represents a collection
## Example usage
```hcl
resource "postman_workspace" "Workspace" {
  name = "My Workspace"
  type = "personal"
  description = "This is my workspace"
}
resource "postman_collection" "example" {
  workspace_id = postman_workspace.Workspace.id
  name = "My Collection"
  description = "This is my collection"
}
```
## Argument Reference
* `workspace_id` - **(Required, ForceNew, String)** The id of the parent workspace.
* `name` - **(Required, String)** The name of the collection.
* `description` - **(Optional, String)** The description of the collection.
## Attribute Reference
* `id` - **(String)** Same as `workspace_id`:`collection_id`
* `collection_id` - **(String)** Id of the collection alone
* `uid` - **(String)** Uid of the collection alone (includes owner id)
## Import
Collections can be imported using a proper value of `id` as described above
