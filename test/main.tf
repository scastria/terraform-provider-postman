terraform {
  required_providers {
    postman = {
      source = "github.com/scastria/postman"
    }
  }
}

provider "postman" {
}

resource "postman_workspace" "Workspace" {
  name = "ShawnTest"
  type = "personal"
}

resource "postman_collection" "Collection" {
  workspace_id = postman_workspace.Workspace.id
  name = "ShawnTest"
  description = "Desc"
  var {
    key = "url_base"
    value = "https://postman-echo.com"
  }
}

resource "postman_folder" "Folder" {
  for_each = toset(["a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"])
  collection_id = postman_collection.Collection.collection_id
  name = each.key
}

resource "postman_collection_sort" "CollectionSort" {
  collection_id = postman_collection.Collection.collection_id
  order = "ASC"
  depends_on = [postman_folder.Folder]
}

# resource "postman_request" "Request" {
#   collection_id = postman_collection.Collection.collection_id
#   folder_id = postman_folder.Folder["a"].folder_id
#   name = "My Request"
#   method = "GET"
#   base_url = "{{url_base}}/get"
#   query_param {
#     key = "p1"
#     value = "v1"
#     description = "My param"
#     enabled = true
#   }
#   query_param {
#     key = "p2"
#     value = "v2"
#     enabled = false
#   }
#   query_param {
#     key = "p3"
#     value = "v3"
#     enabled = true
#   }
# }
