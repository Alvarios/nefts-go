package selectMethods

import (
	"nefts/config"
	"nefts/utils"
	"regexp"
)

// Get match scores for Full Text Searches.
func buildScoreSelector(
	query map[string]string,
	labelOptions map[string]config.LabelOption,
) string {
	scoreString := ""

	if query != nil && len(query) > 0 {
		// Query is an object with search fields per label. {...label: queryString}
		for key, value := range query {
			// Value needs to contain characters other than white space.
			if len(utils.SoftFormat(value)) > 0 {
				// Add separators between score declarations, as required by N1QL syntax.
				if len(scoreString) > 0 {
					scoreString += ", "
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

				if labelParams.Out == "" {
					labelParams.Out = "{key}_query_score"
				}

				reOut := regexp.MustCompile(`\{key}`)
				formattedKey := reOut.ReplaceAllString(labelParams.Out, key)

				// Add score selector (note they will appear in the fetched tuples).
				scoreString += "SEARCH_SCORE(" + key + "_out) AS " + formattedKey
			}
		}
	}

	return scoreString
}

