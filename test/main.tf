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
}
