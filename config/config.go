package config

import "github.com/couchbase/gocb/v2"

type Parameters struct {
	MaxQueryLength int `json:"max_query_length" default:"1000"`
	Debug bool `json:"debug" default:"false"`
}

type LabelOption struct {
	Analyzer string `json:"analyzer"`
	Fuzziness string `json:"fuzziness"`
	Out string `json:"out"`
	PrefixLength string `json:"prefix_length"`
	PhraseMode bool `json:"phrase_mode"`
	RegexpMode bool `json:"regexp_mode"`
	Weight string `json:"weight"`
	Bucket string `json:"bucket"`
}

type Config struct {
	Cluster *gocb.Cluster `json:"cluster"`
	Bucket string `json:"bucket"`
	Parameters Parameters `json:"parameters"`
	Labels map[string][]string `json:"labels"`
	LabelsOptions map[string]LabelOption `json:"labels_options"`
}
