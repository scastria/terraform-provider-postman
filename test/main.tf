terraform {
  required_providers {
    postman = {
      source = "github.com/scastria/postman"
    }
  }
}

provider "postman" {
}

# data "postman_workspace" "Workspace" {
#   name = "Data Team APIs"
#   search_name = "place"
# }

# output "workspace_id" {
#   value = data.postman_workspace.Workspace
# }

# resource "postman_workspace" "Workspace" {
#   name = "ShawnTest"
#   type = "personal"
# }

# resource "postman_global_environment" "Global" {
#   workspace_id = postman_workspace.Workspace.id
#   var {
#     enabled = true
#     key = "cream"
#     type = "default"
#     value = "pie"
#   }
# }

# resource "postman_collection" "Collection" {
#   workspace_id = postman_workspace.Workspace.id
#   name = "ShawnTest"
#   description = "Desc"
#   var {
#     key = "url_base"
#     value = "https://postman-echo.com"
#   }
#   pre_request_script = [
#     "script1",
#     "script2"
#   ]
#   post_response_script = [
#     "script1",
#     "script2"
#   ]
# }

# resource "postman_collection" "Collection2" {
#   workspace_id = postman_workspace.Workspace.id
#   name = "ShawnTest2"
#   description = "Desc2"
#   pre_request_script = [
#     "script1",
#     "script2"
#   ]
# }

# resource "postman_folder" "Folder" {
#   for_each = toset(["a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"])
#   collection_id = postman_collection.Collection.collection_id
#   name = each.key
# }
# resource "postman_folder" "ScriptFolder" {
#   collection_id = postman_collection.Collection.collection_id
#   name = "WithScripts"
#   pre_request_script = [
#       "script1",
#       "script2"
#   ]
# }

# resource "postman_collection_sort" "CollectionSort" {
#   collection_id = postman_collection.Collection.collection_id
#   order = "ASC"
#   hash = timestamp()
#   case_sensitive = true
#   depends_on = [postman_folder.Folder]
# }

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
#   pre_request_script = [
#       "script1",
#       "script2"
#   ]
# }
