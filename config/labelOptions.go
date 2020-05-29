package config

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
