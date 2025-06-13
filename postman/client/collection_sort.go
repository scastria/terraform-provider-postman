package client

type CollectionSortContainer struct {
	Child CollectionSort `json:"collection"`
}

type CollectionSort struct {
	Info      CollectionInfo           `json:"info"`
	Items     []map[string]interface{} `json:"item"`
	Variables []interface{}            `json:"variable"`
	Events    []interface{}            `json:"event"`
}
