package singleservice

type ReadRequest struct {
	K string `json:"k"`
}

type ReadResponse struct {
	V string `json:"v"`
}

type WriteRequest struct {
	K string `json:"k"`
	V string `json:"v"`
}
