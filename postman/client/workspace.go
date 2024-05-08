package client

const (
	WorkspacePath    = "workspaces"
	WorkspacePathGet = WorkspacePath + "/%s"
)

type Workspace struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

type WorkspaceContainer struct {
	Child Workspace `json:"workspace"`
}

type WorkspaceCollection struct {
	Data []Workspace `json:"workspaces"`
}
