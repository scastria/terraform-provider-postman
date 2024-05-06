package client

import "strings"

const (
	CollectionPath    = "collections"
	CollectionPathGet = CollectionPath + "/%s"
	WorkspaceParam    = "workspace"
	CollectionSchema  = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
)

type Collection struct {
	Info      CollectionInfo `json:"info"`
	Variables []Variable     `json:"variable"`
	Events    []Event        `json:"event"`
}

type CollectionUpdate struct {
	Info      CollectionInfo `json:"info"`
	Variables []Variable     `json:"variables"`
	Events    []Event        `json:"events"`
}

type CollectionCreate struct {
	Info      CollectionInfo `json:"info"`
	Items     []interface{}  `json:"item"`
	Variables []Variable     `json:"variables"`
	Events    []Event        `json:"events"`
}

type Event struct {
	Listen string `json:"listen"`
	Script Script `json:"script"`
}

type Script struct {
	Id   string   `json:"id"`
	Type string   `json:"type"`
	Exec []string `json:"exec"`
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

type Variable struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Type     string `json:"type"`
	Disabled bool   `json:"disabled"`
}
type CollectionContainer struct {
	Child Collection `json:"collection"`
}

type CollectionCreateContainer struct {
	Child CollectionCreate `json:"collection"`
}

type CollectionUpdateContainer struct {
	Child CollectionUpdate `json:"collection"`
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
