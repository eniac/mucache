package loadcm

type StartExperiment struct {
	Readers      int `json:"readers"`
	DurationSecs int `json:"duration_secs"`
	//Writers int `json:"writers"`
}

type InvalidationExperimentRequest struct {
	Times   int `json:"times"`
	Timeout int `json:"timeout"`
}
