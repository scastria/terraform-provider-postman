package client

import "strings"

const (
	CollectionPath    = "collections"
	CollectionPathGet = CollectionPath + "/%s"
	WorkspaceParam    = "workspace"
	CollectionSchema  = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
)

type Collection struct {
	Info CollectionInfo `json:"info"`
}

type CollectionCreate struct {
	Info CollectionInfo `json:"info"`
	// Not actually strings, but we don't need to care about the structure
	Items []string `json:"item"`
}

type CollectionInfo struct {
	WorkspaceId string `json:"-"`
	Id          string `json:"_postman_id,omitempty"`
	CreateId    string `json:"id,omitempty"`
	Uid         string `json:"uid,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema,omitempty"`
}

type CollectionContainer struct {
	Child Collection `json:"collection"`
}

type CollectionCreateContainer struct {
	Child CollectionCreate `json:"collection"`
}

type CollectionInfoContainer struct {
	Child CollectionInfo `json:"collection"`
}

func (c *CollectionInfo) CollectionEncodeId() string {
	return c.WorkspaceId + IdSeparator + c.Id
}

func CollectionDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
