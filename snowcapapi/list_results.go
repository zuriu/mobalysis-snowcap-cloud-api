package snowcapapi

type ListResults struct {
	TotalCount int64       `json:"totalCount"`
	Results    interface{} `json:"results"`
}

type EventListResults struct {
	TotalCount int64   `json:"totalCount"`
	Results    []Event `json:"results"`
}
