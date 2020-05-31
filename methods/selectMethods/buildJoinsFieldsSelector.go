package selectMethods

import (
	"github.com/Alvarios/nefts-go/config"
	"strings"
)

// Filter fields in joined tables.
func buildJoinsFieldsSelector(joins []config.Join) (string, *config.Error) {
	joinsString := ""

	// Ignore when no joins
	if joins != nil && len(joins) > 0 {
		// Multiple joins can be assigned to a common key. We need to parse user parameters for better efficiency.
		type Element struct {
			Bucket string
			Fields []string
		}

		// Elements will be grouped by destination key.
		remuxQueryElements := map[string][]Element{}
		var currentElement Element

		// Check each join and add it to the computable object above.
		for _, value := range joins {
			if value.Bucket == "" {
				return "", &config.Error{
					Code:    400,
					Message: "Missing bucket parameter in joined fields.",
				}
			}

			if value.Fields != nil && len(value.Fields) > 0 {
				currentElement = Element{
					Bucket: value.Bucket,
					Fields: value.Fields,
				}

				// Destination key is already assigned : merge results.
				if _, ok := remuxQueryElements[value.DestinationKey]; ok {
					remuxQueryElements[value.DestinationKey] = append(
						remuxQueryElements[value.DestinationKey],
						currentElement,
					)
				} else {
					// Create a new key to assign.
					remuxQueryElements[value.DestinationKey] = []Element{currentElement}
				}
			}
		}

		// Parse destination key, and the fields they are assigned to.
		for key, value := range remuxQueryElements {
			// Add separators between score declarations, as required by N1QL syntax. Values for a given key are grouped
			// under curly braces.
			if len(joinsString) > 0 {
				joinsString += ", {"
			} else {
				joinsString += "{"
			}

			// An element corresponds to a particular set in a specific bucket.
			for _, element := range value {
				// Add each selected field from foreign bucket.
				for _, field := range element.Fields {
					alias := field
					if field == "*" || field == "all" {
						field = "*"
						alias = "data"
					}

					if strings.HasPrefix(field, "meta.") {
						element.Bucket = "META(`" + element.Bucket + "`)"
						field = field[5:]
					} else {
						element.Bucket = "`" + element.Bucket + "`"
					}

					joinsString += "\"" + alias + "\": " + element.Bucket + "." + field + ", "
				}
			}

			// Group under a same key, and remove the last separator.
			joinsString = joinsString[:len(joinsString) - 2] + "} AS " + key
		}
	}

	return joinsString, nil
}
