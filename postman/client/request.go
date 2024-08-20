package client

import "strings"

const (
	RequestPath    = CollectionPathGet + "/requests"
	RequestPathGet = RequestPath + "/%s"
)

type Request struct {
	Data RequestData `json:"data"`
}
type RequestData struct {
	CollectionId string             `json:"collection,omitempty"`
	FolderId     string             `json:"folder,omitempty"`
	Id           string             `json:"id,omitempty"`
	Name         string             `json:"name,omitempty"`
	URL          string             `json:"url"`
	BaseURL      string             `json:"-"`
	Description  string             `json:"description"`
	Method       string             `json:"method"`
	QueryParams  []QueryParamHeader `json:"queryParams"`
	Headers      []QueryParamHeader `json:"headerData"`
	Events       []Event            `json:"events"`
}
type QueryParamHeader struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func (c *Request) RequestEncodeId() string {
	return c.Data.CollectionId + IdSeparator + c.Data.Id
}

func RequestDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
