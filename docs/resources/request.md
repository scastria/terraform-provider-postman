# Resource: postman_request
Represents a request
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
resource "postman_folder" "Folder" {
  collection_id = postman_collection.Collection.collection_id
  name = "My Folder"
  description = "This is my folder"
}
resource "postman_request" "example" {
  collection_id = postman_collection.Collection.collection_id
  folder_id = postman_folder.Folder.folder_id
  name = "My Request"
  method = "GET"
  base_url = "https://postman-echo.com/get"
  body = jsonencode({
    key1 = "value1"
    key2 = "value2"
  })
  query_param {
    key = "p1"
    value = "v1"
    description = "My param"
    enabled = true
  }
  header {
    key = "h1"
    value = "v1"
    description = "My header"
    enabled = true
  }
  pre_request_script = [
    "script line 1",
    "script line 2"
  ]
}
```
## Argument Reference
* `collection_id` - **(Required, ForceNew, String)** The id of the parent collection.
* `name` - **(Required, String)** The name of the request.
* `description` - **(Optional, String)** The description of the request.
* `folder_id` - **(Optional, ForceNew, String)** The parent folder id.
* `method` - **(Optional, String)** The method of the request. Allowed values: `GET`, `PUT`, `POST`, `PATCH`, `DELETE`, `COPY`, `HEAD`, `OPTIONS`, `LINK`, `UNLINK`, `PURGE`, `LOCK`, `UNLOCK`, `PROPFIND`, `VIEW`. Default: `GET`.
* `base_url` - **(Optional, String)** The base url of the request (excluding query params).
* `body` - **(Optional, String)** The raw JSON body of the request.
* `query_param` - **(Optional, list{query_param})** Configuration block for a query_param.  Can be specified multiple times for each query_param.  Each block supports the fields documented below.
* `header` - **(Optional, list{header})** Configuration block for a header.  Can be specified multiple times for each header.  Each block supports the fields documented below.
* `pre_request_script` - **(Optional, List of String)** The JS script to run before the request.
* `post_response_script` - **(Optional, List of String)** The JS script to run after the response (previously called Test scripts).
## query_param
* `key` - **(Required, String)** The name of the query param.
* `value` - **(Optional, String)** The value of the query param.
* `description` - **(Optional, String)** The description of the query param.
* `enabled` - **(Optional, Boolean)** Whether the query param is enabled. Default: `true`
## header
* `key` - **(Required, String)** The name of the header.
* `value` - **(Optional, String)** The value of the header.
* `description` - **(Optional, String)** The description of the header.
* `enabled` - **(Optional, Boolean)** Whether the header is enabled. Default: `true`
## Attribute Reference
* `id` - **(String)** Same as `collection_id`:`request_id`
* `request_id` - **(String)** Id of the request alone
* `url` - **(String)** The url of the request (including enabled query params)
## Import
Requests can be imported using a proper value of `id` as described above
