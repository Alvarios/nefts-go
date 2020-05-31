package selectMethods

import (
	"nefts"
	"github.com/Alvarios/nefts-go/config"
	"strings"
)

func linkJoins(
	joins []config.Join,
	selector string,
) (string, *config.Error) {
	output := ""

	for _, value := range joins {
		// Add separator (white space is the separator for joins).
		if len(output) > 0 {
			output += " "
		}

		if value.JoinQuery == nil {
			output += "LEFT JOIN `" + value.Bucket + "` ON "
		} else {
			nestedQueryString, err := nefts.BuildQueryString(
				value.JoinQuery.Config,
				value.JoinQuery.Start,
				value.JoinQuery.End,
				value.JoinQuery.Options,
			)

			if err != nil {
				return "", err
			}

			output += "JOIN(" + nestedQueryString + ") AS `" + value.DestinationKey + "` ON "
		}

		localSelector := selector

		// Link parent.
		if len(value.ForeignParent) > 0 {
			localSelector = value.ForeignParent
		}

		if value.ForeignKey == "" {
			return "", &config.Error{
				Code:    400,
				Message: "Missing foreign key parameter in joined fields.",
			}
		}

		if strings.HasPrefix(value.ForeignKey, "meta.") {
			value.ForeignKey = value.ForeignKey[5:]
			localSelector = "META(" + localSelector + ")"
		}

		output += localSelector + "." + value.ForeignKey + " = "

		if len(value.JoinKey) > 0 {
			if strings.HasPrefix(value.JoinKey, "meta.") {
				output += "META(`" + value.Bucket + "`)." + value.JoinKey[5:]
			} else {
				output += "`" + value.Bucket + "`." + value.JoinKey
			}
		} else {
			// By default, join on tuple id.
			output += "META(`" + value.Bucket + "`).id"
		}
	}

	return output, nil
}
