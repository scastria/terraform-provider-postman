# Resource: postman_collection_sort
Represents the sorting of a collection's direct contents
## Example usage
```hcl
resource "postman_workspace" "Workspace" {
  name = "My Workspace"
  type = "personal"
  description = "This is my workspace"
}
resource "postman_collection" "Collection" {
  workspace_id = postman_workspace.Workspace.id
  name = "My Collection"
  description = "This is my collection"
}
resource "postman_collection_sort" "example" {
  collection_id = postman_collection.Collection.collection_id
  order = "ASC"
}
```
## Argument Reference
* `collection_id` - **(Required, ForceNew, String)** The id of the collection.
* `order` - **(Optional, String)** The order in which to sort the direct contents. Allowed values: `ASC`, `DESC`, `NONE`. Default: `ASC`.
* `hash` - **(Optional, String)** A hash value to trigger a resort.
## Attribute Reference
* `id` - **(String)** Same as `collection_id`
## Import
Collection sorts can be imported using a proper value of `id` as described above
