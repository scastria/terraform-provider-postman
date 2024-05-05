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
  description = "Desc2"
}

resource "postman_folder" "Folder" {
  collection_id = postman_collection.Collection.collection_id
  name = "ShawnTest"
}

resource "postman_request" "Request" {
  collection_id = postman_collection.Collection.collection_id
  folder_id = postman_folder.Folder.folder_id
  name = "My Request"
  method = "GET"
  base_url = "https://postman-echo.com/get"
  query_param {
    key = "p1"
    value = "v1"
    description = "My param"
    enabled = false
  }
  query_param {
    key = "p2"
    value = "v2"
    enabled = true
  }
  query_param {
    key = "p3"
    value = "v3"
    enabled = true
  }
}
