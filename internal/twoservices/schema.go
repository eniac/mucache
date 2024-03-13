package twoserivces

type ReadRequest struct {
	K int `json:"k"`
}

type ReadResponse struct {
	V int `json:"v"`
}

type ReadBulkRequest struct {
	Ks []int `json:"ks"`
}

type ReadBulkResponse struct {
	Vs []int `json:"vs"`
}

type WriteRequest struct {
	K int `json:"k"`
	V int `json:"v"`
}

type WriteBulkRequest struct {
	Ks []int `json:"ks"`
	Vs []int `json:"vs"`
}

type HitOrMissRequest struct {
	K       int     `json:"k"`
	HitRate float32 `json:"hit_rate"` // [0.0, 1.0)
}
