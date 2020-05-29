package whereMethods

import (
	"nefts/config"
	"nefts/utils"
)

func BuildWhereQuery(
	query map[string]string,
	where []string,
	labelOptions map[string]config.LabelOption,
) string {
	output := ""

	// No need to add filters when query string is empty.
	if query != nil && len(query) > 0 {
		for key, value := range query {
			// Pre-format value for better efficiency.
			formatted := utils.SoftFormat(value)

			// Add separator. String needs to start with WHERE clause.
			if len(output) > 0 {
				output += " AND "
			} else {
				output += "WHERE "
			}

			labelParams := config.LabelOption{}

			// Search if some custom parameter was provided for current label
			if labelOptions != nil {
				match, ok := labelOptions[key]
				if ok {
					labelParams = match
				} else {
					match, ok = labelOptions["global"]
					if ok {
						labelParams = match
					}
				}
			}

			// Add default if absent.
			if labelParams.Analyzer == "" {
				labelParams.Analyzer = "standard"
			}

			if labelParams.Fuzziness == "" {
				labelParams.Fuzziness = "2"
			}

			if labelParams.Bucket == "" {
				labelParams.Out = "b"
			}

			method := "match"

			if labelParams.PhraseMode {
				method = "match_phrase"
			} else if labelParams.RegexpMode {
				method = "regexp"
			}

			// Run search with standard option.
			output += "SEARCH(" + labelParams.Bucket + ", {\"" + method + "\": \"" + formatted + "\""
			// Restrict search to a specific subset.
			if key != "general" {
				output += ", \"field\": \"" + key + "\""
			}

			output += ", \"analyzer\": \"" + labelParams.Analyzer +
				"\", \"fuzziness\": " + labelParams.Fuzziness +
				", \"out\": \"" + key + "_out\"})"
		}
	}

	// Add additional filters.
	if where != nil && len(where) > 0 {
		for _, value := range where {
			if len(output) > 0 {
				output += " AND "
			} else {
				output += "WHERE "
			}

			output += value
		}
	}

	return output
}
