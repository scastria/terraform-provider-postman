package client

const (
	GlobalEnvironmentPath = WorkspacePathGet + "/global-variables"
)

type GlobalEnvironment struct {
	WorkspaceId string        `json:"-"`
	Variables   []EnvVariable `json:"values"`
}

type EnvVariable struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}
