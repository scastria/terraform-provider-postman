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
