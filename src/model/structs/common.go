package structs

type PageListResp struct {
	Data  interface{} `json:"data"`
	Count int64       `json:"count"`
}
