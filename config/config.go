package config

import "github.com/couchbase/gocb/v2"

type Parameters struct {
	MaxQueryLength int `json:"max_query_length" default:"1000"`
	Debug bool `json:"debug" default:"false"`
}



type Config struct {
	Cluster *gocb.Cluster `json:"cluster"`
	Bucket string `json:"bucket"`
	Parameters Parameters `json:"parameters"`
}
