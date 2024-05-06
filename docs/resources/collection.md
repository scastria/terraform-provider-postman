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
  var {
    key = "url_base"
    value = "https://postman-echo.com"
  }
  pre_request_script = [
    "script line 1",
    "script line 2"
  ]
}
```
## Argument Reference
* `workspace_id` - **(Required, ForceNew, String)** The id of the parent workspace.
* `name` - **(Required, String)** The name of the collection.
* `description` - **(Optional, String)** The description of the collection.
* `pre_request_script` - **(Optional, List of String)** The JS script to run before the request.
* `post_response_script` - **(Optional, List of String)** The JS script to run after the response (previously called Test scripts).
* `var` - **(Optional, list{var})** Configuration block for a variable.  Can be specified multiple times for each var.  Each block supports the fields documented below.
## var
* `key` - **(Required, String)** The name of the variable.
* `value` - **(Optional, String)** The value of the variable.
* `disabled` - **(Optional, Boolean)** Whether the variable is disabled. Default: `false`
## Attribute Reference
* `id` - **(String)** Same as `workspace_id`:`collection_id`
* `collection_id` - **(String)** Id of the collection alone
* `uid` - **(String)** Uid of the collection alone (includes owner id)
## Import
Collections can be imported using a proper value of `id` as described above
