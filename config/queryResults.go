package config

type Boundaries struct {
	Start int64 `json:"start"`
	End int64 `json:"end"`
}

type Flags struct {
	BeginningOfResults bool `json:"beginning_of_results"`
	EndOfResults bool `json:"end_of_results"`
}

type QueryResults struct {
	Results []interface{} `json:"results"`
	Boundaries Boundaries `json:"boundaries"`
	Flags Flags `json:"flags"`
}