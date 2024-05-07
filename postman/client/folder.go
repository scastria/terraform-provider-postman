package client

import "strings"

const (
	FolderPath    = CollectionPathGet + "/folders"
	FolderPathGet = FolderPath + "/%s"
)

type Folder struct {
	Data FolderData `json:"data"`
}
type FolderData struct {
	CollectionId   string  `json:"collection,omitempty"`
	ParentFolderId string  `json:"folder,omitempty"`
	Id             string  `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	Description    string  `json:"description"`
	Events         []Event `json:"events"`
}

func (c *Folder) FolderEncodeId() string {
	return c.Data.CollectionId + IdSeparator + c.Data.Id
}

func FolderDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
