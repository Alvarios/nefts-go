package config

type JoinQuery struct {
	Config Config `json:"config"`
	Options Options `json:"options"`
	Start int64 `json:"start"`
	End int64 `json:"end"`
}

type Join struct {
	// Bucket where the joined table is located.
	Bucket string `json:"bucket" required:"true"`
	// Reference key in the parent table.
	ForeignKey string `json:"foreign_key"`
	// Key to assign joined fields in the result tuple.
	DestinationKey string `json:"destination_key"`
	// Fields to filter.
	Fields []string `json:"fields" required:"true"`
	// Specify parent if not main bucket (for nested joins).
	ForeignParent string `json:"foreign_parent"`
	// Reference key in the joined table.
	JoinKey string `json:"join_key"`
	// Nest another request in the join
	JoinQuery *JoinQuery `json:"join_query"`
}
