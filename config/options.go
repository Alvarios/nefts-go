package config

type Options struct {
	Fields []string `json:"fields"`
	Where []string `json:"where"`
	Joins []Join `json:"joins"`
	Order map[string]string `json:"order"`
	QueryString string `json:"query_string"`
	Labels map[string][]string `json:"labels"`
	LabelsOptions map[string]LabelOption `json:"labels_options"`
	Config Config `json:"config"`
}
