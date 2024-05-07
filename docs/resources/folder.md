# Resource: postman_folder
Represents a folder
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
resource "postman_folder" "example" {
  collection_id = postman_collection.Collection.collection_id
  name = "My Folder"
  description = "This is my folder"
  pre_request_script = [
    "script line 1",
    "script line 2"
  ]
}
```
## Argument Reference
* `collection_id` - **(Required, ForceNew, String)** The id of the parent collection.
* `name` - **(Required, String)** The name of the folder.
* `description` - **(Optional, String)** The description of the folder.
* `parent_folder_id` - **(Optional, ForceNew, String)** The parent folder id.
* `pre_request_script` - **(Optional, List of String)** The JS script to run before the request.
* `post_response_script` - **(Optional, List of String)** The JS script to run after the response (previously called Test scripts).
## Attribute Reference
* `id` - **(String)** Same as `collection_id`:`folder_id`
* `folder_id` - **(String)** Id of the folder alone
## Import
Folders can be imported using a proper value of `id` as described above
